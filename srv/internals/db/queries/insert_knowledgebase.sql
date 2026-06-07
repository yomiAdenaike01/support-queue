WITH ins_row AS (
INSERT INTO knowledge_base(source_id, source_type, metadata, content, embedding)
    SELECT 
        $1::uuid, 
        $2::TEXT, 
        COALESCE($3::jsonb, NULL::jsonb), 
        COALESCE($4::jsonb, NULL::jsonb), 
        $5::vector
    WHERE NOT EXISTS(
        SELECT 
            1 
        FROM 
            knowledge_base
        WHERE 
            source_id = $1::uuid
        AND 
            source_type = $2::TEXT
    )
RETURNING 
    id, 
    pk, 
    source_type, 
    metadata, 
    source_id, 
    content
)
SELECT 
    id, 
    pk, 
    source_type, 
    metadata, 
    source_id, 
    content 
FROM 
    ins_row
UNION ALL 
SELECT 
    id, 
    pk, 
    source_type, 
    metadata, 
    source_id, 
    content
FROM 
    knowledge_base 
WHERE 
    source_id = $1::uuid AND source_type = $2::TEXT
