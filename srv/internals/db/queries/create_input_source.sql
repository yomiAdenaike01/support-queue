WITH input_data(
    connection_value, 
    status, 
    source_type, 
    name, 
    config
)
AS (
    VALUES (
        $1::TEXT,
        $2::TEXT,
        $3::TEXT,
        $4::TEXT,
        COALESCE($5::JSONB, NULL::JSONB)
    )
),
insert_input_source AS (
    INSERT INTO input_sources(
        connection_value, 
        status, 
        source_type, 
        name, 
        config, 
        team_id, 
        team_pk
    )
    SELECT
        input_data.*,
        id as team_id,
        pk as team_pk
    FROM 
        input_data
    LEFT JOIN team
    ON team.id = COALESCE($6::UUID, NULL::UUID)
    RETURNING
        id,
        pk,
        connection_value,
        team_pk,
        team_id,
        status,
        source_type,
        name,
        config,
        enabled,
        created_at

)
SELECT 
    i.*,
    t.department as department
FROM 
    insert_input_source i
LEFT JOIN team t
ON t.pk = i.team_pk


