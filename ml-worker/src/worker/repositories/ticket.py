import json
from logging import getLogger
from typing import TYPE_CHECKING, Any, Optional, TypedDict

from httpx import AsyncClient

from ..ticket import Message, Ticket

if TYPE_CHECKING:
    from ..ticket import Message, Ticket


class MessageJSON(TypedDict):
    content: str
    role: str
    id: str


class FindTicketByIdResponse(TypedDict):
    id: str
    customerEmail: str
    subject: str
    status: str
    messages: list["MessageJSON"]
    suggestedResponse: Optional[str]
    category: Optional[str]
    events: list[Any]


logger = getLogger("[ticket-repository]")


class UpdateTicketEvent(TypedDict):
    ticket_id: str
    status: str


class TicketRepository:
    _base_url: str

    def __init__(self, base_url: str):
        self._base_url = base_url

    async def set_ticket_status(self, payload: "UpdateTicketEvent"):
        async with AsyncClient() as http:
            url = f"{self._base_url}/{payload.get('ticket_id')}/events"
            response = await http.post(url, json=payload)
            response.raise_for_status()
            logger.info(
                f"Updated ticket_id={payload.get('ticket_id')} payload={json.dumps(payload)}"
            )

    async def find_by_id(self, ticket_id: str) -> Optional["Ticket"]:
        async with AsyncClient() as http:
            url = f"{self._base_url}/{ticket_id}"
            logger.info(f"Fetching ticket url={url} id={ticket_id}")
            ticket_by_id_response = await http.get(url)
            ticket_by_id_response.raise_for_status()
            logger.info(
                f"Fetched ticket url={url} id={ticket_id} response={ticket_by_id_response}"
            )
            ticket_json: "FindTicketByIdResponse" = ticket_by_id_response.json()

            messages_list: list["Message"] = [
                Message(content=msg.get("content"), role=msg.get("role"))
                for msg in ticket_json.get("messages")
            ]

            return Ticket(
                id=ticket_json.get("id"),
                subject=ticket_json.get("subject"),
                customer_email=ticket_json.get("customerEmail"),
                category=ticket_json.get("category"),
                messages=messages_list,
                suggested_response=ticket_json.get("suggestedResponse"),
            )
