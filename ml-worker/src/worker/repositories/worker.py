from dataclasses import dataclass
from typing import Optional

from httpx import AsyncClient, HTTPStatusError


@dataclass
class ClassificationResult:
    ticket_id: str
    suggested_response: str
    average_sentiment: float
    priority: str
    category: str
    requires_urgency: bool
    urgency_reason: Optional[str] = None


class WorkerRepository:
    _base_url: str

    def __init__(self, base_url: str) -> None:
        self._base_url = base_url

    async def complete_classification(self, worker_result: "ClassificationResult"):
        try:
            async with AsyncClient() as http:
                url = f"{self._base_url}/worker"
                response = await http.post(url, json=worker_result)
                response.raise_for_status()
        except HTTPStatusError:
            raise
        finally:
            return
