INSERT INTO knowledge_base(source_id, source_type, metadata, content, embedding)
VALUES ($1::uuid, $2::TEXT, COALESCE($3::jsonb, NULL::jsonb), $4::TEXT, $5::vector)
RETURNING id, pk, source_type, metadata, source_id, content
