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
    operation varchar(8) not null,
    amount decimal not null,
    createdAt timestamptz default NOW()
);

alter table transactions
    add constraint fk_receiver_id
    foreign key (receiverID) references bank_accounts(userID);

alter table transactions
    add constraint fk_sender_id
    foreign key (senderID) references bank_accounts(userID);

CREATE INDEX if not exists senderid_index ON transactions USING hash ("senderid");
CREATE INDEX if not exists receiverid_index ON transactions USING hash ("receiverid");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists transactions;
drop table if exists bank_accounts;
-- +goose StatementEnd
