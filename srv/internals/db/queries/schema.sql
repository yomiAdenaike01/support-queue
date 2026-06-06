CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";

-- TICKETS
CREATE TABLE IF NOT EXISTS tickets (
    pk SERIAL PRIMARY KEY,
    id UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    customer_email TEXT NOT NULL,
    subject TEXT NOT NULL,
    message TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'PENDING' CHECK (
        status IN (
            'PENDING',
            'PROCESSING',
            'PROCESSED',
            'FAILED',
            'RESOLVED'
        )
    ),
    priority TEXT CHECK (
        priority IN (
            'LOW',
            'MEDIUM',
            'HIGH',
            'URGENT'
        )
    ),
    category TEXT CHECK (
        category IN (
            'BILLING',
            'ACCOUNT_ACCESS',
            'TECHNICAL',
            'DELIVERY',
            'CANCELLATION',
            'GENERAL'
        )
    ),
    assigned_team TEXT,
    suggested_response TEXT,
    retry_count INT NOT NULL DEFAULT 0,
    error_message TEXT,
    urgency_flag BOOLEAN NOT NULL DEFAULT FALSE,
    urgency_reason TEXT,
    sentiment_score FLOAT,
    sentiment_label TEXT CHECK (
        sentiment_label IN (
            'HIGH',
            'MEDIUM',
            'LOW'
        )
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_tickets_id ON tickets(id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_priority ON tickets(priority);
CREATE INDEX IF NOT EXISTS idx_tickets_created_at ON tickets(created_at DESC);

-- TICKET EVENTS
CREATE TABLE IF NOT EXISTS ticket_events (
    pk SERIAL PRIMARY KEY,
    id UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    ticket_pk INT NOT NULL REFERENCES tickets(pk) ON DELETE CASCADE,
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL CHECK (
        event_type IN (
            'TICKET_CREATED',
            'TICKET_PROCESSING',
            'TICKET_CLASSIFIED',
            'TICKET_PROCESSED',
            'TICKET_FAILED',
            'TICKET_REPROCESSED',
            'TICKET_RESOLVED'
        )
    ),
    payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_events_ticket_pk ON ticket_events(ticket_pk);
CREATE INDEX IF NOT EXISTS idx_events_ticket_id ON ticket_events(ticket_id);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON ticket_events(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_events_payload ON ticket_events USING GIN(payload);
CREATE UNIQUE INDEX IF NOT EXISTS idx_events_ticket_id_event_type_unique ON ticket_events(ticket_id, event_type);

-- WORKER STATUS
CREATE TABLE IF NOT EXISTS worker_status (
    pk SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    language TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'OFFLINE' CHECK (
        status IN (
            'ONLINE',
            'OFFLINE',
            'DEGRADED'
        )
    ),
    processed_count INT NOT NULL DEFAULT 0,
    failed_count INT NOT NULL DEFAULT 0,
    last_heartbeat TIMESTAMPTZ
);

-- MESSAGES
CREATE TABLE IF NOT EXISTS messages (
    pk SERIAL PRIMARY KEY,
    id UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    ticket_pk INT NOT NULL REFERENCES tickets(pk) ON DELETE CASCADE,
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (
        role IN (
            'customer',
            'assistant'
        )
    ),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE messages
    ADD COLUMN IF NOT EXISTS ticket_id UUID REFERENCES tickets(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_messages_ticket_pk ON messages(ticket_pk);
CREATE INDEX IF NOT EXISTS idx_messages_ticket_id ON messages(ticket_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at ASC);

-- TEAM
CREATE TABLE IF NOT EXISTS team (
    pk SERIAL PRIMARY KEY,
    id UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    department TEXT NOT NULL UNIQUE,
    integrations JSONB DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_team_department ON team(department);

-- TEAM MEMBERS
CREATE TABLE IF NOT EXISTS team_members (
    pk SERIAL PRIMARY KEY,
    id UUID NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    email_address TEXT UNIQUE DEFAULT NULL,
    phone_number TEXT UNIQUE DEFAULT NULL,
    team_id UUID NOT NULL REFERENCES team(id) ON DELETE CASCADE,
    team_pk INT NOT NULL REFERENCES team(pk) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL UNIQUE,
    role TEXT NOT NULL,
    integrations JSONB DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_team_member_email ON team_members(email_address);
CREATE INDEX IF NOT EXISTS idx_team_member_phonenumber ON team_members(phone_number);
CREATE TABLE IF NOT EXISTS knowledge_base(
    id uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
    pk SERIAL PRIMARY KEY,
    source_id UUID NOT NULL,
    metadata jsonb DEFAULT NULL,
    embedding vector(384) NOT NULL,
    source_type TEXT DEFAULT NULL,
    content TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_knowledge_base_source_id ON knowledge_base(source_id);
