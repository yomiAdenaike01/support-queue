from redis import Redis
from pathlib import Path
from typing import Dict, TYPE_CHECKING
import json
import logging
from datetime import datetime
from .ticket import Message, Ticket

logger = logging.getLogger(__name__)


class EventBus:
    _cache: "Redis"
    _url: str

    def __init__(self, url: str):
        self._url = url

    def await_new_event(self):
        return self._cache.xreadgroup(
            groupname="TICKET_WORKERS",
            consumername="python-worker-1",
            streams={"TICKETS_STREAM": ">"},
            count=1,
            block=5000,  # wait up to 5 seconds for a message
        )

    def register_completion(self, ticket_id: str):
        try:
            self._cache.xack(ticket_id)
        except Exception:
            return

    def connect(self):
        self._cache = Redis.from_url(self._url, decode_responses=True)
        result = self._cache.ping()
        if result is False:
            raise ValueError("Failed to connect to redis")
        logger.info("Connected to Redis at %s", self._url)

    def dispose(self):
        self._cache.close()
        self._cache = None
