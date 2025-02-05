-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bank_accounts (
    userID UUID PRIMARY KEY,
    balance DECIMAL CHECK (balance >= 0)
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    receiverID UUID NOT NULL,
    senderID UUID NOT NULL,
    operation VARCHAR(8) NOT NULL,
    amount DECIMAL NOT NULL,
    createdAt TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

ALTER TABLE transactions
    ADD CONSTRAINT fk_receiver_id
    FOREIGN KEY (receiverID) REFERENCES bank_accounts(userID);

ALTER TABLE transactions
    ADD CONSTRAINT fk_sender_id
    FOREIGN KEY (senderID) REFERENCES bank_accounts(userID);

CREATE INDEX IF NOT EXISTS senderid_index ON transactions USING HASH (senderID);

CREATE INDEX IF NOT EXISTS receiverid_index ON transactions USING HASH (receiverID);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS bank_accounts;
-- +goose StatementEnd