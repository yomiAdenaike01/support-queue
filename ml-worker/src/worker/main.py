import asyncio
from logging import getLogger

from .logging import configure_logging
from .worker import Worker

logger = getLogger("[main]")


async def init_worker():
    configure_logging()
    worker = Worker()
    start_workers_routine = worker.start_workers()
    resolution_listener = asyncio.create_task(worker.on_resolved_ticket_event())
    classification_listener = asyncio.create_task(worker.on_new_ticket())
    await asyncio.gather(
        start_workers_routine, resolution_listener, classification_listener
    )


if __name__ == "__main__":
    import uvloop

    uvloop.run(init_worker())
