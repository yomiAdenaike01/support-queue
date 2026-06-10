from redis import Redis
from typing import Any, TypedDict, Optional
import logging
logger = logging.getLogger(__name__)


class StreamEvent(TypedDict):
    data: dict
    id: str

class EventBus:
    _redis_client: "Redis"
    _url: str

    def __init__(self, redis_client: "Redis"):
        self._redis_client = redis_client

    def listen_resolved_tickets(self) -> Optional["StreamEvent"]:
        event = self._redis_client.xreadgroup(
            consumername='python-worker-1',
            groupname="TICKET_WORKERS",
            streams={"TICKETS:RESOLVED_STREAM": ">"},
            count=1,
            block=5000
        )
        if event is None:
            return
        for _, events in event:
            if events is None or len(events) == 0: 
                return 
            for event_id, event in events:
                return StreamEvent(id=event_id, data=event)
    def listen_new_ticket(self) -> Optional["StreamEvent"]:
        event = self._redis_client.xreadgroup(
            groupname="TICKET_WORKERS",
            consumername="python-worker-1",
            streams={"TICKETS_STREAM": ">"},
            count=1,
            block=5000,  # wait up to 5 seconds for a message
        )
        if event is None:
            return
        for _, messages in event:
            if messages is None or len(messages) == 0:
                return
            for message_id, message in messages:
                return StreamEvent(data=message,id=message_id)
            
    def ack_resolution(self, message_id: str):
        self._redis_client.xack("TICKETS:RESOLVED_STREAM","TICKET_WORKERS", message_id)
    def ack_classification(self, message_id: str):
        try:
            self._redis_client.xack("TICKETS_STREAM", "TICKET_WORKERS" ,message_id, )
        except Exception:
            return