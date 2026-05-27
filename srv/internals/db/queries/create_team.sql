with ins_team as (
    insert into team(
        department,
        integrations
    ) values ($1, NULL::jsonb)
    returning id::text, pk, department, integrations
),
ins_team_member as (
    insert into team_members(
        team_id,
        team_pk,
        role,
        name,
        phone_number,
        email_address
    )
    select
        ins_team.id::uuid as team_id,
        ins_team.pk as team_pk,
        member.role,
        member.name,
        nullif(member.phone_number, ''),
        nullif(member.email_address, '')
    from ins_team
    cross join unnest(
        $2::text[],
        $3::varchar[],
        $4::text[],
        $5::text[]
    ) as member(role, name, phone_number, email_address)
    returning pk, id, email_address, phone_number, name, role, team_pk
)
select
    t.id as team_id,
    t.pk as team_pk,
    t.department,
    t.integrations::text as integrations,
    tm.pk as member_pk,
    tm.id as member_id,
    tm.email_address,
    tm.phone_number,
    tm.name,
    tm.role
from ins_team_member tm
left join ins_team t on tm.team_pk = t.pk


