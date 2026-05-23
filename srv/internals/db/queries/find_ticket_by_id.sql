SELECT t.id::text AS ticket_id,
       t.customer_email,
       t.subject,
       t.status,
       t.priority,
       t.category,
       t.assigned_team,
       t.suggested_response,
       m.id::text AS message_id,
       m.content AS message_content,
       m.role AS message_role
FROM tickets t
LEFT JOIN messages m ON m.ticket_pk = t.pk
WHERE t.id = $1
ORDER BY m.created_at ASC
