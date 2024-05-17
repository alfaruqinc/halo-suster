package repository

import (
	"context"
	"fmt"
	"health-record/internal/domain"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MedicalPatient interface {
	Create(ctx context.Context, db *sqlx.DB, patient domain.MedicalPatient) error
	GetAllMedicalPatients(ctx context.Context, db *sqlx.DB, queryParams domain.MedicalPatientQueryParams) ([]domain.GetMedicalPatient, error)
}

type medicalPatient struct{}

func NewMedicalPatient() MedicalPatient {
	return &medicalPatient{}
}

func (mp *medicalPatient) Create(ctx context.Context, db *sqlx.DB, patient domain.MedicalPatient) error {
	query := `INSERT INTO medical_patients (id, created_at, identity_number, phone_number, name, birth_date, gender, id_card_img)
	VALUES (:id, :created_at, :identity_number, :phone_number, :name, :birth_date, :gender, :id_card_img)`
	_, err := db.NamedExecContext(ctx, query, patient)
	if err != nil {
		return err
	}

	return nil
}

func (mp *medicalPatient) GetAllMedicalPatients(ctx context.Context, db *sqlx.DB, queryParams domain.MedicalPatientQueryParams) ([]domain.GetMedicalPatient, error) {
	var whereClause []string
	var args []any
	argPos := 1

	if queryParams.IdentityNumber != "" {
		if _, err := strconv.Atoi(queryParams.IdentityNumber); err == nil {
			whereClause = append(whereClause, fmt.Sprintf("identity_number = $%d", argPos))
			args = append(args, queryParams.IdentityNumber)
			argPos++
		}
	}

	if queryParams.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("name ILIKE $%d", argPos))
		args = append(args, "%"+queryParams.Name+"%")
		argPos++
	}

	if queryParams.PhoneNumber != "" {
		whereClause = append(whereClause, fmt.Sprintf("phone_number LIKE $%d", argPos))
		args = append(args, "+"+queryParams.PhoneNumber+"%")
		argPos++
	}

	var orderClause []string
	createdAt := "DESC"
	if queryParams.CreatedAt == "asc" {
		createdAt = "ASC"
	}
	orderClause = append(orderClause, fmt.Sprintf("created_at %s", createdAt))

	var limitOffsetClause []string
	limit := "5"
	if queryParams.Limit != "" {
		if _, err := strconv.Atoi(queryParams.Limit); err == nil {
			limit = queryParams.Limit
		}
	}
	limitOffsetClause = append(limitOffsetClause, fmt.Sprintf("LIMIT $%d", argPos))
	args = append(args, limit)
	argPos++
	offset := "0"
	if queryParams.Offset != "" {
		if _, err := strconv.Atoi(queryParams.Offset); err == nil {
			offset = queryParams.Offset
		}
	}
	limitOffsetClause = append(limitOffsetClause, fmt.Sprintf("OFFSET $%d", argPos))
	args = append(args, offset)
	argPos++

	var queryCondition string
	if len(whereClause) > 0 {
		queryCondition += "\nWHERE " + strings.Join(whereClause, " AND ")
	}
	queryCondition += "\nORDER BY " + strings.Join(orderClause, ", ")
	queryCondition += "\n" + strings.Join(limitOffsetClause, " ")

	query := `SELECT created_at, identity_number, phone_number, name, birth_date, gender
	FROM medical_patients`
	query += queryCondition

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	medicalPatients := []domain.GetMedicalPatient{}
	for rows.Next() {
		var mp domain.GetMedicalPatient
		err := rows.Scan(
			&mp.CreatedAt, &mp.IdentityNumber, &mp.PhoneNumber,
			&mp.Name, &mp.BirthDate, &mp.Gender,
		)
		if err != nil {
			return nil, err
		}
		medicalPatients = append(medicalPatients, mp)
	}

	return medicalPatients, nil
}
