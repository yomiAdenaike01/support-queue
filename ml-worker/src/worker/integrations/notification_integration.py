from ..team import Contact


class NotificationIntegration:
    def __init__(self):
        pass

    async def send_many(self, contacts: list[Contact]):
        print("Sending many notifications to contacts -> ", contacts)
