from .llm_integration import LLMIntegration, IntegrationOptions
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from redis import Redis

class Integrations:
    llm: LLMIntegration

    def __init__(self, cache:"Redis"):
        self.llm = LLMIntegration(
            options=IntegrationOptions(
                base_url="http://localhost:11434/v1", timeout=200.0
            ),
            cache=cache
        )
