from typing import TYPE_CHECKING
from .llm_integration import LLMIntegration, IntegrationOptions
from .encoders import Encoders

if TYPE_CHECKING:
    from redis import Redis
class Integrations:
    llm: "LLMIntegration"
    encoders: "Encoders"
    
    def __init__(self, cache:"Redis"):
        self.llm = LLMIntegration(
            options=IntegrationOptions(
                base_url="http://localhost:11434/v1", timeout=200.0
            ),
            cache=cache
        )
        self.encoders = Encoders(cache=cache)
