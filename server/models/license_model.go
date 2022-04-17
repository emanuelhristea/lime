package models

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/emanuelhristea/lime/config"
	"github.com/emanuelhristea/lime/license"
	"github.com/jinzhu/gorm"
)

// License is a ...
type License struct {
	ID             uint64       `gorm:"primary_key;auto_increment" json:"id"`
	SubscriptionID uint64       `sql:"unique_index:idx_sub_mac;type:int REFERENCES subscriptions(id) ON DELETE CASCADE" json:"subscription_id"`
	Mac            string       `gorm:"unique_index:idx_sub_mac;size:255;not null;unique" json:"mac"`
	License        []byte       `gorm:"null" json:"license"`
	Hash           string       `gorm:"null" json:"hash"`
	Status         bool         `gorm:"false" json:"status"`
	CreatedAt      time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
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
func (l *License) FindLicenseByID(uid uint64, relations ...string) (*License, error) {
	db := config.DB.Model(License{}).Where("id = ?", uid)
	for _, rel := range relations {
		db = db.Preload(rel, func(db *gorm.DB) *gorm.DB {
			return db.Order("ID asc")
		})
	}
	err := db.Take(&l).Error
	if err != nil {
		return &License{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &License{}, ErrKeyNotFound
	}
	return l, err
}

// FindLicenseByMac is a ...
func (l *License) FindLicenseByMac(mac string, relations ...string) (*License, error) {
	db := config.DB.Model(License{}).Where("mac = ?", mac)
	for _, rel := range relations {
		db = db.Preload(rel, func(db *gorm.DB) *gorm.DB {
			return db.Order("ID asc")
		})
	}
	err := db.Take(&l).Error
	if err != nil {
		return &License{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &License{}, ErrKeyNotFound
	}
	return l, err
}

// FindLicense is a ...
func (l *License) FindLicense(key []byte, mac string, relations ...string) (*License, error) {
	db := config.DB.Model(License{}).Where("license = ?", key).Where("mac = ?", mac)
	for _, rel := range relations {
		db = db.Preload(rel, func(db *gorm.DB) *gorm.DB {
			return db.Order("ID asc")
		})
	}
	err := db.Take(&l).Error
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
			"mac":             l.Mac,
			"license":         l.License,
			"status":          l.Status,
			"hash":            l.Hash,
			"updated_at":      time.Now(),
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

func (l *License) UpdateLicenseExpiry(expiry time.Time) (*License, error) {
	existing, err := license.Decode([]byte(l.License), license.GetPublicKey())
	if err != nil {
		return l, err
	}

	_license := &license.License{
		Iss: existing.Iss,
		Cus: existing.Cus,
		Sub: existing.Sub,
		Typ: existing.Typ,
		Lim: existing.Lim,
		Dat: existing.Dat,
		Exp: expiry,
		Iat: existing.Iat,
	}

	encoded, err := _license.Encode(license.GetPrivateKey())
	if err != nil {
		return l, err
	}

	l.License = encoded

	hash := md5.Sum([]byte(encoded))
	licenseHash := hex.EncodeToString(hash[:])
	l.Hash = licenseHash
	return l.UpdateLicense(l.ID)
}

// DeleteLicense is a ...
func DeleteLicense(uid uint64) (int64, error) {
	db := config.DB.Model(&License{}).Where("id = ?", uid).Take(&License{}).Delete(&License{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// SetLicenseStatusBySubID is a ...
func SetLicenseStatusBySubID(uid uint64, status bool) error {
	db := config.DB.Model(&License{}).Where("subscription_id = ?", uid).UpdateColumns(
		map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// SetLicenseStatusBySubID is a ...
func SetLicenseStatusByMac(mac string, status bool) error {
	db := config.DB.Model(&License{}).Where("mac = ?", mac).UpdateColumns(
		map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

// DeactivateLicenseBySubID is a ...
func SetLicenseStatusByID(uid uint64, status bool) error {
	db := config.DB.Model(&License{}).Where("id = ?", uid).UpdateColumns(
		map[string]interface{}{
			"status":     status,
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

// SubscriptionsList is a ...
func LicensesList(subscriptionID string, relations ...string) *[]License {
	db := config.DB.Model(&License{}).Where("subscription_id=?", subscriptionID)
	for _, rel := range relations {
		db = db.Preload(rel, func(db *gorm.DB) *gorm.DB {
			return db.Order("ID asc")
		})
	}

	licenses := []License{}
	db = db.Find(&licenses).Order("ID asc")
	if db.Error != nil {
		return &licenses
	}
	return &licenses
}
