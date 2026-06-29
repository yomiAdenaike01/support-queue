import asyncio
from logging import getLogger

from .logging import configure_logging
from .worker import Worker

logger = getLogger("[main]")


async def init_worker():
    configure_logging()
    worker = Worker()
    start_workers_routine = worker.start_workers()
    pending_resolution_worker = asyncio.create_task(
        worker.on_pending_resolved_tickets()
    )
    pending_ticket_worker = asyncio.create_task(worker.on_pending_tickets())
    resolution_listener = asyncio.create_task(worker.on_resolved_ticket_event())
    classification_listener = asyncio.create_task(worker.on_new_ticket())
    await asyncio.gather(
        pending_resolution_worker,
        start_workers_routine,
        resolution_listener,
        classification_listener,
        pending_ticket_worker,
    )


if __name__ == "__main__":
    import uvloop

    uvloop.run(init_worker())
