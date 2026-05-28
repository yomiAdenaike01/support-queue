from .llm_integration import LLMIntegration, IntegrationOptions


class Integrations:
    llm: LLMIntegration

    def __init__(self):
        self.llm = LLMIntegration(
            options=IntegrationOptions(
                base_url="http://localhost:11434/v1", timeout=200.0
            )
        )
