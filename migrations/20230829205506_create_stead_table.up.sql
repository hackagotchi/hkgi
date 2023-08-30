CREATE TABLE Stead (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    inventory JSONB,
    ephemeral_statuses TEXT[]
);

CREATE TABLE stead_plant (
    stead INTEGER references stead(id),
    plant INTEGER references plant(id)
);
