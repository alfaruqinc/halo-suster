package repository

import (
	"context"
	"fmt"
	"health-record/internal/domain"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MedicalRecord interface {
	Create(ctx context.Context, db *sqlx.DB, record domain.MedicalRecord) (int64, error)
	GetAllMedicalRecords(ctx context.Context, db *sqlx.DB, queryParams domain.MedicalRecordQueryParams) ([]domain.GetMedicalRecord, error)
}

type medicalRecord struct{}

func NewMedicalRecord() MedicalRecord {
	return &medicalRecord{}
}

func (mr *medicalRecord) Create(ctx context.Context, db *sqlx.DB, record domain.MedicalRecord) (int64, error) {
	query := `INSERT INTO medical_records (medical_patient_id, id, created_at, identity_number, symptoms, medications, created_by_id)
	SELECT mp.id, :id, :created_at, :identity_number, :symptoms, :medications, :created_by_id
	FROM medical_patients mp
	WHERE identity_number = :identity_number`
	res, err := db.NamedExecContext(ctx, query, record)
	if err != nil {
		return 0, err
	}
	affRow, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affRow, nil
}

func (mr *medicalRecord) GetAllMedicalRecords(ctx context.Context, db *sqlx.DB, queryParams domain.MedicalRecordQueryParams) ([]domain.GetMedicalRecord, error) {
	var whereClause []string
	var args []any
	argPos := 1

	if queryParams.IdentityDetail.IdentityNumber != "" {
		if _, err := strconv.Atoi(queryParams.IdentityDetail.IdentityNumber); err == nil {
			whereClause = append(whereClause, fmt.Sprintf("mr.identity_number = $%d", argPos))
			args = append(args, queryParams.IdentityDetail.IdentityNumber)
			argPos++
		}
	}

	if queryParams.CreatedBy.ID != "" {
		whereClause = append(whereClause, fmt.Sprintf("mr.created_by_id = $%d", argPos))
		args = append(args, queryParams.CreatedBy.ID)
		argPos++
	}

	if queryParams.CreatedBy.NIP != "" {
		if _, err := strconv.Atoi(queryParams.CreatedBy.NIP); err == nil {
			whereClause = append(whereClause, fmt.Sprintf("u.nip = $%d", argPos))
			args = append(args, queryParams.CreatedBy.NIP)
			argPos++
		}
	}

	var orderClause []string
	createdAt := "DESC"
	if queryParams.CreatedAt == "asc" {
		createdAt = "ASC"
	}
	orderClause = append(orderClause, fmt.Sprintf("mr.created_at %s", createdAt))
	orderClause = append(orderClause, "mr.id DESC")

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

	var queryCondition string
	var joinWhere string
	if len(whereClause) > 0 {
		joinWhere = "WHERE " + strings.Join(whereClause, " AND ")
	}
	queryCondition += "\nORDER BY " + strings.Join(orderClause, ", ")

	subqueryMedicalRecord := fmt.Sprintf(`
		WITH paginated_mr AS (
			SELECT mr.id, mr.created_at, mr.identity_number, mr.symptoms, 
				mr.medications, mr.medical_patient_id, mr.created_by_id,
				u.nip
			FROM medical_records mr
			JOIN users u ON u.id = mr.created_by_id
			%s
			%s
		)
	`, joinWhere, strings.Join(limitOffsetClause, " "))

	query := `SELECT mr.created_at, mr.symptoms, mr.medications,
		mp.identity_number, mp.phone_number, mp.name, mp.birth_date, mp.gender, mp.id_card_img,
		u.id, u.nip, u.name
	FROM paginated_mr mr
	JOIN medical_patients mp ON mp.id = mr.medical_patient_id
	JOIN users u ON u.id = mr.created_by_id`
	query = subqueryMedicalRecord + query + queryCondition

	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	records := []domain.GetMedicalRecord{}
	for rows.Next() {
		r := domain.GetMedicalRecord{}
		err := rows.Scan(
			&r.CreatedAt, &r.Symptoms, &r.Medications,
			&r.IdentityDetail.IdentityNumber, &r.IdentityDetail.PhoneNumber, &r.IdentityDetail.Name,
			&r.IdentityDetail.BirthDate, &r.IdentityDetail.Gender, &r.IdentityDetail.IDCardImg,
			&r.CreatedBy.ID, &r.CreatedBy.NIP, &r.CreatedBy.Name,
		)
		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil
}
