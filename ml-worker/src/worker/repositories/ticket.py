import json
from dataclasses import dataclass
from logging import getLogger
from typing import TYPE_CHECKING, Any, Optional, TypedDict

from httpx import AsyncClient, HTTPStatusError

from ..ticket import Message, Ticket

if TYPE_CHECKING:
    from ..ticket import Message, Ticket


class FindTicketByIdResponse(TypedDict):
    id: str
    customerEmail: str
    subject: str
    status: str
    messages: list["Message"]
    suggestedResponse: Optional[str]
    category: Optional[str]
    events: list[Any]


logger = getLogger("[ticket-repository]")


@dataclass
class UpdateTicketEvent:
    ticket_id: str
    status: str


class TicketRepository:
    _base_url: str

    def __init__(self, base_url: str):
        self._base_url = base_url

    async def set_ticket_status(self, payload: UpdateTicketEvent):
        try:
            async with AsyncClient() as http:
                response = await http.post(f"{self._base_url}/events", json=payload)
                logger.info(
                    f"Successfully updated ticket_id={payload.ticket_id} payload={json.dumps(payload)}"
                )
                response.raise_for_status()

        except HTTPStatusError as status_error:
            logger.info(f"Failed to update ticket status id={payload.ticket_id}")
            if status_error.response.status_code != 409:
                raise

    async def find_by_id(self, ticket_id: str) -> Optional["Ticket"]:
        async with AsyncClient() as http:
            try:
                url = f"{self._base_url}/{ticket_id}"
                logger.info(f"Fetching ticket url={url} id={ticket_id}")
                ticket_by_id_response = await http.get(url)
                ticket_by_id_response.raise_for_status()
                logger.info(
                    f"Fetched ticket url={url} id={ticket_id} response={ticket_by_id_response}"
                )
                ticket_json: FindTicketByIdResponse = ticket_by_id_response.json()
                messages_list = [
                    Message(content=msg.content, role=msg.role)
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
            except HTTPStatusError:
                return None
