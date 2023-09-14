CREATE TYPE plant_kind AS ENUM('bbc', 'hvv', 'cyl', 'dirt');
CREATE TABLE Plant (
    id SERIAL PRIMARY KEY,
    stead_owner INTEGER REFERENCES stead(id),
    kind plant_kind NOT NULL DEFAULT 'dirt',
    xp INTEGER NOT NULL DEFAULT 0,
    xp_multiplier FLOAT NOT NULL DEFAULT 1,
    next_yield TIMESTAMP
);
