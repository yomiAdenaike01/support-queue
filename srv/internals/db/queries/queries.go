package queries

import _ "embed"

//go:embed schema.sql
var SCHEMA_SQL_QUERY string

//go:embed create_ticket.sql
var CREATE_TICKET_QUERY string

//go:embed create_team.sql
var CREATE_TEAM_QUERY string

//go:embed find_ticket_by_id.sql
var FIND_TICKET_BY_ID_QUERY string

//go:embed create_message.sql
var INSERT_MESSAGE_QUERY string

//go:embed find_team.sql
var SEARCH_TEAM_QUERY string

//go:embed update_team.sql
var UPDATE_TEAM_QUERY string

//go:embed dashboard_summary.sql
var DASHBOARD_SUMMARY string

//go:embed find_tickets.sql
var FIND_TICKETS string

//go:embed insert_ticket_event.sql
var INSERT_TICKET_EVENT string

//go:embed insert_knowledgebase.sql
var INSERT_KNOWLEDGE_BASE string
