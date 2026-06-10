from dataclasses import dataclass
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from ..config import Config
    from ..event_bus import EventBus
    from ..pipeline import Pipeline
    from ..repositories import TicketRepository, WorkerRepository


@dataclass(frozen=True)
class QueueDependencies:
    ticket_repository: "TicketRepository"
    worker_repository: "WorkerRepository"
    config: "Config"
    event_bus: "EventBus"
    pipeline: "Pipeline"
