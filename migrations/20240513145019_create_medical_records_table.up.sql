CREATE TABLE IF NOT EXISTS medical_records (
    id VARCHAR(26) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    identity_number BIGINT NOT NULL,
    symptoms VARCHAR(2000) NOT NULL,
    medications VARCHAR(2000) NOT NULL,
    medical_patient_id VARCHAR(26) NOT NULL,
    created_by_id VARCHAR(26) NOT NULL,

    CONSTRAINT medical_patient_id_medical_record_fk FOREIGN KEY(medical_patient_id) REFERENCES medical_patients(id),
    CONSTRAINT created_by_id_users_fk FOREIGN KEY(created_by_id) REFERENCES users(id)
);
