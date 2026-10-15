CREATE TABLE IF NOT EXISTS balances (
    account_id VARCHAR(255) NOT NULL,
    balance DOUBLE NOT NULL DEFAULT 0,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (account_id)
);
