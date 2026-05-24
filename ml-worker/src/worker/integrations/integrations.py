from .llm_integration import LLMIntegration, IntegrationOptions
from .notification_integration import NotificationIntegration


class Integrations:
    llm: LLMIntegration
    notifications: NotificationIntegration

    def __init__(self):
        self.llm = LLMIntegration(
            options=IntegrationOptions(
                base_url="http://localhost:11434/v1", timeout=200.0
            )
        )
        self.notifications = NotificationIntegration()
