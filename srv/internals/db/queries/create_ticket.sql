WITH ins_ticket AS (
    INSERT INTO tickets (
        customer_email,
        subject,
        message
    )
    VALUES (
        $1,
        $2,
        $3
    )
    ON CONFLICT DO NOTHING
    RETURNING
        id,
        pk,
        customer_email
),
ins_event as (
    INSERT INTO ticket_events (
        ticket_id,
        ticket_pk,
        payload,
        event_type
    )
    SELECT 
        ins_ticket.id as ticket_id,
        ins_ticket.pk as ticket_pk,
        jsonb_build_object('customer_email',$1::TEXT,'subject',$2::TEXT,'message_content',$3::TEXT) as payload,
        'TICKET_CREATED'::TEXT
    FROM 
        ins_ticket
),
ins_message AS (
    INSERT INTO messages (
        ticket_id,
        ticket_pk,
        content,
        role
    )
    SELECT
        ins_ticket.id AS ticket_id,
        ins_ticket.pk AS ticket_pk,
        $3::TEXT AS content,
        'customer' AS role
    FROM ins_ticket
    RETURNING
        id,
        pk,
        content,
        role,
        ticket_pk
)
SELECT
    ins_ticket.id AS ticket_id,
    ins_ticket.pk AS ticket_pk,
    ins_ticket.customer_email AS customer_email,
    ins_message.id AS message_id,
    ins_message.pk AS message_pk,
    ins_message.content AS message_content,
    ins_message.role AS role
FROM ins_ticket
LEFT JOIN ins_message ON ins_message.ticket_pk = ins_ticket.pk;