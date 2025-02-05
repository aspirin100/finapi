-- +goose Up
-- +goose StatementBegin
INSERT INTO bank_accounts (userID, balance)
VALUES ('3fec06e9-29cc-4ff4-9ae7-fb0e7c757b61', 0);

INSERT INTO bank_accounts (userID, balance)
VALUES ('4178f61f-2ff9-4ab5-afa5-f30dc16e6ad9', 100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM transactions;
DELETE FROM bank_accounts;
-- +goose StatementEnd