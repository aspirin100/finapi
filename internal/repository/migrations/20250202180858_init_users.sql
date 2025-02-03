-- +goose Up
-- +goose StatementBegin
create table if not exists bank_accounts(
    userID uuid primary key,
    balance decimal check(balance >= 0)
);

create table if not exists transactions(
    id uuid primary key,
    receiverID uuid not null,
    senderID uuid not null,
    amount decimal not null,
    createdAt timestamptz default NOW()
);

alter table transactions
    add constraint fk_receiver_id
    foreign key (receiverID) references bank_accounts(userID);

alter table transactions
    add constraint fk_sender_id
    foreign key (senderID) references bank_accounts(userID);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists transactions;
drop table if exists bank_accounts;
-- +goose StatementEnd
