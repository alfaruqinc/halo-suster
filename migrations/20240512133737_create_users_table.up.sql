CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(26) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    nip INT NOT NULL,
    name VARCHAR(50) NOT NULL,
    password VARCHAR(33) NULL,
    role VARCHAR(5) NOT NULL,
    id_card_img VARCHAR NULL
);
