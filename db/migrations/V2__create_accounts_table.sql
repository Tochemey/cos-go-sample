-- accounts relation
CREATE TABLE sample.accounts(
    account_id VARCHAR(255) NOT NULL,
    account_balance NUMERIC(19, 2) NOT NULL,
    account_owner VARCHAR(255) NOT NULL,
    is_closed BOOLEAN NOT NULL,

    PRIMARY KEY (account_id)
);