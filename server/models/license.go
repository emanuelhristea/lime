package models

import (
	"time"

	"github.com/emanuelhristea/lime/config"
	"github.com/jinzhu/gorm"
)

// License is a ...
type License struct {
	ID             uint64       `gorm:"primary_key;auto_increment" json:"id"`
	SubscriptionID uint64       `sql:"type:int REFERENCES subscriptions(id) ON DELETE CASCADE" json:"subscription_id"`
	License        []byte       `gorm:"null" json:"license"`
	Hash           string       `gorm:"null" json:"hash"`
	Status         bool         `gorm:"false" json:"status"`
	CreatedAt      time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt      *time.Time   `sql:"index" json:"deleted_at"`
	Subscription   Subscription `json:"subscription"`
}

// SaveLicense is a ...
func (l *License) SaveLicense() (*License, error) {
	err := config.DB.Create(&l).Error
	if err != nil {
		return &License{}, err
	}
	return l, nil
}

// FindLicenseByID is a ...
func (l *License) FindLicenseByID(uid uint64) (*License, error) {
	err := config.DB.Model(License{}).Where("id = ?", uid).Take(&l).Error
	if err != nil {
		return &License{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &License{}, ErrKeyNotFound
	}
	return l, err
}

// FindLicense is a ...
func (l *License) FindLicense(key []byte) (*License, error) {
	err := config.DB.Model(License{}).Where("license = ?", key).Take(&l).Error
	if err != nil {
		return &License{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &License{}, ErrLicenseNotFound
	}
	return l, err
}

// UpdateLicense is a ...
func (l *License) UpdateLicense(uid uint64) (*License, error) {
	db := config.DB.Model(&License{}).Where("id = ?", uid).Take(&License{}).UpdateColumns(
		map[string]interface{}{
			"subscription_id": l.SubscriptionID,
			"license":         l.License,
			"status":          l.Status,
			"update_at":       time.Now(),
		},
	)
	if db.Error != nil {
		return &License{}, db.Error
	}

	err := config.DB.Model(&License{}).Where("id = ?", uid).Take(&l).Error
	if err != nil {
		return &License{}, err
	}
	return l, nil
}

// DeleteLicense is a ...
func (l *License) DeleteLicense(uid uint64) (int64, error) {
	db := config.DB.Model(&License{}).Where("id = ?", uid).Take(&License{}).Delete(&License{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeactivateLicenseBySubID is a ...
func DeactivateLicenseBySubID(uid uint64) error {
	db := config.DB.Model(&License{}).Where("subscription_id = ?", uid).UpdateColumns(
		map[string]interface{}{
			"status":     false,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// LicensesListBySubscriptionID is a ...
func LicensesListBySubscriptionID(uid uint64) *[]License {
	licenses := []License{}
	db := config.DB.Where("subscription_id = ?", uid).Order("created_at DESC").Find(&licenses)
	if db.Error != nil {
		return &licenses
	}
	return &licenses
}
