import asyncio
import os
from pathlib import Path
from dotenv import load_dotenv
from logging import getLogger
from typing import TYPE_CHECKING, TypedDict
from .event_bus import EventBus
from .pipeline import Pipeline, PipelineException
from .integrations import Integrations
from .logging import configure_logging
from .repositories import TicketRepository
from .config import create_config

if TYPE_CHECKING:
    from .pipeline import WorkerContext

class WorkerResult(TypedDict):
    suggested_response: str
    average_sentiment: float
    priority: str
    category: str
    requires_urgency: bool
    urgency_reason: str  

logger = getLogger(__name__)

async def health_check(base_url: str):
    from httpx import AsyncClient
    max_attempts = 3
    default_backoff = 5
    
    async with AsyncClient() as http:
        for attempt in range(max_attempts):
            try:
                response = await http.get(f"{base_url}/healthz")
                
                if response.status_code == 200:
                    print("Health check passed!")
                    return True
                    
            except Exception as e:
                print(f"Request failed: {e}")
            
            if attempt < max_attempts - 1:
                delay = default_backoff * (attempt + 1)
                print(f"Attempt {attempt + 1} failed. Retrying in {delay}s...")
                await asyncio.sleep(delay)
                
        return False      



async def init_worker():
    config = create_config()

    base_url = config.get("BASE_URL")
    ticket_repo_url = f"{base_url}/ticket"
    event_bus_url = "redis://localhost:6379"

    configure_logging()
    
    await health_check(base_url)
    
    ticket_repository = TicketRepository(base_url=ticket_repo_url)
    
    event_bus = EventBus(url=event_bus_url)
    event_bus.connect()
    integrations = Integrations()
    pipeline = Pipeline(integrations, cache=event_bus.get_cache())

    MAX_ATTEMPTS = int(config.get("MAX_ATTEMPTS"))
    queue = []
    while True:
        event = event_bus.await_new_event()
        if event is None: 
            continue
        for _,messages in event:
            for __, message in messages:
                ticket_id = message.get("ticket_id", None)
                
                if ticket_id is None:
                    continue
                queue.append(ticket_id)
                ticket = await ticket_repository.find_by_id(ticket_id)

                if ticket is None:
                    continue
            for attempt in range(MAX_ATTEMPTS):
                try:
                    logger.info("================ PIPELINE ATTEMPT %d =================", attempt)
                    pipeline_result = await pipeline.run(ticket)
                    event_bus.register_completion(ticket_id=ticket.id)
                    await post_pipeline_result(base_url, WorkerResult(
                        average_sentiment=pipeline_result.average_sentiment,
                        suggested_response=pipeline_result.suggested_response,
                        priority=pipeline_result.priority,
                        category=pipeline_result.category,
                        requires_urgency=pipeline_result.requires_urgency,
                        urgency_reason=pipeline_result.urgency_reason
                    ))
                    break
                except Exception as e:
                    logger.exception("Failed attempt reason=%s", str(e))
                    continue
              
    
async def post_pipeline_result(base_url: str, worker_result: "WorkerResult"):
    from httpx import AsyncClient
    try:
        async with AsyncClient() as http:
            url = f"{base_url}/worker"
            response = await http.post(url, json=worker_result)
            response.raise_for_status()
    except:
        raise

if __name__ == "__main__":
    import uvloop
    uvloop.run(init_worker())
