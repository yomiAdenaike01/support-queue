from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from ..pipeline import WorkerContext


class PipelineException(Exception):
    __ctx: "WorkerContext"

    def __init__(self, ctx: "WorkerContext", message: str):
        super().__init__(message)
        self.__ctx = ctx
