import json
from typing import TYPE_CHECKING
from logging import getLogger
from .models import ResolvedTicketSummary, Pipeline

if TYPE_CHECKING:
    from ..integrations import Integrations
    from ..ticket import Ticket
    from ..knowledge_base import KnowledgeBase

logger = getLogger('[resolution-pipeline]')


RESOLUTION_SYSTEM_PROMPT = """You are a support ticket resolution summarizer.

Your job is to summarize a resolved support ticket for a knowledge base.

Return valid JSON only. No markdown, no preamble, no explanation.

Rules:
- core_issue: one sentence describing the customer's original problem.
- steps_taken: one sentence describing what support or the customer tried.
- resolution: one sentence describing exactly how the issue was resolved.
- If a field cannot be determined, use "unknown".
- Do not invent facts that are not present in the conversation.
- Each field must be a string.
- Do not include extra keys.

Return this exact JSON shape:
{"core_issue":"","steps_taken":"","resolution":""}"""

class ResolutionPipeline(Pipeline):
    _integrations: "Integrations"
    _knowledge_base: "KnowledgeBase"

    def __init__(self, integrations: "Integrations", knowledge_base: "KnowledgeBase"):
        self._integrations = integrations
        self._knowledge_base = knowledge_base

    def _build_summarise_prompt(self, ticket: "Ticket") -> str:
        conversation = "\n".join([
            f"{msg.role.upper()}: {msg.content}"
            for msg in ticket.messages
        ])
        return f"""
            Subject: {ticket.subject}
            Conversation: {conversation}
        """

    async def _summarise(self, ticket: "Ticket") -> "ResolvedTicketSummary":
        summary = self._build_summarise_prompt(ticket)
        logger.info(f"Summarising ticket id={ticket.id}")
        response = await self._integrations.llm.prompt(system_prompt=RESOLUTION_SYSTEM_PROMPT, prompt=summary)
        return ResolvedTicketSummary(**json.loads(response))
    
    async def run(self, ticket: "Ticket"):
        logger.info(f"Running pipeline on ticket={ticket.id}")
        summary: "ResolvedTicketSummary" = await self._summarise(ticket)
        embedding_json = summary.to_json(ticket)
        logger.info(f"Created embedding str id={ticket.id} embeddingstr={embedding_json}")
        embedding = self._integrations.encoders.from_str_to_embedding(embedding_json)
        logger.info(f"Successfully created embedding ticket id={ticket.id}")
        await self._knowledge_base.insert_resolved_ticket(id=ticket.id, content=embedding_json, embedding=embedding)

