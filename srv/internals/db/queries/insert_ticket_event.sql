INSERT INTO ticket_events (ticket_id, ticket_pk, payload, event_type)
SELECT
    id,
    pk,
    $1::jsonb,
    $2::TEXT
FROM tickets
WHERE id = $3::uuid
ON CONFLICT (ticket_id, event_type) DO NOTHING
RETURNING id as event_id, pk as event_pk, ticket_id, ticket_pk, payload, event_type, created_at
