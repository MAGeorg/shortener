-- +migration Up

CREATE TABLE shot_url (
    id         SERIAL NOT NULL,
    hash_value BIGSERIAL NOT NULL,
    origin_url VARCHAR(2048) UNIQUE,
    PRIMARY KEY (id)
);

-- +migration Down
DROP TABLE shot_url;