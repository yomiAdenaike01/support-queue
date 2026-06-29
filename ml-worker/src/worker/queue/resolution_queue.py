import asyncio
from logging import getLogger
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from ..event_bus import StreamEvent
    from ..queue import QueueDependencies

logger = getLogger("[resolution-queue]")


class ResolutionQueue:
    _queue: asyncio.Queue["StreamEvent"]
    _deps: "QueueDependencies"

    def __init__(self, deps: "QueueDependencies"):
        self._queue = asyncio.Queue(maxsize=20)
        self._deps = deps

    async def __work(self):
        queue = self._queue
        while True:
            event = await queue.get()
            ticket_id = str(event.data.get("ticket_id"))
            try:
                logger.info(f"[resolution-worker]: Processing ticket id={ticket_id}...")
                ticket = await self._deps.ticket_repository.find_by_id(ticket_id)

                if ticket is None:
                    continue

                logger.info(
                    f"[resolution-worker]: Successfully found ticket by id={ticket_id}..."
                )
                await self._deps.pipeline.run(ticket)
                self._deps.event_bus.ack_resolution(message_id=event.id)
            except Exception as e:
                logger.exception(
                    "[resolution-worker]: Failed processing ticket id=%s reason=%s",
                    ticket_id,
                    str(e),
                )
            finally:
                queue.task_done()

    def begin_workers(self, num_workers: int = 3):
        logger.info("[resolution-worker]: Resolution workers initialised.")
        return [asyncio.create_task(self.__work())]

    async def add_to_queue(self, event: "StreamEvent"):
        await self._queue.put(event)
