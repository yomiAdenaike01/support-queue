import json
import logging
from dataclasses import dataclass
from typing import TYPE_CHECKING, Optional, TypedDict

from httpx import AsyncClient

from ..utils import Timer

logger = logging.getLogger(__name__)

if TYPE_CHECKING:
    from redis import Redis


class LLMMessage(TypedDict):
    role: str
    content: str


class LLMRequestBody(TypedDict):
    model: str
    messages: list[LLMMessage]


@dataclass
class IntegrationCreds:
    _client_secret: str
    _client_id: str
    _api_key: str

    def get_client_id(self):
        return self._client_id

    def get_api_key(self):
        return self._api_key

    def get_client_secret(self):
        return self._client_secret


@dataclass
class IntegrationOptions:
    base_url: str
    timeout: float


class LLMIntegration:
    _options: "IntegrationOptions"
    _cache: "Redis"

    def __init__(self, cache: "Redis", options: "IntegrationOptions"):
        self._options = options
        self._cache = cache

    def _hash(self, prompt_body: str):
        from hashlib import md5

        return md5(prompt_body.encode("utf-8")).hexdigest()

    def _get_cached_response(self, body_str: str) -> Optional[str]:
        hash: str = f"prompt:{self._hash(body_str)}"
        return str(self._cache.get(hash))

    async def prompt(self, system_prompt: str, prompt: str) -> str:

        url = f"{self._options.base_url}/chat/completions"
        body: LLMRequestBody = {
            "model": "llama3.2",
            "messages": [
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": prompt},
            ],
        }
        body_str = json.dumps(body)
        cached_response = self._get_cached_response(body_str)

        if cached_response is not None:
            return cached_response

        timer = Timer()
        logger.info(
            f"Sending LLM prompt - {body_str}",
            extra={
                "base_url": self._options.base_url,
                "model": body["model"],
                "system_prompt_length": len(system_prompt),
                "prompt_length": len(prompt),
            },
        )
        try:
            async with AsyncClient(timeout=self._options.timeout) as client:
                response = await client.post(url, json=body)
                duration_ms = timer.elapsed()
                logger.info(
                    "Received LLM response in %sms",
                    duration_ms,
                    extra={
                        "base_url": self._options.base_url,
                        "model": body["model"],
                        "status_code": response.status_code,
                        "duration_ms": duration_ms,
                    },
                )
                response.raise_for_status()
                data = response.json()
                content = data["choices"][0]["message"]["content"]

                self._cache.set(f"prompt:{self._hash(body_str)}", content)

                logger.info(f"LLM Response: {content}")
                logger.info(
                    "Parsed LLM response in %sms",
                    duration_ms,
                    extra={
                        "base_url": self._options.base_url,
                        "model": body["model"],
                        "response_length": len(content),
                        "duration_ms": duration_ms,
                    },
                )
                return content
        except Exception:
            duration_ms = timer.elapsed()
            logger.exception(
                "LLM prompt failed after %sms",
                duration_ms,
                extra={
                    "base_url": self._options.base_url,
                    "model": body["model"],
                    "duration_ms": duration_ms,
                },
            )
            raise
