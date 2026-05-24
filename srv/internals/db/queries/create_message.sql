insert into messages (ticket_pk, ticket_id, content, role)
select pk,
       id,
       $1 as content,
       $2 as role
from tickets
where id = $3
    and customer_email = $4 returning id,
                                      content,
                                      ticket_id,
                                      ticket_pk,
                                      role