from __future__ import annotations

import hashlib
import json
import logging
from typing import TYPE_CHECKING, Optional

from textblob import TextBlob

from ..integrations import Integrations
from ..knowledge_base import SearchFilter
from .models import (
    ClassificationJson,
    ClassificationSystemPromptInput,
    Pipeline,
    Priority,
    ResolvedTicketSummary,
    WorkerContext,
)
from .pipeline_exception import PipelineException
from .pipeline_stage import StageRegister

if TYPE_CHECKING:
    from ..knowledge_base import KnowledgeBase
    from ..ticket import Ticket

logger = logging.getLogger(__name__)


# run stage class
class ClassificationPipeline(Pipeline):
    _integrations: "Integrations"
    _stage_register: "StageRegister"
    _knowledge_base: "KnowledgeBase"

    def __init__(self, integrations: "Integrations", knowledge_base: "KnowledgeBase"):
        self._integrations = integrations
        self._stage_register = StageRegister()

        self._knowledge_base = knowledge_base

    def _sentiment_per_message(self, ctx: "WorkerContext") -> None:
        stage = self._stage_register.start_new_stage(
            name="sentiment-analysis", input=ctx
        )
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

    def _hash_prompt(self, prompt: str) -> str:
        return hashlib.md5(prompt.encode("utf-8")).hexdigest()

    async def _fetch_similar_tickets(
        self, embedding: list[list[float]]
    ) -> list["ResolvedTicketSummary"]:
        related_knowledge_list = await self._knowledge_base.search(
            SearchFilter(source_type="resolved_ticket", embedding=embedding)
        )
        return [
            ResolvedTicketSummary(**json.dumps(item.content), id=item.get("sourceId"))
            for item in related_knowledge_list
        ]

    async def _get_similar_cases(
        self, ctx_ticket: "Ticket", ticket_embed_str: str
    ) -> list["ResolvedTicketSummary"]:
        embedding: list[list[float]] = (
            self._integrations.encoders.from_str_to_embedding(ticket_embed_str)
        )
        related_resolved_tickets = await self._fetch_similar_tickets(embedding)
        if len(related_resolved_tickets) > 0:
            ticket_map: dict[str, "ResolvedTicketSummary"] = {
                ticket.to_json(ctx_ticket): ticket
                for ticket in related_resolved_tickets
            }
            pairings = [
                (ticket_embed_str, ticket.to_json(ctx_ticket))
                for ticket in related_resolved_tickets
            ]
            ranked_pairings = self._integrations.encoders.rank_pairings(
                keys=[ticket_json for _, ticket_json in pairings], pairings=pairings
            )

            filtered_list: list["ResolvedTicketSummary"] = [
                ticket
                for ticket_json, _ in ranked_pairings[:3]
                if (ticket := ticket_map.get(ticket_json, None)) is not None
            ]
            return filtered_list
        return []

    async def _get_llm_classification(self, ctx: "WorkerContext") -> ClassificationJson:
        stage = self._stage_register.start_new_stage(
            name="llm-classification", input=ctx
        )
        logger.info(
            "Starting LLM ticket classification",
            extra={
                "ticket_id": ctx.ticket.id,
                "message_count": len(ctx.ticket.messages),
            },
        )
        ticket_embed = ctx.ticket.to_embed_str()
        resolved_prev_cases = []
        if ticket_embed is not None:
            resolved_prev_cases = await self._get_similar_cases(
                ctx.ticket, ticket_embed
            )
        system_prompt = ClassificationSystemPromptInput(
            sentiment=ctx.average_sentiment,
            previous_cases=resolved_prev_cases,
        )

        logger.info("Cache MISS starting manual classification")
        classification_result: "ClassificationJson" = json.loads(
            await self._integrations.llm.prompt(
                system_prompt=system_prompt.to_prompt_text(ctx.ticket),
                prompt=str(ticket_embed),
            )
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

    async def _run_classifications(
        self, ctx: "WorkerContext", attempt: int = 1
    ) -> Optional["WorkerContext"]:
        stage = self._stage_register.start_new_stage(name="classifcations", input=ctx)
        classification = await self._get_llm_classification(ctx)
        if len(classification.get("category")) == 0 and attempt + 1 <= 3:
            await self._run_classifications(ctx, attempt=attempt + 1)
            return
        ctx.category = classification.get("category")
        ctx.priority = Priority(classification.get("priority"))
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

    async def run(self, ticket: "Ticket") -> "WorkerContext":
        stage = self._stage_register.start_new_stage("FULL_PIPELINE", input=ticket)
        logger.info("Starting ticket pipeline", extra={"ticket_id": ticket.id})
        ctx = WorkerContext(ticket=ticket)
        try:
            self._sentiment_per_message(ctx)
            await self._run_classifications(ctx)
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
