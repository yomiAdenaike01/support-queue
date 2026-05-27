import json
import asyncio
from typing import TYPE_CHECKING, Dict, TypedDict
from httpx import AsyncClient
from pathlib import Path
from logging import getLogger
from datetime import datetime
from ..ticket import Ticket, Message
from ..utils import read_json, JSONFilenames

if TYPE_CHECKING:
    from ..ticket import Ticket, Message

class FindTicketByIdResponse(TypedDict):
    id: str
    customer_email: str
    subject: str
    status: str
    messages: list["Message"]

logger = getLogger(__name__);

class TicketRepository:
    _base_url: str

    def __init__(self, base_url: str):
        self._base_url = base_url

    async def find_by_id(self, ticket_id: str) -> Ticket:
        async with AsyncClient() as http:
            ticket_by_id_response = await http.get(f"{self._base_url}/{ticket_id}")
            ticket_by_id_response.raise_for_status()
            ticket_json: FindTicketByIdResponse = ticket_by_id_response.json()
            messages_list = [Message(content=msg.get("content"),role=msg.get("role")) for msg in ticket_json.get("messages")]
            
            return Ticket(id=ticket_json.get("id"),subject=ticket_json.get("subject"), customer_email=ticket_json.get("customer_email"), messages=messages_list)
            
            


