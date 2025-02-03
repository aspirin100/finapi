-- +goose Up
-- +goose StatementBegin
create table if not exists bank_accounts(
    userID uuid primary key,
    balance decimal
);

create table if not exists transactions(
    id uuid primary key,
    receiverID uuid not null,
    senderID uuid not null,
    amount decimal not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists bank_accounts;
drop table if exists transactions;
-- +goose StatementEnd
