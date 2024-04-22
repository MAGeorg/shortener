-- +migration Up

CREATE TABLE shot_url (
    id         SERIAL  NOT NULL,
    hash_value INTEGER NOT NULL,
    origin_url TEXT    UNIQUE,
    PRIMARY KEY (id)
);

-- +migration Down
DROP TABLE shot_url;