CREATE TABLE Stead (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    inventory JSONB,
    ephemeral_statuses JSONB
);


