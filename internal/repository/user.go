package repository

import (
	"context"
	"fmt"
	"health-record/internal/domain"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type User interface {
	GetAllUser(ctx context.Context, db *sqlx.DB, queryParams domain.UserQueryParams) ([]domain.UserResponse, error)
}

type user struct{}

func NewUser() User {
	return &user{}
}

func (u *user) GetAllUser(ctx context.Context, db *sqlx.DB, queryParams domain.UserQueryParams) ([]domain.UserResponse, error) {
	var whereClause []string
	var args []any
	argPos := 1

	if queryParams.ID != "" {
		whereClause = append(whereClause, fmt.Sprintf("id = $%d", argPos))
		args = append(args, queryParams.ID)
		argPos++
	}

	if queryParams.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("name ILIKE $%d", argPos))
		args = append(args, "%"+queryParams.Name+"%")
		argPos++
	}

	if queryParams.NIP != "" {
		if _, err := strconv.Atoi(queryParams.NIP); err == nil {
			whereClause = append(whereClause, fmt.Sprintf("CAST(nip AS VARCHAR(13)) LIKE $%d", argPos))
			args = append(args, queryParams.NIP+"%")
			argPos++
		}
	}

	prefixNIP := ""
	if queryParams.Role == "it" {
		prefixNIP = "615"
	}
	if queryParams.Role == "nurse" {
		prefixNIP = "303"
	}
	if prefixNIP != "" {
		whereClause = append(whereClause, fmt.Sprintf("CAST(nip AS VARCHAR(13)) LIKE $%d", argPos))
		args = append(args, prefixNIP+"%")
		argPos++
	}

	// ORDER
	var orderClause []string
	createdAt := "desc"
	if queryParams.CreatedAt == "asc" {
		createdAt = "asc"
	}
	orderClause = append(orderClause, fmt.Sprintf("created_at %s", createdAt))

	// PAGINATION
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

	query := `SELECT id, created_at, nip, name
	FROM users`
	query += queryCondition

	users := []domain.UserResponse{}
	rows, err := db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var u domain.UserResponse
		err = rows.Scan(&u.ID, &u.CreatedAt, &u.NIP, &u.Name)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}
