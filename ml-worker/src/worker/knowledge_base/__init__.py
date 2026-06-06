from httpx import AsyncClient

class KnowledgeBase:
    _base_url: str
    def __init__(self,base_url: str):
        self._base_url = base_url

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

__all__ = ['KnowledgeBase']
