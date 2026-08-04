CREATE TABLE IF NOT EXISTS agents (

    id UUID PRIMARY KEY,

    name TEXT NOT NULL,

    hostname TEXT NOT NULL,

    ip_address TEXT NOT NULL,

    operating_system TEXT,

    version TEXT,

    status TEXT,

    created_at TIMESTAMP NOT NULL,

    updated_at TIMESTAMP NOT NULL,

    deleted_at TIMESTAMP
);