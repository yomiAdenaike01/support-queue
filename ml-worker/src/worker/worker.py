import asyncio
import os
from logging import getLogger
from time import sleep

from redis import Redis

from .config import create_config
from .event_bus import EventBus
from .integrations import Integrations
from .knowledge_base import KnowledgeBase
from .pipeline import ClassificationPipeline, ResolutionPipeline
from .queue import (
    ClassificationQueue,
    QueueDependencies,
    ResolutionQueue,
)
from .repositories import TicketRepository, WorkerRepository

logger = getLogger("[worker]")


def init_redis_client(url: str) -> "Redis":
    redis_client: "Redis" = Redis.from_url(url, decode_responses=True)  # type: ignore
    result: bool = redis_client.ping()  # type: ignore
    if result is False:
        raise ValueError("Failed to connect to redis")
    logger.info("Connected to Redis at %s", url)
    return redis_client


class Worker:
    _classification_queue: "ClassificationQueue"
    _resolution_queue: "ResolutionQueue"

    def __init__(self):
        logger.info("[Application]: Bootstraping application...")
        worker_config = create_config()

        base_url = worker_config.get("BASE_URL")
        redis_url = os.environ.get("REDIS_ADDR", "redis://localhost:6379")

        health_check = self._health_check(base_url)
        if health_check is False:
            raise RuntimeError("Worker health check failed")

        redis_client = init_redis_client(redis_url)
        integrations = Integrations(cache=redis_client)
        knowledge_base = KnowledgeBase(base_url=f"{base_url}/knowledge")
        ticket_repository = TicketRepository(base_url=f"{base_url}/tickets")
        worker_repository = WorkerRepository(base_url=base_url)
        event_bus = EventBus(redis_client=redis_client)

        self._event_bus = event_bus

        resolution_queue_deps = QueueDependencies(
            ticket_repository=ticket_repository,
            worker_repository=worker_repository,
            config=worker_config,
            event_bus=event_bus,
            pipeline=ResolutionPipeline(integrations, knowledge_base=knowledge_base),
        )

        classification_queue_deps = QueueDependencies(
            ticket_repository=ticket_repository,
            worker_repository=worker_repository,
            config=worker_config,
            event_bus=event_bus,
            pipeline=ClassificationPipeline(
                integrations, knowledge_base=knowledge_base
            ),
        )

        self._classification_queue = ClassificationQueue(deps=classification_queue_deps)
        self._resolution_queue = ResolutionQueue(deps=resolution_queue_deps)
        logger.info("[Application]: Bootstraping complete!")

    async def start_workers(self):
        resolution_workers = self._resolution_queue.begin_workers()
        classification_workers = self._classification_queue.begin_workers()
        await asyncio.gather(*resolution_workers, *classification_workers)

    def _health_check(self, base_url: str):
        from httpx import Client

        max_attempts = 3
        default_backoff = 5
        with Client() as http:
            for attempt in range(max_attempts):
                try:
                    response = http.get(f"{base_url}/healthz")

                    if response.status_code == 200:
                        logger.info("[Application]: Health check passed!")
                        return True

                except Exception as e:
                    logger.error("[Application]: Health check failed!", e)

                if attempt < max_attempts - 1:
                    delay = default_backoff * (attempt + 1)
                    logger.info(
                        f"[Application]: Attempt {attempt + 1} failed. Retrying in {delay}s..."
                    )
                    sleep(delay)

        return False

    async def on_resolved_ticket_event(self):
        logger.info(
            "[Application:on_resolved_ticket_event]: Waiting for resolved ticket event..."
        )
        while True:
            evnt = await asyncio.to_thread(self._event_bus.listen_resolved_tickets)
            if evnt is None:
                continue
            resolved_ticket_id = evnt.data.get("ticket_id")
            if resolved_ticket_id is None:
                return
            await self._resolution_queue.add_to_queue(evnt)

    async def on_pending_resolved_tickets(self):
        logger.info(
            "[Application:on_resolved_ticket_event]: Waiting for pending ticket event..."
        )
        while True:
            evnt = await asyncio.to_thread(
                self._event_bus.get_incomplete_resolved_tickets
            )
            if evnt is None:
                break
            resolved_ticket_id = evnt.data.get("ticket_id")
            if resolved_ticket_id is None:
                break
            await self._resolution_queue.add_to_queue(evnt)
            await asyncio.sleep(1)

    async def on_pending_tickets(self):
        logger.info(
            "[Application:on_pending_tickets]: Waiting for pending ticket event..."
        )
        while True:
            evnt = await asyncio.to_thread(self._event_bus.listen_incomplete_tickets)
            if evnt is None:
                continue
            resolved_ticket_id = evnt.data.get("ticket_id")
            if resolved_ticket_id is None:
                continue
            await self._classification_queue.add_to_queue(evnt)
            await asyncio.sleep(1)

    async def on_new_ticket(self):
        logger.info("[Application:on_new_ticket]: Waiting for ticket event...")
        while True:
            evnt = await asyncio.to_thread(self._event_bus.listen_new_ticket)
            if evnt is None:
                continue
            ticket_id = evnt.data.get("ticket_id", None)
            if ticket_id is None:
                continue
            await self._classification_queue.add_to_queue(evnt)
