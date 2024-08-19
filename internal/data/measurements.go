package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"householdingindex.homecatalogue.net/internal/validator"
)

type MeasurementModel struct {
	DB *sql.DB
}

func (mm MeasurementModel) Insert(measurement *Measurement) error {
	query := `
		INSERT INTO measurements (name)
		VALUES ($1)
		RETURNING id, created_at, version`

	args := []interface{}{measurement.Name}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return mm.DB.QueryRowContext(ctx, query, args...).Scan(&measurement.ID, &measurement.CreatedAt, &measurement.Version)
}

func (it MeasurementModel) Get(id int64) (*Measurement, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, name, version
		FROM measurements
		WHERE id = $1`

	var measurement Measurement

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := it.DB.QueryRowContext(ctx, query, id).Scan(
		&measurement.ID,
		&measurement.CreatedAt,
		&measurement.Name,
		&measurement.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &measurement, nil
}

func (mm MeasurementModel) GetAll(name string, filters Filters) ([]*Measurement, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, name, version
		FROM measurements
		WHERE (name = $1 OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, filters.limit(), filters.offset()}

	rows, err := mm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	measurements := []*Measurement{}

	for rows.Next() {
		var measurement Measurement

		err := rows.Scan(
			&totalRecords,
			&measurement.ID,
			&measurement.CreatedAt,
			&measurement.Name,
			&measurement.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		measurements = append(measurements, &measurement)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return measurements, metadata, nil
}

func (mm MeasurementModel) Update(measurement *Measurement) error {
	query := `
		UPDATE measurements
		SET name = $1, version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version`

	args := []interface{}{
		measurement.Name,
		measurement.ID,
		measurement.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := mm.DB.QueryRowContext(ctx, query, args...).Scan(&measurement.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (mm MeasurementModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM measurements
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := mm.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

type Measurement struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Version   int32     `json:"version"`
}

func ValidateMeasurement(v *validator.Validator, measurement *Measurement) {
	v.Check(measurement.Name != "", "name", "must be provided")
	v.Check(len(measurement.Name) <= 500, "name", "must not be more than 500 bytes long")
	//v.Check(validator.Unique(input.Name), "name", "must not contain duplicate values")
}
