-- +migration Up

CREATE TABLE shot_url (
    id         SERIAL  NOT NULL,
    hashValue  INTEGER NOT NULL,
    origin_url TEXT    UNIQUE,
    hash_url   TEXT    NOT NULL,
    PRIMARY KEY (id)
);

-- +migration Down
DROP TABLE shot_url;