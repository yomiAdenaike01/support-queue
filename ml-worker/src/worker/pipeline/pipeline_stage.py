from logging import getLogger
from typing import TYPE_CHECKING, Any, Callable, Optional, TypedDict

from ..utils import Timer

if TYPE_CHECKING:
    from ..pipeline import WorkerContext

logger = getLogger(__name__)


class StageDict(TypedDict):
    stage_name: str
    stage_input: str
    starttime: float
    endtime: float


class Stage:
    _timer: "Timer"
    _name: str
    _output: "WorkerContext"
    _end: float
    _input: "WorkerContext"

    def __init__(
        self,
        name: str,
        input: "WorkerContext",
        oncomplete: Optional[Callable[["Stage"], Any]] = None,
    ):
        self._timer = Timer()
        self._name = name
        self._input = input
        self._oncomplete = oncomplete

    def get_start(self):
        return self._timer.get_start()

    def complete(self, output: Any) -> float:
        self._output = output
        print("self._oncomplete ->, ", self._oncomplete)
        if self._oncomplete is not None:
            self._oncomplete(self)
        end = self._timer.elapsed()
        self._end = end
        return end

    def to_dict(self) -> StageDict:
        return {
            "stage_name": self._name,
            "stage_input": self._input.to_json(),
            "starttime": self._timer.get_start(),
            "endtime": self._timer.elapsed(),
        }

    def to_json(self):
        import json

        return json.dumps(self.to_dict())

    def get_name(self):
        return self._name

    def get_output(self):
        return self._output


class StageRegister:
    stages_by_name: dict[str, Any]

    def __init__(self):
        self.stages_by_name = {}

    def _oncomplete_stage(self, stage: "Stage"):
        logger.info(f"Completed stage:{stage.get_name()} output:{stage.get_output()}")
        from pathlib import Path

        stages_path = Path(__file__).parents[3] / "stages.json"
        with open(stages_path, "a") as f:
            f.write(stage.to_json())

    def start_new_stage(self, name: str, input: Any):
        logger.info(f"Starting new stage name:{name} input:{input}")
        stage = Stage(name=name, input=input, oncomplete=self._oncomplete_stage)
        self.stages_by_name[name] = stage
        return stage
