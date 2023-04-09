CREATE TABLE "users" (
    id varchar(64) primary key unique,
    created_at timestamptz,
    updated_at timestamptz,
    name varchar(100),
    balance float4
);

CREATE TABLE "transactions" (
    id serial primary key unique,
    created_at timestamptz,
    amount float4,
    transaction_id varchar(64),
    user_id varchar(64),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

-- transaction table indices
CREATE UNIQUE INDEX idx_transaction_id ON transactions(transaction_id);
CREATE INDEX idx_user_id ON transactions(user_id);

-- load sample data
INSERT into users (id, created_at, updated_at, name, balance) VALUES ('6d7750a1-c3f2-4765-bf8f-33bc80f3f809', now(), now(), 'Test', 100);
