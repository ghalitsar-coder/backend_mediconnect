package postgres

import (
  "context"
  "mediconnect/internal/domain"
  "gorm.io/gorm"
)

type FacilityRepository struct {
  db *gorm.DB
}

func NewFacilityRepository(db *gorm.DB) domain.FacilityRepository {
  return &FacilityRepository{db: db}
}

func (r *FacilityRepository) GetFacilities(ctx context.Context, filter domain.FacilityFilter) ([]domain.Facility, error) {
  var facilities []domain.Facility
  query := r.db.WithContext(ctx).Table("facilities")
  
  if filter.Type != "" {
    query = query.Where("type = ?", filter.Type)
  }
  if filter.DistrictID != "" {
    query = query.Where("district_id = ?", filter.DistrictID)
  }
  
  if err := query.Find(&facilities).Error; err != nil {
    return nil, err
  }
  return facilities, nil
}
