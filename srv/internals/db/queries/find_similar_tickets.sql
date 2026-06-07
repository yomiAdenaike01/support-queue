SELECT 
    id,
    pk,
    source_type,
    metadata,
    content
FROM knowledge_base 
WHERE ($1 IS NULL OR $1::uuid = source_id)
AND ($2 IS NULL OR $2::TEXT = source_type)
ORDER BY embedding <=> $2::vector
OFFSET $3
LIMIT $4