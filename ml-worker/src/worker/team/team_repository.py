from dataclasses import dataclass
from httpx import AsyncClient
from ..utils import read_json, JSONFilenames
from typing import TypedDict

@dataclass
class TeamMember:
    id: str
    email: str
    integrations: str
    phone: str


@dataclass
class Team:
    id: str
    department: str
    members: list['TeamMember']


class TeamMemberResponse(TypedDict):
    id: str
    email_address: str
    integrations: str
    phone_number: str
class TeamsResponse(TypedDict, Team):
    id: str
    department: str
    members: list['TeamMemberResponse']

TEAMS_BY_CATEGORY: dict[str, list[str]] = {
    "BILLING": ["billing", "logistics"],
    "TECHNICAL": ["tech-support"],
    "ACCOUNT_ACCESS": ["billing", "retention"],
    "GENERAL": ["customer-success"],
    "CANCELLATION": ["billing", "logistics"],
}

class TeamRepository:
    _base_url: str
    def __init__(self, base_url: str):
        self._base_url = base_url

    async def get_teams_by_category(self, category: str) -> list["Team"]:
        async with AsyncClient() as http:
            try:
                url = f'{self._base_url}/team/search?departments={','.join(TEAMS_BY_CATEGORY.get(category))}'
                response = await http.get(url=url)
                response.raise_for_status()
                raw_teams: list['TeamsResponse'] = response.json()
                teams = [Team(id=t.get("id"),department=t.get('department'), members=[TeamMember(email=m.get("email"), id=m.get("id"), phone=m.get("phone_number"), integrations=m.get("integrations")) for m in t.members]) for t in raw_teams]
                return teams
            except:
                raise
