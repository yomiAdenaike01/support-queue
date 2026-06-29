SELECT 
    i.id as id,
    i.pk as pk,
    connection_value,
    team_pk,
    team_id,
    status,
    source_type,
    name,
    config,
    enabled,
    created_at,
    t.department as department
FROM input_sources i
LEFT JOIN team t ON t.pk = i.team_pk
WHERE ($1::TEXT IS NULL OR t.department = $1::TEXT)
AND ($2::TEXT IS NULL OR source_type = $2::TEXT)
AND ($3::UUID IS NULL OR i.id = $3::UUID)
LIMIT COALESCE($4::INT, 10)