-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cases (
    id UUID PRIMARY KEY,
    version INTEGER NOT NULL
);

CREATE TABLE slides (
    id UUID PRIMARY KEY,
    version INTEGER NOT NULL,
    preparation_status SMALLINT NOT NULL,
    case_id UUID
);

CREATE TABLE events (
    id UUID PRIMARY KEY,
    type SMALLINT NOT NULL,
    case_id UUID REFERENCES cases(id) ON DELETE CASCADE,
    slide_id UUID REFERENCES slides(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    published BOOLEAN NOT NULL
);

CREATE TABLE case_projections (
    id UUID PRIMARY KEY,
    status SMALLINT NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS slides;
DROP TABLE IF EXISTS cases;
DROP TABLE IF EXISTS case_projections;

DROP EXTENSION IF EXISTS "uuid-ossp";
