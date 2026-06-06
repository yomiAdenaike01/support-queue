import json
import asyncio
from typing import TYPE_CHECKING, Dict, TypedDict
from httpx import AsyncClient, HTTPStatusError
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

logger = getLogger('[ticket-repository]');

class UpdateTicketEvent(TypedDict):
    ticket_id: str
    status: str
    payload: str
class TicketRepository:
    _base_url: str

    def __init__(self, base_url: str):
        self._base_url = base_url
    async def set_ticket_status(self, payload:UpdateTicketEvent ):
        try:
            async with AsyncClient() as http:
                response = await http.post(f"{self._base_url}/events",json=payload)
                logger.info(f"Successfully updated ticket_id={payload.ticket_id} payload={json.dumps(payload)}")
                response.raise_for_status()
        except HTTPStatusError as status_error:
            logger.info(f"Failed to update ticket status id={payload.ticket_id}")
            if status_error.response.status_code != 409:
                raise

    async def find_by_id(self, ticket_id: str) -> Ticket:
        async with AsyncClient() as http:
            url = f"{self._base_url}/{ticket_id}"
            logger.info(f"Fetching ticket url={url} id={ticket_id}")
            ticket_by_id_response = await http.get(url)
            ticket_by_id_response.raise_for_status()
            logger.info(f"Fetched ticket url={url} id={ticket_id} response={ticket_by_id_response}")
            ticket_json: FindTicketByIdResponse = ticket_by_id_response.json()
            messages_list = [Message(content=msg.get("content"),role=msg.get("role")) for msg in ticket_json.get("messages")]
            
            return Ticket(
                id=ticket_json.get("id"),
                subject=ticket_json.get("subject"),
                customer_email=ticket_json.get("customer_email"), 
                messages=messages_list,
                )
            
            


