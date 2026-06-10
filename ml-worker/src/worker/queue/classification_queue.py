import asyncio
from logging import getLogger
from typing import TYPE_CHECKING

from ..repositories import ClassificationResult

if TYPE_CHECKING:
    from ..event_bus import StreamEvent
    from ..pipeline import WorkerContext
    from ..repositories import UpdateTicketEvent
    from .models import QueueDependencies
logger = getLogger("[classification-queue]")


class ClassificationQueue:
    _queue: asyncio.Queue["StreamEvent"]
    _deps: "QueueDependencies"

    def __init__(self, deps: "QueueDependencies"):
        self._queue = asyncio.Queue(maxsize=20)
        self._deps = deps

    def begin_workers(self, num_workers: int = 3):
        return [asyncio.create_task(self.__work()) for _ in range(num_workers)]

    def _to_classification_result(self, ctx: "WorkerContext") -> "ClassificationResult":
        return ClassificationResult(
            ticket_id=ctx.ticket.id,
            suggested_response=ctx.suggested_response,
            average_sentiment=ctx.average_sentiment,
            priority=str(ctx.priority),
            category=ctx.category,
            requires_urgency=ctx.requires_urgency,
            urgency_reason=ctx.urgency_reason,
        )

    async def __work(self):
        max_attempts = int(self._deps.config.get("MAX_ATTEMPTS"))
        while True:
            event = await self._queue.get()
            ticket_id = event.data.get("ticket_id", "")
            try:
                ticket = await self._deps.ticket_repository.find_by_id(ticket_id)
                if ticket is None:
                    continue
                for attempt in range(max_attempts):
                    try:
                        logger.info(
                            "================ Classication Attempt %d =================",
                            attempt,
                        )
                        await self._deps.ticket_repository.set_ticket_status(
                            UpdateTicketEvent(
                                ticket_id=ticket.id, status="TICKET_PROCESSED"
                            )
                        )
                        result = await self._deps.pipeline.run(ticket)
                        if result is None:
                            continue

                        await self._deps.worker_repository.complete_classification(
                            self._to_classification_result(result)
                        )
                        self._deps.event_bus.ack_classification(message_id=event.id)
                        break
                    except Exception as e:
                        logger.exception(
                            "[classification-queue]: Failed attempt reason=%s", str(e)
                        )
                        continue
            finally:
                self._queue.task_done()

    async def add_to_queue(self, data: "StreamEvent"):
        await self._queue.put(data)
