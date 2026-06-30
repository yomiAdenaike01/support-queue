from dataclasses import dataclass
from typing import TYPE_CHECKING, Optional

if TYPE_CHECKING:
    from .repositories import FindTicketByIdResponse, MessageJSON


@dataclass
class Message:
    content: str
    role: str

    @classmethod
    def to_messages_list(cls, messages: list["MessageJSON"]) -> list["Message"]:
        return (
            [
                Message(content=msg.get("content"), role=msg.get("role"))
                for msg in messages
            ]
            if len(messages) > 0 and messages is not None  # type: ignore
            else []
        )


@dataclass
class Ticket:
    id: str
    subject: str
    customer_email: str
    messages: list[Message]
    category: Optional[str] = None
    suggested_response: Optional[str] = None

    @classmethod
    def from_json(cls, ticket_json: "FindTicketByIdResponse") -> Optional["Ticket"]:
        return (
            cls(
                id=ticket_json.get("id"),
                subject=ticket_json.get("subject"),
                customer_email=ticket_json.get("customerEmail"),
                category=ticket_json.get("category"),
                messages=Message.to_messages_list(ticket_json.get("messages")),
                suggested_response=ticket_json.get("suggestedResponse"),
            )
            if ticket_json.get("id") is not None  # type: ignore
            else None
        )

    def to_json(self) -> str:
        import json

        return json.dumps(
            {
                "id": self.id,
                "subject": self.subject,
                "customer_email": self.customer_email,
                "suggested_response": self.suggested_response,
                "category": self.category,
                "messages": [m.content for m in self.messages],
            }
        )

    def to_embed_str(self) -> Optional[str]:
        if len(self.messages) == 0:
            return None
        return f"""
            subject: {self.subject}
            message: {self.messages[0].content}
        """
