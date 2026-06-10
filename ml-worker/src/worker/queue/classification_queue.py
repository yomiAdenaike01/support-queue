import asyncio
from typing import TYPE_CHECKING, TypedDict
from logging import getLogger
from httpx import AsyncClient

if TYPE_CHECKING:
    from ..event_bus import StreamEvent
    from .models import QueueDependencies

logger = getLogger('[classification-queue]')


class ClassificationResult(TypedDict):
    ticket_id: str
    suggested_response: str
    average_sentiment: float
    priority: str
    category: str
    requires_urgency: bool
    urgency_reason: str  

class ClassificationQueue:
    _queue: asyncio.Queue["StreamEvent"]
    _deps: "QueueDependencies"

    def __init__(
        self,
        deps: "QueueDependencies"
    ):
        self._queue = asyncio.Queue(maxsize=20)
        self._deps = deps

    def begin_workers(self, num_workers = 3):
        return [asyncio.create_task(self.__work()) for _ in range(num_workers)]
    
    async def __work(self):
        max_attempts = int(self._deps.config.get("MAX_ATTEMPTS"))
        while True:
            event = await self._queue.get()
            ticket_id = event.get("data").get("ticket_id")
            try:
                ticket = await self._deps.ticket_repository.find_by_id(ticket_id)
                if ticket is None:
                    continue
                for attempt in range(max_attempts):
                    try:
                        logger.info("================ Classication Attempt %d =================", attempt)
                        await self._deps.ticket_repository.set_ticket_status({
                            "ticket_id": ticket.id,
                            "status": "TICKET_PROCESSING",
                        })
                        result = await self._deps.pipeline.run(ticket)
                        await self._on_complete_classification(self._deps.config.get("API_BASE_URL"), ClassificationResult(
                            ticket_id=ticket.id,
                            average_sentiment=result.average_sentiment,
                            suggested_response=result.suggested_response,
                            priority=result.priority,
                            category=result.category,
                            requires_urgency=result.requires_urgency,
                            urgency_reason=result.urgency_reason
                        ))
                        self._deps.event_bus.ack_classification(message_id=event.get("id"))
                        break
                    except Exception as e:
                        logger.exception("[classification-queue]: Failed attempt reason=%s", str(e))
                        continue
            finally:
                self._queue.task_done()
    
    async def add_to_queue(self, data:"StreamEvent"):
        await self._queue.put(data)

    async def _on_complete_classification(self, worker_result: "ClassificationResult"):
        try:
            async with AsyncClient() as http:
                url = f"{self._deps.config.get("API_BASE_URL")}/worker"
                response = await http.post(url, json=worker_result)
                response.raise_for_status()
        except:
            raise
