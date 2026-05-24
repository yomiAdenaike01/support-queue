from .event_bus import EventBus
from .pipeline import Pipeline
from .team import TeamRepository
from .integrations import Integrations
from .logging import configure_logging
from .services import TicketService
from logging import getLogger

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
    configure_logging()
    await health_check(base_url)
    ticket_service = TicketService(base_url=f"{base_url}/support/ticket")
    event_bus = EventBus(url="redis://localhost:6379")
    event_bus.connect()

    pipeline = Pipeline(TeamRepository(),Integrations())
    while True:
        event = event_bus.await_new_event()
        if event is None: 
            continue
        for _,messages in event:
            for __, message in messages:
                ticket_id = message.get("ticket_id", None)
                
                if ticket_id is None:
                    continue

                ticket = await ticket_service.find_ticket_by_id(ticket_id)
                print(ticket)
                if ticket is None:
                    continue
                
                await pipeline.run(ticket)
    

if __name__ == "__main__":
    import asyncio
    asyncio.run(init_worker())
