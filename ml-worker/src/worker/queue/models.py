from typing import TYPE_CHECKING
from dataclasses import dataclass
if TYPE_CHECKING:
    from ..repositories import TicketRepository, WorkerRepository
    from ..event_bus import EventBus
    from ..pipeline import Pipeline

@dataclass(frozen=True)
class QueueDependencies:
    ticket_repository: "TicketRepository"
    worker_repository: "WorkerRepository"
    config: dict
    event_bus: "EventBus"
    pipeline: "Pipeline"