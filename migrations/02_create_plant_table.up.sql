
CREATE TABLE Plant (
    id SERIAL PRIMARY KEY,
    stead_owner INTEGER REFERENCES stead(id),
    kind TEXT NOT NULL DEFAULT 'dirt',
    xp INTEGER NOT NULL DEFAULT 0,
    xp_multiplier FLOAT NOT NULL DEFAULT 1,
    next_yield TIMESTAMP
);
