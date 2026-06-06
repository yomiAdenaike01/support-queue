SELECT
    id as ticket_id,
    pk as ticket_pk,
    category,
    status,
    priority,
    customer_email,
    suggested_response,
    assigned_team,
    created_at
FROM tickets 
WHERE ($1::TEXT IS NULL OR subject ILIKE '%' || $1::TEXT || '%')
AND ($1::TEXT IS NULL OR customer_email ILIKE '%' || $1::TEXT || '%')
AND ($2::TEXT IS NULL OR status = $2::TEXT)
AND ($3::TEXT IS NULL OR priority = $3::TEXT)
AND ($4::TEXT IS NULL OR category = $4::TEXT)
ORDER BY created_at DESC
OFFSET 0
LIMIT $5