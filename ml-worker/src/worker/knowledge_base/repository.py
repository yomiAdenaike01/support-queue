import json
from typing import TYPE_CHECKING, TypedDict
from httpx import AsyncClient
from .models import Filter

if TYPE_CHECKING:
    from ..pipeline import ResolvedTicketSummary

class RelatedKnowledgeResponse(TypedDict):
    id: str
    sourceId: str
    metadata: str
    sourceType: str
    content: str


class KnowledgeBase:
    _base_url: str
    def __init__(self,base_url: str):
        self._base_url = base_url
    
    async def search(self, filter: "Filter") -> list["RelatedKnowledgeResponse"]:
        try:
            async with AsyncClient() as http:
                url = f"{self._base_url}/knowledge?{filter.to_qs()}"
                response = await http.get(url)
                response.raise_for_status()
                related_knowledge_list: list["RelatedKnowledgeResponse"] = response.json()
                return related_knowledge_list
        except:
            raise

    async def insert_resolved_ticket(self,id: str,content: str,  embedding: list) -> str:
        async with AsyncClient() as http:
            url = f"{self._base_url}/knowledge/"
            response = await http.post(url=url,json={
                "source_id":id,
                "content": content,
                "source_event_type":"resolved_ticket",
                "embedding": embedding
            })
            response.raise_for_status()
            return
