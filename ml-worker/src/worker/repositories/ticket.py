import asyncio
from enum import Enum
from logging import getLogger
from typing import Any, Optional, TypedDict

from httpx import AsyncClient

from ..ticket import Ticket


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
    event_status: "TicketEventStatus"
    status: "TicketStatus"


class TicketStatus(str, Enum):
    PROCESSING = "PROCESSING"
    PENDING = "PENDING"
    PROCESSED = "PROCESSED"
    FAILED = "FAILED"
    RESOLVED = "RESOLVED"


class TicketEventStatus(str, Enum):
    CREATED = "TICKET_CREATED"
    PROCESSING = "TICKET_PROCESSING"
    CLASSIFIED = "TICKET_CLASSIFIED"
    PROCESSED = "TICKET_PROCESSED"
    FAILED = "TICKET_FAILED"
    REPROCESSED = "TICKET_REPROCESSED"
    RESOLVED = "TICKET_RESOLVED"


class TicketRepository:
    _base_url: str

    def __init__(self, base_url: str):
        self._base_url = base_url

    async def add_ticket_event(self, payload: "UpdateTicketEvent"):
        async with AsyncClient() as http:
            url = f"{self._base_url}/{payload.get('ticket_id')}/events"
            response = await http.post(
                url, json={**payload, "status": payload.get("event_status")}
            )
            response.raise_for_status()
            logger.info(f"Added event ticket_id={payload.get('ticket_id')}")

    async def set_status(self, payload: "UpdateTicketEvent"):
        async with AsyncClient() as http:
            url = f"{self._base_url}/{payload.get('ticket_id')}"
            response = await http.patch(url, json={"status": payload.get("status")})
            response.raise_for_status()
            logger.info(
                f"Updated ticket_id={payload.get('ticket_id')} status={payload.get('status')}"
            )

    async def set_ticket_status(self, payload: "UpdateTicketEvent"):
        await asyncio.gather(
            self.add_ticket_event(payload),
            self.set_status(payload),
            return_exceptions=True,
        )
        logger.info(f"Updated ticket_id={payload.get('ticket_id')}")

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

            return Ticket.from_json(ticket_json)
