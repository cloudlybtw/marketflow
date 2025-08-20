CREATE TABLE prices (
    pair_name TEXT NOT NULL,
    exchange TEXT NOT NULL,
    timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    average_price FLOAT NOT NULL,
    min_price FLOAT NOT NULL,
    max_price FLOAT NOT NULL
);