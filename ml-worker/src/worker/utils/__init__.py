from time import perf_counter
from pathlib import Path
import json
from enum import Enum

__ROOT_PATH = Path(__file__).parents[4].resolve()

class Timer:
    def __init__(self):
        self._start = perf_counter()
    def elapsed(self) -> float:
        return round((perf_counter() - self._start) * 1000, 2) 

class JSONFilenames(str, Enum):
    TICKETS = 'tickets'
    TEAM = 'teams'
    STAGES = "stages"

def read_json(file_name: JSONFilenames):
    file_path = __ROOT_PATH / f"{file_name}.json"
    print(file_path)
    try:
        if file_path.exists() is False:
            raise FileNotFoundError
        return json.loads(file_path.read_text(encoding="utf-8"))
    except Exception:
        raise

__all__ = ['Timer','read_json',"JSONFilenames"]