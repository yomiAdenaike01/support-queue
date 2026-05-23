import json
from pathlib import Path
from dataclasses import dataclass, field
from ..utils import read_json, JSONFilenames

TEAMS_BY_CATEGORY: dict[str, list[str]] = {
    "BILLING": ["billing-team", "logistics-team"],
    "TECHNICAL": ["tech-support"],
    "ACCOUNT_ACCESS": ["billing-team", "retention-team"],
    "GENERAL": ["customer-success"],
    "CANCELLATION": ["billing-team", "logistics-team"],
}


@dataclass
class Contact:
    email: str
    slack: str
    phone: str


@dataclass
class Team:
    teamId: int
    slug: str
    name: str
    contact: Contact
    hours: list[int] = field(default_factory=list)


class TeamRepository:
    _teams: dict[str, "Team"]

    def __init__(self):
        self._teams = {}

    def new(self):
        teams: dict = read_json(JSONFilenames.TEAM)

        for team_slug, team_info in teams.items():
            contact = Contact(
                email=team_info.get("email"),
                slack=team_info.get("slack"),
                phone=team_info.get("phone"),
            )
            self._teams[team_slug] = Team(
                name=team_info.get("name"),
                teamId=team_info.get("teamId"),
                hours=[0, 24],
                slug=team_slug,
                contact=contact,
            )

    async def gather_team_contacts_by_category(self, category: str) -> list["Team"]:
        potential_teams: list[str] = TEAMS_BY_CATEGORY.get(category)
        if potential_teams is None:
            return None
        print("resolved-teams ->", potential_teams)
        return [
            team
            for team_slug, team in self._teams.items()
            if team_slug in potential_teams
        ]
