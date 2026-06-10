import json
from base64 import standard_b64encode
from dataclasses import asdict, dataclass
from typing import Optional


@dataclass
class Filter:
    source_type: str
    embedding: Optional[list[list[float]]] = None
    limit: int = 20

    def _encode_embedding(self) -> Optional[str]:
        if self.embedding is None:
            return None
        byte_data = json.dumps(self.embedding).encode("utf-8")
        return standard_b64encode(byte_data).decode("utf-8")

    def to_qs(self) -> str:
        from urllib.parse import urlencode

        params = asdict(self)
        params["embedding"] = self._encode_embedding()
        return urlencode(
            {key: value for key, value in params.items() if value is not None}
        )
