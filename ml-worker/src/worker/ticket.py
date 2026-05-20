from dataclasses import dataclass
from datetime import datetime


@dataclass
class Message:
    content: str
    created_at: datetime


@dataclass
class Ticket:
    id: int
    created_at: int
    messages: list[Message]
