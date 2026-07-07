SELECT
    id,
    pk,
    source_id,
    source_type,
    metadata,
    content
FROM knowledge_base 
WHERE ($1::uuid IS NULL OR $1::uuid = source_id)
AND ($3::TEXT IS NULL OR $3::TEXT = source_type)
ORDER BY embedding <=> $2::vector
OFFSET $4
LIMIT $5