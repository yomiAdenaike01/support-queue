UPDATE input_sources
SET
    connection_value = $2::TEXT,
    status = $3::TEXT,
    source_type = $4::TEXT,
    name = $5::TEXT,
    config = COALESCE($6::JSONB, '{}'::JSONB),
    enabled = $7::BOOLEAN,
    team_id = team.id,
    team_pk = team.pk
FROM (SELECT $8::UUID AS team_id) input_data
LEFT JOIN team ON team.id = input_data.team_id
WHERE input_sources.id = $1::UUID
RETURNING
    input_sources.id,
    input_sources.pk,
    input_sources.connection_value,
    input_sources.team_pk,
    input_sources.team_id,
    input_sources.status,
    input_sources.source_type,
    input_sources.name,
    input_sources.config,
    input_sources.enabled,
    input_sources.created_at;
