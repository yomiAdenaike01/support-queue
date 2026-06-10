import asyncio
from logging import getLogger
from .worker import Worker
from .logging import configure_logging

logger = getLogger('[main]')

async def init_worker():
    configure_logging()
    worker = Worker() 
    resolution_listener = asyncio.create_task(worker.on_resolved_ticket_event())
    classification_listener = asyncio.create_task(worker.on_new_ticket())
    await asyncio.gather(*worker.start_pipeline_workers(), resolution_listener, classification_listener) 

if __name__ == "__main__":
    import uvloop
    uvloop.run(init_worker())
