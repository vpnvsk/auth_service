CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- INSERT INTO users (id, login, password_hash)
-- VALUES (uuid_generate_v4(), 'new_user', 'password_hash');

CREATE TABLE IF NOT EXISTS users
(
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    login      varchar(255) not null unique,
    password_hash BYTEA not null
);
CREATE INDEX IF NOT EXISTS idx_login ON users (login);

CREATE TABLE IF NOT EXISTS apps
(
    id     SERIAL PRIMARY KEY,
    name   varchar(255) NOT NULL UNIQUE,
    secret varchar(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS admins
(
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    app_id int REFERENCES apps(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, app_id)
);