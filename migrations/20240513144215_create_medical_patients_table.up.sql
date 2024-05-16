CREATE TABLE IF NOT EXISTS medical_patients (
    id VARCHAR(26) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    identity_number BIGINT NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    name VARCHAR(30) NOT NULL,
    birth_date VARCHAR NOT NULL,
    gender VARCHAR(6) NOT NULL,
    id_card_img VARCHAR NOT NULL,

    CONSTRAINT indentity_number_medical_patients_unique UNIQUE(identity_number),
    CONSTRAINT gender_medical_patients_check CHECK(gender IN ('male', 'female'))
);
