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

class TicketService:
    _base_url: str
    _tickets: Dict[str, "Ticket"] = {}
    _is_dry_run: bool

    def __init__(self, base_url: str, is_dry_run = False):
        self._base_url = base_url
        self._is_dry_run = is_dry_run
        self._tickets = {}
        if is_dry_run:
            self._load_from_file()
    
    def _load_from_file(self):
        tickets = read_json(JSONFilenames.TICKETS)
        for ticket in tickets:
            ticket_messages: list["Message"] = [
                Message(content=msg["content"], created_at=datetime.now())
                for msg in ticket["messages"]
            ]
            ticket = Ticket(
                id=str(ticket["id"]),
                created_at=ticket["created_at"],
                messages=ticket_messages,
            )
            self._tickets[ticket.id] = ticket

    async def find_ticket_by_id(self, ticket_id: str) -> Ticket:
        if self._is_dry_run:
            return self._tickets.get(ticket_id, None)
        async with AsyncClient() as http:
            try:
               ticket_by_id_response = await http.get(f"{self._base_url}/{ticket_id}")
               ticket_json: FindTicketByIdResponse = ticket_by_id_response.json()
               messages_list = [Message(content=msg.get("content"),role=msg.get("role")) for msg in ticket_json.get("messages")]
               
               return Ticket(id=ticket_json.get("id"),subject=ticket_json.get("subject"), customer_email=ticket_json.get("customer_email"), messages=messages_list)
            except Exception:
                raise
                
                


