CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id              UUID DEFAULT uuid_generate_v1(),
    username        VARCHAR(250) NOT NULL,
    password        VARCHAR(250) NOT NULL,

    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS notes (
    id      UUID DEFAULT uuid_generate_v1(),
    user_id UUID,
    text    TEXT,

    PRIMARY KEY (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES "users"(id)
);


INSERT INTO users (username, password)
    VALUES ('tymbaca', 'longin');
