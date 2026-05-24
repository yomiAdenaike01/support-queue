with paginated_messages as
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
     limit $3)
SELECT t.id::text AS ticket_id,
       t.pk as ticket_pk,
       t.customer_email,
       t.subject,
       t.status,
       t.priority,
       t.category,
       t.assigned_team,
       t.suggested_response,
       m.id AS message_id,
       m.content AS message_content,
       m.role AS message_role
FROM tickets t
LEFT JOIN paginated_messages m ON m.ticket_pk = t.pk
WHERE t.id = $1