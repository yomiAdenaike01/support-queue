from typing import Any, Optional, Callable, TypeVar, Protocol
from logging import getLogger
from ..utils import Timer

class JSONSerializable(Protocol):
    def to_json() -> str:
        ...

T = TypeVar("T", bound=JSONSerializable)

logger = getLogger(__name__)

class Stage:
    _timer: "Timer"
    _name: str
    _output: T
    _end: float
    _input: T

    def __init__(self, name: str, input: T, oncomplete: Optional[Callable[["Stage"], Any]] = None):
        self._timer = Timer()
        self._name = name
        self._input = input
        self._oncomplete = oncomplete

    def complete(self, output: T) -> float:
        self._output = output
        print('self._oncomplete ->, ', self._oncomplete)
        if self._oncomplete is not None:
            self._oncomplete(self)
        end = self._timer.elapsed()
        self._end = end
        return end
    def to_json(self):
        import json
        return json.dumps({
            "name":self._name,
            "output": self._output.to_json(),
            "input": self._input.to_json()
        })
        

class StageRegister:
    stages_by_name: dict[str, Any]

    def __init__(self):
        self.stages_by_name = {}

    def _oncomplete_stage(self, stage: "Stage"):
        logger.info(f"Completed stage:{stage._name} output:{stage._output}")
        from pathlib import Path
        import json
        try:    
            stages_path = Path(__file__).parents[3] / "stages.json"
            with open(stages_path,'a') as f:
                f.write(json.dumps({
                "stage_name": stage._name,
                "stage_input": stage._input,
                "stage_output": stage._output,
                "starttime": stage._start,
                "endtime": stage._end
                }))
        except:
            pass
        
    def start_new_stage(self, name: str, input: Any):
        logger.info(f"Starting new stage name:{name} input:{input}")
        stage = Stage(name=name, input=input, oncomplete=self._oncomplete_stage)
        self.stages_by_name[name] = stage
        return stage

    