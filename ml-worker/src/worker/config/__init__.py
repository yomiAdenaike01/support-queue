import os
from typing import TypedDict
from pathlib import Path
from dotenv import load_dotenv

class Config(TypedDict):
    BASE_URL: str
    MAX_ATTEMPTS: int

def create_config() -> Config: 
    env_path = Path(__file__).parents[2] / ".env"
    load_dotenv(env_path)
    return {
        "BASE_URL": os.environ.get("BASE_URL", ""),
        "MAX_ATTEMPTS": int(os.environ.get("MAX_ATTEMPTS", 3))
    }        

__all__ = ['create_config']