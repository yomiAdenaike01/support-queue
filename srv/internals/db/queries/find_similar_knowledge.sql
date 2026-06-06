SELECT
    id,
    pk,
    source_id,
    source_type,
    metadata,
    content
FROM knowledge_base
ORDER BY embedding <=> $1::vector
LIMIT 5