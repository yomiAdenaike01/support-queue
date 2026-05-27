import asyncio
from logging import getLogger
from typing import TYPE_CHECKING
from .event_bus import EventBus
from .pipeline import Pipeline
from .team import TeamRepository
from .integrations import Integrations
from .logging import configure_logging
from .repositories import TicketRepository

if TYPE_CHECKING:
    from .pipeline import WorkerContext

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
    base_url = "http://localhost:2342/api/v1"
    ticket_repo_url = f"{base_url}/ticket"
    team_repo_url = f"{base_url}/team"
    event_bus_url = "redis://localhost:6379"

    configure_logging()
    
    await health_check(base_url)
    
    ticket_repository = TicketRepository(base_url=ticket_repo_url)
    team_repository = TeamRepository(base_url=team_repo_url)
    
    event_bus = EventBus(url=event_bus_url)
    event_bus.connect()
    integrations = Integrations()
    pipeline = Pipeline(team_repository,integrations)
    while True:
        event = event_bus.await_new_event()
        if event is None: 
            continue
        for _,messages in event:
            for __, message in messages:
                ticket_id = message.get("ticket_id", None)
                
                if ticket_id is None:
                    continue
                
                ticket = await ticket_repository.find_by_id(ticket_id)

                if ticket is None:
                    continue
                
                pipeline_result = await pipeline.run(ticket)
                await post_pipeline_result(base_url, pipeline_result)
                
    
async def post_pipeline_result(base_url: str, ctx: "WorkerContext"):
    from httpx import AsyncClient
    async with AsyncClient() as http:
        url = f"{base_url}/ml"
        response = await http.post(url, json=ctx)
        response.raise_for_status()

if __name__ == "__main__":
    import uvloop
    uvloop.run(init_worker())
