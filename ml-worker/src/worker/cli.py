from .store import Store
from .integrations import Integrations
from .team import TeamRepository
from .pipeline import TicketPipeline
from .logging import configure_logging

async def cli():
    configure_logging()
    store = Store(url="redis://localhost:6379")
    integrations = Integrations()
    team_repository = TeamRepository()
    team_repository.new()

    pipeline = TicketPipeline(team_repository, integrations)
    store.load()
    store.connect()

    ticket = store.find_ticket_by_ref("3")
    ctx = await pipeline.run(ticket)
    store.ack(ctx.ticket.id)
    print("ticket pipeline complete", ctx)

if __name__ == "__main__":
    import asyncio

    asyncio.run(cli())
