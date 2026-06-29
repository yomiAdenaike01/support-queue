import json
from hashlib import md5
from logging import getLogger

from redis import Redis
from sentence_transformers import SentenceTransformer
from sentence_transformers.cross_encoder import CrossEncoder

logger = getLogger("[encoders]")


class Encoders:
    _sentence_transformer: "SentenceTransformer"

    def __init__(self, cache: "Redis"):
        self._cache = cache
        self._transformer = SentenceTransformer("all-MiniLM-L6-v2")
        self._cross_encoder = CrossEncoder("cross-encoder/ms-marco-MiniLM-L-6-v2")

    def from_str_to_embedding(self, str: str) -> list[list[float]]:
        cache_key = f"embedding:{md5(str.encode('utf-8')).hexdigest()}"
        result = self._cache.get(cache_key)
        if result is not None:
            return list(json.loads(result))
        encoded_list: list[list[float]] = list(self._transformer.encode(str).tolist())
        try:
            self._cache.set(cache_key, json.dumps(encoded_list))
        except Exception as error:
            logger.error(f"Failed to create embedding error={str(error)}")
        return encoded_list

    def rank_pairings(self, keys: list[str], pairings: list[tuple[str, str]]):
        scores = self._cross_encoder.predict(pairings)
        return sorted(zip(keys, scores), key=lambda item: item[1], reverse=True)
