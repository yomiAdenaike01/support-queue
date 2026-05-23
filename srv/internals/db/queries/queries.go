package queries

import _ "embed"

//go:embed schema.sql
var SCHEMA_SQL_QUERY string

//go:embed create_ticket.sql
var CREATE_TICKET_QUERY string

//go:embed find_ticket_by_id.sql
var FIND_TICKET_BY_ID_QUERY string
