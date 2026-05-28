from __future__ import annotations
import json
import hashlib
import logging
from textblob import TextBlob
from dataclasses import dataclass
from enum import Enum
from typing import Optional, TypedDict, TYPE_CHECKING
from ..integrations import Integrations
from .pipeline_stage import StageRegister
from .pipeline_exception import PipelineException


if TYPE_CHECKING:
    from ..ticket import Ticket
    from redis import Redis

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


class Pipeline:
    _integrations: "Integrations"
    _stage_register: "StageRegister"
    _cache: "Redis"

    def __init__(self, integrations: "Integrations", cache: "Redis"):
        self._integrations = integrations
        self._stage_register = StageRegister()
        self._cache = cache

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
        # TODO: Append similar messages to prompt to ensure correct categorisation
        categories_as_str = " | ".join(cat.upper() for cat in categories)
        return f"""You are a support ticket classifier. Reply ONLY with JSON, no other text.

    Categories: {categories_as_str}
    Priority: URGENT|HIGH|MEDIUM|LOW
    Sentiment score: {sentiment} (range -1.0 to 1.0, negative is bad)

    Flag urgency_flag as true if messages contain: legal threats, fraud, chargeback, sarcasm, extreme distress, or account compromise. Set urgency_reason to one sentence explaining why if urgency_flag should be set to true.
    Ensure suggested response is in a professional and customer service tone.
    category, suggested_response and priority are required. Urgency reason is only required is urgency_flag should be set to true.
    {{"category":"","priority":"","suggested_response":"","urgency_flag":false,"urgency_reason":""}}
    """

    def _hash_prompt(self, prompt: str) -> str:
        return hashlib.md5(prompt.encode('utf-8')).hexdigest()
    
    def _store_classification_result(self, prompt: str,  classification: "ClassificationJson"):    
        key = self._hash_prompt(prompt)
        self._cache.set(key, json.dumps(classification), 1 * 60 * 10)

    def _get_stored_classification_result(self, prompt: str) -> "ClassificationJson" | None:
        entry = self._cache.get(self._hash_prompt(prompt))
        if entry is None:
            return None
        return ClassificationJson(**json.loads(entry))
        
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
        cached_result = self._get_stored_classification_result(system_prompt)

        if cached_result is not None:
             duration_ms = stage.complete(cached_result)
             logger.info(
            "Cache HIT Finished LLM ticket classification in %sms",
            duration_ms,
            extra={
                "ticket_id": ctx.ticket.id,
                "category": cached_result.get("category"),
                "priority": cached_result.get("priority"),
                "requires_urgency": cached_result.get("urgency_flag"),
                "duration_ms": duration_ms,
            }
             )
             return cached_result
        
        logger.info("Cache MISS starting manual classification")
        messages = ctx.ticket.messages
        last_msgs = " | ".join(msg.content for msg in messages[-3:])
        classification_result: ClassificationJson = json.loads(
            await self._integrations.llm.prompt(system_prompt, last_msgs)
        )
        # store classification result
        self._store_classification_result(system_prompt, classification_result)
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

    async def _run_classifications(self, ctx: WorkerContext, attempt = 0) -> WorkerContext:
        stage = self._stage_register.start_new_stage(name="classifcations", input=ctx)
        classification = await self._get_llm_classification(ctx)
        if len(classification.get('category')) == 0 and attempt + 1 <= 3:
            await self._run_classifications(ctx, attempt=attempt + 1)
            return
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



    async def run(self, ticket: Ticket) -> WorkerContext:
        stage = self._stage_register.start_new_stage("FULL_PIPELINE", input=ticket)
        logger.info("Starting ticket pipeline", extra={"ticket_id": ticket.id})
        ctx = WorkerContext(ticket=ticket)
        try:
            self._sentiment_per_message(ctx)
            await self._run_classifications(ctx)
        except Exception as e:
            elapsed = stage.complete(output=ctx)
            logger.exception(
                "Ticket pipeline failed after %sms",
                elapsed,
                extra={
                    "ticket_id": ticket.id,
                    "duration_ms": elapsed,
                },
            )
            raise PipelineException(ctx=ctx, message="Failed to run pipeline")
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
