import asyncio
from logging import getLogger
from typing import TYPE_CHECKING, Any
from uuid import uuid4

from ..repositories import (
    ClassificationResult,
    TicketEventStatus,
    TicketStatus,
    UpdateTicketEvent,
)

if TYPE_CHECKING:
    from ..event_bus import StreamEvent
    from ..pipeline import WorkerContext
    from .models import QueueDependencies

logger = getLogger("[classification-queue]")


class ClassificationQueue:
    _queue: asyncio.Queue["StreamEvent"]
    _deps: "QueueDependencies"
    _seen: set[Any]

    def __init__(self, deps: "QueueDependencies"):
        self._queue = asyncio.Queue(maxsize=20)
        self._deps = deps
        self._seen = set()

    def begin_workers(self, num_workers: int = 3):
        return [asyncio.create_task(self.__work(uuid4().hex))]

    def _to_classification_result(self, ctx: "WorkerContext") -> "ClassificationResult":
        return ClassificationResult(
            ticket_id=ctx.ticket.id,
            suggested_response=ctx.suggested_response,
            average_sentiment=ctx.average_sentiment,
            priority=ctx.priority.value if ctx.priority is not None else None,
            category=ctx.category,
            requires_urgency=ctx.requires_urgency,
            urgency_reason=ctx.urgency_reason,
        )

    async def __work(self, id: str):
        max_attempts = int(self._deps.config.get("MAX_ATTEMPTS"))
        ticket_id = None
        while True:
            event = await self._queue.get()
            for attempt in range(max_attempts):
                logger.info(
                    f"[classification-worker-{id}] working on attempt={attempt}"
                )
                try:
                    ticket_id = event.data.get("ticket_id", "")
                    ticket = await self._deps.ticket_repository.find_by_id(ticket_id)
                    if ticket is None:
                        break

                    await self._deps.ticket_repository.set_ticket_status(
                        UpdateTicketEvent(
                            ticket_id=ticket.id,
                            event_status=TicketEventStatus.PROCESSING,
                            status=TicketStatus.PROCESSING,
                        )
                    )
                    result = await self._deps.pipeline.run(ticket)
                    if result is None:
                        continue
                    await self._deps.worker_repository.complete_classification(
                        self._to_classification_result(result)
                    )
                    self._deps.event_bus.ack_classification(message_id=event.id)
                    self._queue.task_done()

                    await self._deps.ticket_repository.set_ticket_status(
                        UpdateTicketEvent(
                            ticket_id=ticket_id,
                            status=TicketStatus.PROCESSED,
                            event_status=TicketEventStatus.CLASSIFIED,
                        )
                    )
                    break
                except Exception as error:
                    logger.exception(
                        f"[classification-worker-{id}]: Failed worker task error=%s",
                        str(error),
                    )
                    if attempt + 1 == max_attempts and ticket_id is not None:
                        await self._deps.ticket_repository.set_ticket_status(
                            UpdateTicketEvent(
                                ticket_id=ticket_id,
                                event_status=TicketEventStatus.FAILED,
                                status=TicketStatus.FAILED,
                            )
                        )
                        self._queue.task_done()
                        break
                    continue

    async def add_to_queue(self, data: "StreamEvent"):
        if data.id in self._seen:
            return
        self._seen.add(data.id)
        await self._queue.put(data)
