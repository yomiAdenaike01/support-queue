from dataclasses import dataclass


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

    def to_json(self) -> str:
        import json
        return json.dumps({
            "id": self.id,
            "subject": self.subject,
            "customer_email": self.customer_email,
            "messages": [m.content for m in self.messages]
        })
