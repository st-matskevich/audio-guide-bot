CREATE TABLE tickets(
    ticket_id BIGSERIAL PRIMARY KEY, 
    code VARCHAR(64) NOT NULL,
    used BOOLEAN NOT NULL DEFAULT false);