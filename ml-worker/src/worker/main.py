from .event_bus import EventBus
from .pipeline import Pipeline
from .team import TeamRepository
from .integrations import Integrations
from .logging import configure_logging
from .services import TicketService
from logging import getLogger
import json
logger = getLogger(__name__)

async def init_worker():
    configure_logging()
    ticket_service = TicketService(base_url="http://localhost:2342/api/v1/tickets", is_dry_run=True)
    event_bus = EventBus(url="redis://localhost:6379")
    event_bus.connect()

    pipeline = Pipeline(TeamRepository(),Integrations())
    while True:
        event = event_bus.await_new_event()
        logger.info("Listening for messages on TICKETS_STREAM")
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
