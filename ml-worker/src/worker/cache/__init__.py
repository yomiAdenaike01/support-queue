from redis import Redis

class CacheClient:
    _redis_client: "Redis"
    def __init__(self, redis_client: "Redis"):
        self._redis_client = redis_client
    def _hash_key(key: str):
        from hashlib import md5
        return md5(key.encode('utf-8')).hexdigest()
    def get(self, key: str):
        
        return self._redis_client.get(self._hash_key(key))
    
    def set(self, key: str, data):
        self._redis_client.set(self._hash_key(key), data)
        
__all__ = ["CacheClient"]