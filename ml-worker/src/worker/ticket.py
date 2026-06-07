from dataclasses import dataclass
from typing import Optional

@dataclass
class Message:
    content: str
    role: str


@dataclass
class Ticket:
    id: str
    subject: str
    customer_email: str
    messages: list[Message]
    category: Optional[str] = None
    suggested_response: Optional[str] = None

    def to_json(self) -> str:
        import json
        return json.dumps({
            "id": self.id,
            "subject": self.subject,
            "customer_email": self.customer_email,
            "suggested_response": self.suggested_response,
            "category": self.category,
            "messages": [m.content for m in self.messages]
        })
    def to_embed_str(self) -> Optional[str]:
        if len(self.messages) == 0:
            return None
        return f"""
            subject: {self.subject}
            message: {self.messages[0].content}
        """