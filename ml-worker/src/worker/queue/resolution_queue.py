import asyncio
from logging import getLogger
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from ..pipeline import Pipeline
    from ..event_bus import StreamEvent
    from ..event_bus import EventBus
    from ..repositories import TicketRepository

logger = getLogger('[resolution-queue]')

class ResolutionQueue:
    _queue: asyncio.Queue["StreamEvent"]
    _pipeline: "Pipeline"
    _event_bus: "EventBus"
    _ticket_repository: "TicketRepository"

    def __init__(self, ticket_repository: "TicketRepository", event_bus:"EventBus", resolution_pipeline: "Pipeline"):
        self._queue = asyncio.Queue(maxsize=20)
        self._pipeline = resolution_pipeline
        self._event_bus = event_bus
        self._ticket_repository = ticket_repository
        

    async def __work(self):
        queue = self._queue
        while True:
            event = await queue.get()
            ticket_id = event.get("data").get("ticket_id")
            try:
                logger.info(f'[resolution-worker]: Processing ticket id={ticket_id}...')
                ticket = await self._ticket_repository.find_by_id(ticket_id)
                logger.info(f'[resolution-worker]: Successfully found ticket by id={ticket_id}...')
                await self._pipeline.run(ticket)
                self._event_bus.ack_resolution(message_id=event.get("id"))
            except Exception as e:
                logger.exception("[resolution-worker]: Failed processing ticket id=%s reason=%s", ticket_id, str(e))
            finally:
                queue.task_done()

    def begin_workers(self, num_workers = 3):
        logger.info("[resolution-worker]: Resolution workers initialised.")
        return [asyncio.create_task(self.__work()) for _ in range(num_workers)]
    
    async def add_to_queue(self, event: "StreamEvent"):
        await self._queue.put(event)