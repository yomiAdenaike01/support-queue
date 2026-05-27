WITH selected_team_members AS (
    SELECT 
        email_address,
        phone_number,
        id,
        pk,
        team_id,
        team_pk,
        name,
        role,
        integrations
    FROM team_members
    WHERE ($1::uuid IS NULL OR team_id = $1::uuid)
)
SELECT 
    t.id as team_id,
    t.pk as team_pk,
    department,
    t.integrations::text as integrations,
    tm.id as member_id,
    tm.pk as member_pk,
    tm.name AS name,
    tm.role AS role,
    tm.phone_number as phone_number,
    tm.email_address as email_address
FROM team t
LEFT JOIN selected_team_members tm 
ON t.pk = tm.team_pk
WHERE ($1::uuid IS NULL OR t.id = $1::uuid)
AND (
    $2::text[] IS NULL
    OR cardinality($2::text[]) = 0
    OR t.department = ANY($2::text[])
)

