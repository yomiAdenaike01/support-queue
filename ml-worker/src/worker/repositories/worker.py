from dataclasses import asdict, dataclass
from logging import getLogger
from typing import Optional

from httpx import AsyncClient

logger = getLogger("[classifcation-result]")


@dataclass
class ClassificationResult:
    ticket_id: str
    suggested_response: str
    average_sentiment: float
    category: str
    requires_urgency: bool
    priority: Optional[str] = None
    urgency_reason: Optional[str] = None


class WorkerRepository:
    _base_url: str

    def __init__(self, base_url: str) -> None:
        self._base_url = base_url

    async def complete_classification(self, worker_result: "ClassificationResult"):
        try:
            async with AsyncClient() as http:
                url = f"{self._base_url}/worker"
                response = await http.post(url, json=asdict(worker_result))
                response.raise_for_status()
        except Exception as error:
            logger.exception("Failed to complete classification reason=%s", str(error))
            raise error
