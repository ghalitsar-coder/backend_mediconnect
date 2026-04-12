package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"mediconnect/internal/domain"
)

type facilityRepository struct {
	db *pgxpool.Pool
}

// NewFacilityRepository returns a domain.FacilityRepository backed by PostgreSQL.
func NewFacilityRepository(db *pgxpool.Pool) domain.FacilityRepository {
	return &facilityRepository{db: db}
}

func (r *facilityRepository) GetFacilities(ctx context.Context, filter domain.FacilityFilter) ([]domain.Facility, error) {
	query := `
		SELECT id, name, address, lat, lng, type, district_id, is_active, created_at, updated_at
		FROM   facilities
		WHERE  is_active = true`

	var args []any

	if filter.DistrictID != "" {
		args = append(args, filter.DistrictID)
		query += fmt.Sprintf(" AND district_id = $%d", len(args))
	}
	if filter.Type != "" {
		args = append(args, filter.Type)
		query += fmt.Sprintf(" AND type = $%d", len(args))
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query facilities: %w", err)
	}
	defer rows.Close()

	var facilities []domain.Facility
	for rows.Next() {
		var f domain.Facility
		if err := rows.Scan(
			&f.ID, &f.Name, &f.Address, &f.Lat, &f.Lng,
			&f.Type, &f.DistrictID, &f.IsActive,
			&f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan facility row: %w", err)
		}
		facilities = append(facilities, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate facility rows: %w", err)
	}

	return facilities, nil
}
