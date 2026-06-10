import json
from hashlib import md5

from redis import Redis
from sentence_transformers import SentenceTransformer
from sentence_transformers.cross_encoder import CrossEncoder


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
        encoded_list = self._transformer.encode(str).tolist()
        self._cache.set(cache_key, json.dumps(list))
        return encoded_list

    def rank_pairings(self, keys: list[str], pairings: list[tuple[str, str]]):
        scores = self._cross_encoder.predict(pairings)
        return sorted(zip(keys, scores), key=lambda item: item[1], reverse=True)
