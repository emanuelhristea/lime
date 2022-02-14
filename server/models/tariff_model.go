package models

import (
	"time"

	"github.com/emanuelhristea/lime/config"
	"github.com/jinzhu/gorm"
)

// Tariff is a ...
type Tariff struct {
	ID            uint64         `gorm:"primary_key;auto_increment" json:"id"`
	Name          string         `gorm:"size:255;not null;unique" json:"name"`
	Price         int            `gorm:"size:6;not null" json:"price"`
	Tandem        bool           `gorm:"size:1;not null" json:"crossbar"`
	Triaxis       bool           `gorm:"size:1;not null" json:"triaxis"`
	Robots        bool           `gorm:"size:1;not null" json:"robots"`
	Period        int            `gorm:"size:6;not null" json:"period"`
	Users         int            `gorm:"size:6;not null" json:"users"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Subscriptions []Subscription `json:"subscriptions"`
}

// SaveTariff is a ...
func (t *Tariff) SaveTariff() (*Tariff, error) {
	err := config.DB.Create(&t).Error
	if err != nil {
		return &Tariff{}, err
	}
	return t, nil
}

// FindTariffByID is a ...
func (t *Tariff) FindTariffByID(uid uint64) (*Tariff, error) {
	err := config.DB.Model(Tariff{}).Where("id = ?", uid).Take(&t).Error
	if err != nil {
		return &Tariff{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Tariff{}, ErrTariffNotFound
	}
	return t, err
}

// UpdateTariff is a ...
func (t *Tariff) UpdateTariff(uid uint64) (*Tariff, error) {
	db := config.DB.Model(&Tariff{}).Where("id = ?", uid).Take(&Tariff{}).UpdateColumns(
		map[string]interface{}{
			"name":       t.Name,
			"price":      t.Price,
			"tandem":     t.Tandem,
			"triaxis":    t.Triaxis,
			"robots":     t.Robots,
			"period":     t.Period,
			"users":      t.Users,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Tariff{}, db.Error
	}

	err := db.Model(&Tariff{}).Where("id = ?", uid).Take(&t).Error
	if err != nil {
		return &Tariff{}, err
	}
	return t, nil
}

// DeleteTariff is a ...
func DeleteTariff(uid uint64) (int64, error) {
	db := config.DB.Model(&Tariff{}).Where("id = ?", uid).Take(&Tariff{}).Delete(&Tariff{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// TariffsList is a ...
func TariffsList(relations ...string) *[]Tariff {
	db := config.DB.Model(&Tariff{}).Order("ID asc")
	for _, rel := range relations {
		db = db.Preload(rel)
	}

	tariffs := []Tariff{}
	db = db.Find(&tariffs)
	if db.Error != nil {
		return &tariffs
	}
	return &tariffs
}
