from redis import Redis
from .ticket import Ticket, Message
from pathlib import Path
from typing import Dict
import json
import logging
from datetime import datetime

logger = logging.getLogger(__name__)

class Store:
    _cache: Redis
    _url: str
    _tickets: Dict[str, Ticket] = {}

    def __init__(self, url: str):

        self._url = url

    def complete_ticket(self, ticket_id: str):
        self._cache.xack(ticket_id)

    def connect(self):
        self._cache = Redis(host=self._url)

    def load(self):
        base_path = Path(__file__).resolve().parents[2]

        test_file_path = Path(base_path) / "tickets.json"

        if test_file_path.exists():
            print("Successfully found tickets.json")

        contents = test_file_path.read_text(encoding="utf-8")
        tickets = json.loads(contents)
        for ticket in tickets:
            ticket_messages: list[Message] = [
                Message(content=msg["content"], created_at=datetime.now())
                for msg in ticket["messages"]
            ]
            ticket = Ticket(
                id=str(ticket["id"]),
                created_at=ticket["created_at"],
                messages=ticket_messages,
            )
            self._tickets[ticket.id] = ticket

    def dispose(self):
        self._cache.close()
        self._tickets = None
        self._cache = None

    def find_ticket_by_ref(self, id: str):
        # make sure that the ticket exists in stream
        # maybe fetch from sql to make sure that it exists
        return self._tickets.get(id, None)
