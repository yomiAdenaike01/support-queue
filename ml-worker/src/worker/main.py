from fastapi import FastAPI, APIRouter, Request, status
from contextlib import asynccontextmanager
from pydantic import BaseModel
from .store import Store
from .team import TeamRepository
from .pipeline import TicketPipeline
from .integrations import Integrations
from .logging import configure_logging
class TicketIdBody(BaseModel):
    ticket_id: str


@asynccontextmanager
async def lifespan(app: FastAPI):
    configure_logging()
    store = Store(url="redis://localhost:6379")
    integrations = Integrations()
    team_repository = TeamRepository()
    team_repository.new()
    pipeline = TicketPipeline(team_repository, integrations)
    store.load()
    store.connect()
    app.state.store = store
    app.state.pipeline = pipeline
    yield
    store.dispose()


app = FastAPI(lifespan=lifespan)
v1_router = APIRouter(prefix="/v1")


def get_store(req: Request) -> Store:
    return req.app.state.store


def get_pipeline(req: Request) -> TicketPipeline:
    return req.app.state.pipeline


@v1_router.get("/healthz")
def healthz():
    return "ok"


@v1_router.post("/ml")
async def pick_ticket(req: Request, body: TicketIdBody):
    store = get_store(req)
    pipeline = get_pipeline(req)
    ticket = store.find_ticket_by_ref(body.ticket_id)
    if ticket is None:
        return status.HTTP_404_NOT_FOUND

    await pipeline.run(ticket=ticket)
    return status.HTTP_200_OK


app.include_router(v1_router)

