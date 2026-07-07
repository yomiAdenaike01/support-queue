from typing import Optional, TypedDict

from httpx import AsyncClient, HTTPStatusError

from .models import Filter


class RelatedKnowledgeResponse(TypedDict):
    id: str
    sourceId: str
    metadata: str
    sourceType: str
    content: dict[str, object]


class KnowledgeBase:
    _base_url: str

    def __init__(self, base_url: str):
        self._base_url = base_url

    async def search(self, filter: "Filter") -> list["RelatedKnowledgeResponse"]:
        try:
            async with AsyncClient() as http:
                url = f"{self._base_url}/search/"
                response = await http.post(url, json=filter.to_dict())
                response.raise_for_status()
                related_knowledge_list: list["RelatedKnowledgeResponse"] = (
                    response.json()
                )
                return related_knowledge_list
        except HTTPStatusError as status_error:
            if status_error.response.status_code != 404:
                raise
            return []

    async def insert_resolved_ticket(
        self, id: str, content: dict[str, object], embedding: list[list[float]]
    ) -> Optional[str]:
        async with AsyncClient() as http:
            url = f"{self._base_url}/"
            response = await http.post(
                url=url,
                json={
                    "source_id": id,
                    "content": content,
                    "source_event_type": "resolved_ticket",
                    "embedding": embedding,
                },
            )
            response.raise_for_status()
            return None
