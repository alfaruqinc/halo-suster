CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(26) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    nip BIGINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    password VARCHAR NULL,
    role VARCHAR(5) NOT NULL,
    id_card_img VARCHAR NULL,

    CONSTRAINT nip_users_unique UNIQUE(nip),
    CONSTRAINT role_users_check CHECK(role IN ('it', 'nurse'))
);
