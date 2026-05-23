from dataclasses import dataclass
from datetime import datetime


@dataclass
class Message:
    content: str
    created_at: datetime
    role: str


@dataclass
class Ticket:
    id: int
    subject: str
    created_at: int
    customer_email: str
    messages: list[Message]
