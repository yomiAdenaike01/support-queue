import json
from dataclasses import dataclass, field
from typing import TYPE_CHECKING, TypedDict, Optional
from enum import Enum

if TYPE_CHECKING:
    from ..ticket import Ticket
    from .resolution_pipeline import ResolvedTicketSummary


CATEGORIES = [
    "BILLING",
    "ACCOUNT_ACCESS",
    "TECHNICAL",
    "DELIVERY",
    "CANCELLATION",
    "GENERAL"
]


@dataclass
class ClassificationSystemPromptInput:
    sentiment: float
    categories: list[str] = field(default_factory=list(CATEGORIES))
    previous_cases: list["ResolvedTicketSummary"] = field(default_factory=list)
        
    def to_prompt_text(self):
        previous_cases_list = "" if len(self.previous_cases) == 0 else "\n".join([previous_case.to_json() for previous_case in self.previous_cases])
        categories_as_str = " | ".join(cat.upper() for cat in self.categories)
        return f"""You are a support ticket classifier. Reply ONLY with JSON, no other text.
        {f"Here is a list of similar cases: {previous_cases_list}" if len(previous_cases_list) > 0 else ""}
    Categories: {categories_as_str}
    Priority: URGENT|HIGH|MEDIUM|LOW
    Sentiment score: {self.sentiment} (range -1.0 to 1.0, negative is bad)

    Flag urgency_flag as true if messages contain: legal threats, fraud, chargeback, sarcasm, extreme distress, or account compromise. Set urgency_reason to one sentence explaining why if urgency_flag should be set to true.
    Ensure suggested response is in a professional and customer service tone.
    category, suggested_response and priority are required. Urgency reason is only required is urgency_flag should be set to true.
    {{"category":"","priority":"","suggested_response":"","urgency_flag":false,"urgency_reason":""}}
    """


class ClassificationJson(TypedDict):
    category: str
    priority: str
    suggested_response: str
    urgency_flag: bool
    urgency_reason: str


@dataclass
class Classification:
    category: str
    priority: str
    suggested_response: str
    urgency_flag: str
    urgency_reason: str




class Priority(str, Enum):
    LOW = "LOW"
    MEDIUM = "MEDIUM"
    HIGH = "HIGH"


@dataclass
class WorkerContext:
    ticket: "Ticket"
    average_sentiment: float = 0.0
    priority: Optional["Priority"] = None
    category: str = ""
    confidence_score: float = 0.0
    suggested_response: str = ""
    requires_urgency: bool = False
    urgency_reason: Optional[str] = None

    def to_json(self):
        return json.dumps({
            "average_sentiment": self.average_sentiment,
            "priority":self.priority,
            "suggested_response": self.suggested_response,
            "urgency_reason": self.urgency_reason,
            "requires_urgency": self.requires_urgency,
            })
@dataclass
class ResolvedTicketSummary:
    core_issue: str       
    steps_taken: str       
    resolution: str
    id: Optional[str] = None

    def to_json(self, ticket: "Ticket") -> str:
        return json.dumps({
            "core_issue": self.core_issue,
            "resolution": self.resolution,
            "steps_taken": self.steps_taken,
            "subject": ticket.subject,
            "suggested_response": ticket.suggested_response,
            "category": ticket.category
        })
