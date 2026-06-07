import asyncio
import os
from pathlib import Path
from dotenv import load_dotenv
from logging import getLogger
from typing import TYPE_CHECKING, TypedDict
from .event_bus import EventBus
from .pipeline import ClassificationPipeline, ResolutionPipeline
from .integrations import Integrations
from .logging import configure_logging
from .repositories import TicketRepository
from .knowledge_base import KnowledgeBase
from .config import create_config
from time import sleep

if TYPE_CHECKING:
    from .pipeline import WorkerContext

class ClassificationResult(TypedDict):
    ticket_id: str
    suggested_response: str
    average_sentiment: float
    priority: str
    category: str
    requires_urgency: bool
    urgency_reason: str  

logger = getLogger(__name__)

class Application:
    def __init__(self):
        self._config = create_config()
        base_url = self._config.get("BASE_URL")
        
        ticket_repo_url = f"{base_url}/tickets"
        event_bus_url = "redis://localhost:6379"

        configure_logging()
        
        
        self._health_check(base_url)
        self._event_bus = EventBus(url=event_bus_url)
        self._event_bus.connect()
        self._cache = self._event_bus.get_cache()
        self._integrations = Integrations(cache=self._cache)
        knowledge_base = KnowledgeBase(base_url=base_url)
        self._ticket_repository = TicketRepository(base_url=ticket_repo_url)
        
        
        self._classifier = ClassificationPipeline(self._integrations, cache=self._cache, knowledge_base=knowledge_base)
        
        self._resolver = ResolutionPipeline(self._integrations, knowledge_base=knowledge_base)

    def on_create_ticket(self):
        return self._event_bus.await_new_event()
    
    def _health_check(self, base_url: str):
        from httpx import Client
        max_attempts = 3
        default_backoff = 5
        with Client() as http:
            for attempt in range(max_attempts):
                try:
                    response =  http.get(f"{base_url}/healthz")
                    
                    if response.status_code == 200:
                        print("Health check passed!")
                        return True
                        
                except Exception as e:
                    print(f"Request failed: {e}")
                
                if attempt < max_attempts - 1:
                    delay = default_backoff * (attempt + 1)
                    print(f"Attempt {attempt + 1} failed. Retrying in {delay}s...")
                    sleep(delay)
                
        return False 
    
    async def await_resolution(self, queue:"asyncio.Queue"):
        while True:
            ticket_id = await queue.get()
            try:
                if ticket_id is None:
                    continue
                logger.info(f'[resolution-pipeline]: Processing ticket id={ticket_id}...')
                ticket = await self._ticket_repository.find_by_id(ticket_id)
                logger.info(f'[resolution-pipeline]: Successfully found ticket by id={ticket_id}...')
                await self._resolver.run(ticket)
            except Exception as e:
                logger.exception("[resolution-pipeline]: Failed processing ticket id=%s reason=%s", ticket_id, str(e))
            finally:
                queue.task_done()

    async def run_classification_pipeline(self, queue: "asyncio.Queue"):
        MAX_ATTEMPTS = int(self._config.get("MAX_ATTEMPTS"))
        ticket_id = await queue.get()
        if ticket_id is None:
            return
        
        ticket = await self._ticket_repository.find_by_id(ticket_id)
        
        if ticket is None:
            return
        
        for attempt in range(MAX_ATTEMPTS):
            try:
                logger.info("================ PIPELINE ATTEMPT %d =================", attempt)
                await self._ticket_repository.set_ticket_status(ticket_id, 'PROCESSING')
                classification_result = await self._classifier.run(ticket)
                self._event_bus.ack_classification_complete(ticket_id=ticket.id)
                queue.task_done()
                await self._on_complete_classification(self._base_url, ClassificationResult(
                    ticket_id=ticket.id,
                    average_sentiment=classification_result.average_sentiment,
                    suggested_response=classification_result.suggested_response,
                    priority=classification_result.priority,
                    category=classification_result.category,
                    requires_urgency=classification_result.requires_urgency,
                    urgency_reason=classification_result.urgency_reason
                ))
                break
            except Exception as e:
                logger.exception("Failed attempt reason=%s", str(e))
                continue
    async def _on_complete_classification(self, base_url: str, worker_result: "ClassificationResult"):
        from httpx import AsyncClient
        try:
            async with AsyncClient() as http:
                url = f"{base_url}/worker"
                response = await http.post(url, json=worker_result)
                response.raise_for_status()
        except:
            raise
    async def listen(self, queue:"asyncio.Queue"):
        while True:
                event = await asyncio.to_thread(self.on_create_ticket)
                if event is None: 
                    continue
                for _,messages in event:
                    for __, message in messages:
                        ticket_id = message.get("ticket_id", None)
                        
                        if ticket_id is None:
                            continue
                        await queue.put(ticket_id)

async def init_worker():
    configure_logging()
    application = Application()    
    resolution_queue = asyncio.Queue(maxsize=20)

    resolution_workers = [asyncio.create_task(application.await_resolution(resolution_queue)) for _ in range(3)]

    # resolution_listener = asyncio.create_task(application.listen(resolution_queue))
    
    logger.info("Inserting TEST-ID into resolution queue...")
    await resolution_queue.put("1c653cd0-75a0-4cb9-842f-e6b174ff2cec")
    
    await resolution_queue.join()
    await asyncio.gather(*resolution_workers, return_exceptions=True)

if __name__ == "__main__":
    import uvloop
    uvloop.run(init_worker())
