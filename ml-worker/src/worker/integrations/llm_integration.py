import logging
from dataclasses import dataclass
from time import perf_counter

from httpx import AsyncClient

logger = logging.getLogger(__name__)


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
    _options: IntegrationOptions

    def __init__(self, options: IntegrationOptions):
        self._options = options

    async def prompt(self, system_prompt: str, prompt: str):
        url = f"{self._options.base_url}/chat/completions"
        body = {
            "model": "llama3.2",
            "messages": [
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": prompt},
            ],
        }
        started_at = perf_counter()
        logger.info(
            "Sending LLM prompt",
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
                duration_ms = round((perf_counter() - started_at) * 1000, 2)
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
                duration_ms = round((perf_counter() - started_at) * 1000, 2)
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
            duration_ms = round((perf_counter() - started_at) * 1000, 2)
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
