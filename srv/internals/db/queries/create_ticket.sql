with ins_ticket as
    (INSERT INTO tickets (customer_email, subject, message)
     VALUES (:customer_email,
             :subject,
             :message_content) ON CONFLICT DO NOTHING RETURNING id,
                                                                pk,
                                                                customer_email),
     ins_message as
    (INSERT INTO messages (ticket_id, ticket_pk, content, role) SELECT ins_ticket.id as ticket_id,
                                                                       ins_ticket.pk as ticket_pk,
                                                                       :message_content as content,
                                                                       'customer' as role
     from ins_ticket returning id,
                               pk,
                               content,
                               role,
                               ticket_pk)
select ins_ticket.id as ticket_id,
       ins_ticket.pk as ticket_pk,
       ins_ticket.customer_email as customer_email,
       ins_message.id as message_id,
       ins_message.pk as message_pk,
       ins_message.content as message_content,
       ins_message.role as role
from ins_ticket
left join ins_message on ins_message.ticket_pk = ins_ticket.pk