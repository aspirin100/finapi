-- +goose Up
-- +goose StatementBegin
insert into bank_accounts(userID, balance) values ('3fec06e9-29cc-4ff4-9ae7-fb0e7c757b61', 0);
insert into bank_accounts(userID, balance) values ('4178f61f-2ff9-4ab5-afa5-f30dc16e6ad9', 100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
