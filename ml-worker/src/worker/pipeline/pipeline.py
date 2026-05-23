from __future__ import annotations
import json
import asyncio
import logging
from textblob import TextBlob
from dataclasses import dataclass, field
from enum import Enum
from typing import Optional, TypedDict, TYPE_CHECKING
from ..integrations import Integrations
from .pipeline_stage import Stage, StageRegister


if TYPE_CHECKING:
    from ..team import TeamRepository, Team
    from ..ticket import Ticket

logger = logging.getLogger(__name__)

# run stage class


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


CATEGORIES = [
    "BILLING",
    "ACCOUNT_ACCESS",
    "TECHNICAL",
    "DELIVERY",
    "CANCELLATION",
    "GENERAL",
]


class Priority(str, Enum):
    LOW = "LOW"
    MEDIUM = "MEDIUM"
    HIGH = "HIGH"


@dataclass
class WorkerContext:
    ticket: Ticket
    average_sentiment: float = 0.0
    priority: Optional["Priority"] = None
    category: str = ""
    confidence_score: float = 0.0
    suggested_teams: list["Team"] = field(default_factory=list)
    suggested_response: str = ""
    requires_urgency: bool = False
    urgency_reason: Optional[str] = None

    def to_json(self):
        return json.dumps({
            "average_sentiment": self.average_sentiment,
            "priority":self.priority,
            "suggested_teams": self.suggested_teams,
            "suggested_response": self.suggested_response,
            "urgency_reason": self.urgency_reason,
            "requires_urgency": self.requires_urgency,
            })


class Pipeline:
    _team_repository: "TeamRepository"
    _integrations: "Integrations"
    _stage_register: "StageRegister"

    def __init__(self, team_repository: "TeamRepository", integrations: "Integrations"):
        self._team_repository = team_repository
        self._integrations = integrations
        self._stage_register = StageRegister()

    def _sentiment_per_message(self, ctx: WorkerContext) -> WorkerContext:
        stage = self._stage_register.start_new_stage(name="sentiment-analysis", input=ctx)
        average_sentiment = sum(
            [TextBlob(msg.content).sentiment.polarity for msg in ctx.ticket.messages]
        ) / len(ctx.ticket.messages)
        ctx.average_sentiment = average_sentiment
        duration_ms = stage.complete(average_sentiment)
        logger.info(
            "Calculated ticket sentiment in %sms",
            duration_ms,
            extra={
                "ticket_id": ctx.ticket.id,
                "message_count": len(ctx.ticket.messages),
                "average_sentiment": average_sentiment,
                "duration_ms": duration_ms,
            },
        )

    def _get_prompt(self, sentiment: float, categories: list[str]) -> str:
        categories_as_str = " | ".join(cat.upper() for cat in categories)
        return f"""You are a support ticket classifier. Reply ONLY with JSON, no other text.

    Categories: {categories_as_str}
    Priority: URGENT|HIGH|MEDIUM|LOW
    Sentiment score: {sentiment} (range -1.0 to 1.0, negative is bad)

    Flag urgency_flag as true if messages contain: legal threats, fraud, chargeback, sarcasm, extreme distress, or account compromise. Set urgency_reason to one sentence explaining why if urgency_flag should be set to true.
    Ensure suggested response is in a professional and customer service tone, is it a required field.
    {{"category":"","priority":"","suggested_response":"","urgency_flag":false,"urgency_reason":""}}
    """

    async def _get_llm_classification(self, ctx: WorkerContext) -> ClassificationJson:
        stage = self._stage_register.start_new_stage(name="llm-classification", input=ctx)
        logger.info(
            "Starting LLM ticket classification",
            extra={
                "ticket_id": ctx.ticket.id,
                "message_count": len(ctx.ticket.messages),
            },
        )
        system_prompt = self._get_prompt(ctx.average_sentiment, categories=CATEGORIES)
        messages = ctx.ticket.messages
        last_msgs = " | ".join(msg.content for msg in messages[-3:])
        classification_result: ClassificationJson = json.loads(
            await self._integrations.llm.prompt(system_prompt, last_msgs)
        )
        duration_ms = stage.complete(classification_result)
        logger.info(
            "Finished LLM ticket classification in %sms",
            duration_ms,
            extra={
                "ticket_id": ctx.ticket.id,
                "category": classification_result.get("category"),
                "priority": classification_result.get("priority"),
                "requires_urgency": classification_result.get("urgency_flag"),
                "duration_ms": duration_ms,
            },
        )

        return classification_result

    async def _run_classifications(self, ctx: WorkerContext) -> WorkerContext:
        stage = self._stage_register.start_new_stage(name="classifcations", input=ctx)
        classification = await self._get_llm_classification(ctx)

        ctx.category = classification.get("category")
        ctx.priority = classification.get("priority")
        ctx.suggested_response = classification.get("suggested_response")
        ctx.requires_urgency = classification.get("urgency_flag")
        ctx.urgency_reason = classification.get("urgency_reason")
        elapsed = stage.complete(ctx)
        logger.info(
            "Applied ticket classification in %sms",
            elapsed,
            extra={
                "ticket_id": ctx.ticket.id,
                "category": ctx.category,
                "priority": ctx.priority,
                "requires_urgency": ctx.requires_urgency,
                "duration_ms": elapsed,
            },
        )

    async def _resolve_team(self, ctx: WorkerContext):
        stage = self._stage_register.start_new_stage(name="team-resolution",input=ctx)
        if ctx.category == "" or ctx.category is None:
            raise ValueError("No category found on ctx")
        suggested_teams = await self._team_repository.gather_team_contacts_by_category(
            ctx.category
        )
        ctx.suggested_teams = suggested_teams
        elapsed = stage.complete(ctx)
        logger.info(
            "Resolved suggested teams in %sms",
            elapsed,
            extra={
                "ticket_id": ctx.ticket.id,
                "category": ctx.category,
                "team_count": len(suggested_teams),
                "duration_ms": elapsed,
            },
        )

    async def _send_team_notifications(self, ctx: WorkerContext):
        stage = self._stage_register.start_new_stage(name="send-team-notifications", input=ctx)

        if ctx.suggested_response is None or len(ctx.suggested_teams) < 1:
            elapsed = stage.complete(ctx)
            logger.info(
                "Skipped team notifications in %sms",
                elapsed,
                extra={
                    "ticket_id": ctx.ticket.id,
                    "team_count": len(ctx.suggested_teams),
                    "duration_ms": elapsed,
                },
            )
            return
        contacts = [team.contact for team in ctx.suggested_teams]
        
        task = asyncio.create_task(self._integrations.notifications.send_many(contacts))
        task.add_done_callback(self._log_notification_task_result)
        elapsed = stage.complete(ctx)

        logger.info(
            "Sent team notifications in %sms",
            elapsed,
            extra={
                "ticket_id": ctx.ticket.id,
                "contact_count": len(contacts),
                "duration_ms": elapsed,
            },
        )

    def _log_notification_task_result(self, task: asyncio.Task):
        if task.cancelled():
            logger.info("Team notification task was cancelled")
            return

        error = task.exception()
        if error is not None:
            logger.error(
                "Team notification task failed",
                exc_info=(type(error), error, error.__traceback__),
            )

    async def run(self, ticket: Ticket) -> WorkerContext:
        
        stage = self._stage_register.start_new_stage("FULL_PIPELINE", input=ticket)
        logger.info("Starting ticket pipeline", extra={"ticket_id": ticket.id})
        ctx = WorkerContext(ticket=ticket)
        try:
            self._sentiment_per_message(ctx)
            await self._run_classifications(ctx)
            await self._resolve_team(ctx)
        except Exception:
            elapsed = stage.complete(output=ctx)
            logger.exception(
                "Ticket pipeline failed after %sms",
                elapsed,
                extra={
                    "ticket_id": ticket.id,
                    "duration_ms": elapsed,
                },
            )
            raise

        task = asyncio.create_task(self._send_team_notifications(ctx))
        task.add_done_callback(self._log_notification_task_result)
        elapsed = stage.complete(ctx)
        logger.info(
            "Finished ticket pipeline in %sms",
            elapsed,
            extra={
                "ticket_id": ticket.id,
                "category": ctx.category,
                "priority": ctx.priority,
                "requires_urgency": ctx.requires_urgency,
                "duration_ms": elapsed,
            },
        )
        return ctx
