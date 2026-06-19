WITH events_by_ticket AS (
    SELECT 
        ticket_pk,
        event_type,
        created_at,
        payload
    FROM ticket_events
    WHERE ticket_id = $1
),
 paginated_messages as
    (select content,
            ticket_pk,
            ticket_id,
            id,
            role,
            created_at
     from messages
     where ticket_id = $1
     order by pk desc
     offset $2
     limit $3
    )
SELECT t.id::text AS ticket_id,
       t.pk as ticket_pk,
       t.customer_email,
       t.subject,
       t.status,
       t.priority,
       t.category,
       t.assigned_team,
       t.suggested_response,
       t.created_at::timestamp as created_at,
       m.id AS message_id,
       m.content AS message_content,
       m.role AS message_role,
       m.created_at::timestamp as message_created_at,
       jsonb_agg(jsonb_build_object('eventType', e.event_type, 'createdAt', e.created_at, 'payload', e.payload::jsonb)) FILTER (WHERE e.ticket_pk IS NOT NULL) as events
FROM tickets t
LEFT JOIN paginated_messages m ON m.ticket_pk = t.pk
LEFT JOIN events_by_ticket e ON e.ticket_pk = t.pk
WHERE t.id = $1
GROUP BY t.id, t.pk, t.customer_email, t.subject, t.status, t.priority, t.category, t.assigned_team, t.suggested_response, m.id, m.content, m.role, t.created_at, m.created_at