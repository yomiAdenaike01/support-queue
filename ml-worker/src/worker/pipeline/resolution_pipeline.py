import json
from typing import TYPE_CHECKING
from dataclasses import dataclass
from logging import getLogger
from sentence_transformers import SentenceTransformer
if TYPE_CHECKING:
    from ..integrations import Integrations
    from ..ticket import Ticket
    from ..knowledge_base import KnowledgeBase
    from redis import Redis

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

@dataclass
class ResolvedTicketSummary:
    core_issue: str       
    steps_taken: str       
    resolution: str

    def to_text(self, ticket: "Ticket") -> list[str]:
        return f"""
            "core_issue": {self.core_issue}
            "resolution": {self.resolution}
            "steps_taken": {self.steps_taken}
            "subject": {ticket.subject}
            "suggested_response": {ticket.suggested_response}
            "category": {ticket.category}
        """

class ResolutionPipeline:
    _integrations: "Integrations"
    _knowledge_base: "KnowledgeBase"
    _transformer: "SentenceTransformer"

    def __init__(self, integrations: "Integrations", knowledge_base: "KnowledgeBase", cache: "Redis"):
        self._integrations = integrations
        self._knowledge_base = knowledge_base
        self._transformer = SentenceTransformer("all-MiniLM-L6-v2")
        self._cache = cache

    def _build_summarise_prompt(self, ticket: "Ticket") -> str:
        conversation = "\n".join([
            f"{msg.role.upper()}: {msg.content}"
            for msg in ticket.messages
        ])
        return f"""
            Ticket Subject: {ticket.subject}
            Conversation: {conversation}
        """

    async def _summarise(self, ticket: "Ticket") -> "ResolvedTicketSummary":
        summary = self._build_summarise_prompt(ticket)
        logger.info(f"Summarising ticket id={ticket.id}")
        response = await self._integrations.llm.prompt(system_prompt=RESOLUTION_SYSTEM_PROMPT, prompt=summary)
        return ResolvedTicketSummary(**json.loads(response))
    
    def _get_embedding(self, embedding_sentences: str):
            from hashlib import md5
            hash = md5("\n".join(embedding_sentences).encode("utf-8")).hexdigest()
            embedding = self._cache.get(hash)
            if embedding is not None:
                return json.loads(embedding)
            
            embedding = self._transformer.encode(sentences=embedding_sentences, show_progress_bar=True)
            embedding_as_list = embedding.tolist()
            self._cache.set(hash, json.dumps(embedding_as_list))
            return embedding_as_list
    
    async def run(self, ticket: "Ticket"):
        logger.info(f"Running pipeline on ticket={ticket}")
        summary: "ResolvedTicketSummary" = await self._summarise(ticket)
        embedding_sentences = summary.to_text(ticket)
        logger.info(f"Created embedding str id={ticket.id} embeddingstr={embedding_sentences}")
        embedding = self._get_embedding(embedding_sentences)
        logger.info(f"Successfully created embedding ticket id={ticket.id}")
        await self._knowledge_base.insert_resolved_ticket(id=ticket.id, content=embedding_sentences, embedding=embedding)

