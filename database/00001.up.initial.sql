BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Fixed: Changed table name from 'user' to 'users' for consistency
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name TEXT NOT NULL,
                       email TEXT UNIQUE NOT NULL,
                       phone_no TEXT NOT NULL,
                       password TEXT NOT NULL,
                       created_at TIMESTAMPTZ DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          message TEXT NOT NULL,
                          pop_up_time TIMESTAMPTZ NOT NULL,
                          status TEXT DEFAULT 'pending',
                          created_at TIMESTAMPTZ DEFAULT NOW()
);

create Table session (
    id int not null primary key ,
    user_id int references users(id),
    createdAt timestamptz not null,
    expiresAt timestamptz not null ,
    session_token text not null unique

);

CREATE INDEX idx_messages_user_id ON messages(user_id);
CREATE INDEX idx_messages_pop_up_time ON messages(pop_up_time);

COMMIT;
