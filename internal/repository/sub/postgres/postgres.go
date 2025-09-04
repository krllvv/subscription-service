package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"subscription-service/config"
	"subscription-service/internal/model"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
	ErrNotFound = errors.New("requested item not found")
	ErrDatabase = errors.New("database error")
)

type SubPostgresRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewSubPostgresRepository(cfg *config.Config, logger *log.Logger) (*SubPostgresRepository, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var db *sql.DB
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Println("Could not create database connection:", err)
		return nil, ErrDatabase
	}

	err = db.Ping()
	if err != nil {
		logger.Println("Could not connect to PostgreSQL:", err)
		return nil, ErrDatabase
	}

	logger.Println("Connected to PostgreSQL")
	return &SubPostgresRepository{db: db, logger: logger}, nil
}

func (r *SubPostgresRepository) Create(sub *model.Subscription) error {
	if sub.ID == uuid.Nil {
		sub.ID = uuid.New()
	}

	_, err := r.db.Exec(
		"INSERT INTO subs (id, name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)",
		sub.ID, sub.Name, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	)
	if err != nil {
		r.logger.Println("Failed to create subscription:", err)
		return ErrDatabase
	}

	r.logger.Printf("Successfully created subscription with ID %s", sub.ID)
	return nil
}

func (r *SubPostgresRepository) GetByID(id uuid.UUID) (*model.Subscription, error) {
	sub := &model.Subscription{}

	err := r.db.
		QueryRow("SELECT id, name, price, user_id, start_date, end_date FROM subs WHERE id = $1", id).
		Scan(&sub.ID, &sub.Name, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Printf("Subscription with ID %s not found: %v", sub.ID, err)
			return nil, ErrNotFound
		}
		r.logger.Printf("Failed to get subscription with ID %s: %v", sub.ID, err)
		return nil, ErrDatabase
	}

	r.logger.Printf("Successfully got subscription with ID %s", sub.ID)
	return sub, nil
}

func (r *SubPostgresRepository) Update(id uuid.UUID, sub *model.Subscription) error {
	res, err := r.db.Exec(
		"UPDATE subs SET name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6",
		sub.Name, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, id,
	)

	if err != nil {
		r.logger.Println("Failed to update subscription:", err)
		return ErrDatabase
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Println("Failed to get affected rows for update:", err)
		return ErrDatabase
	}
	if rows == 0 {
		r.logger.Printf("Subscription with ID %s not found: %v", id, err)
		return ErrNotFound
	}

	r.logger.Printf("Successfully updated subscription with ID %s", id)
	return nil
}

func (r *SubPostgresRepository) Delete(id uuid.UUID) error {
	res, err := r.db.Exec("DELETE FROM subs WHERE id = $1", id)
	if err != nil {
		r.logger.Println("Failed to delete subscription", err)
		return ErrDatabase
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.logger.Println("Failed to get affected rows for delete:", err)
		return ErrDatabase
	}
	if rows == 0 {
		r.logger.Printf("Subscription with ID %s not found: %v", id, err)
		return ErrNotFound
	}

	r.logger.Printf("Successfully deleted subscription with ID %s", id)
	return nil
}

func (r *SubPostgresRepository) GetAll() ([]model.Subscription, error) {
	rows, err := r.db.Query("SELECT id, name, price, user_id, start_date, end_date FROM subs")
	if err != nil {
		r.logger.Println("Failed to get all subscriptions:", err)
		return nil, ErrDatabase
	}
	defer rows.Close()

	subs := make([]model.Subscription, 0)
	for rows.Next() {
		var sub model.Subscription
		err = rows.Scan(&sub.ID, &sub.Name, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			r.logger.Println("Failed to scan row while getting subscriptions", err)
			return nil, ErrDatabase
		}
		subs = append(subs, sub)
	}

	if err = rows.Err(); err != nil {
		r.logger.Println("Failed iterating rows while getting subscriptions:", err)
		return nil, ErrDatabase

	}

	r.logger.Printf("Successfully found %d subscriptions", len(subs))
	return subs, nil
}

func (r *SubPostgresRepository) GetTotalSum(start, end string, userID uuid.UUID, name string) (int, error) {
	conditions := []string{
		"TO_DATE('01-' || start_date, 'DD-MM-YYYY') <= TO_DATE('01-' || $1, 'DD-MM-YYYY')",
		"(TO_DATE('01-' || end_date, 'DD-MM-YYYY') >= TO_DATE('01-' || $2, 'DD-MM-YYYY') OR end_date IS NULL)",
	}
	args := []interface{}{end, start}

	if userID != uuid.Nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, userID)
	}

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, name)
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT SUM(price) FROM subs WHERE ")
	queryBuilder.WriteString(strings.Join(conditions, " AND "))

	query := queryBuilder.String()
	var totalSum sql.NullInt64
	err := r.db.QueryRow(query, args...).Scan(&totalSum)
	if err != nil {
		r.logger.Println("Error calculate total sum:", err)
		return 0, ErrDatabase
	}

	if !totalSum.Valid {
		return 0, nil
	}

	res := int(totalSum.Int64)
	r.logger.Println("Calculated total sum:", res)
	return res, nil
}
