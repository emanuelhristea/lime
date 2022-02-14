package seed

import (
	"log"

	"github.com/emanuelhristea/lime/server/models"
	"github.com/jinzhu/gorm"
)

var tariffs = []models.Tariff{
	{
		Name:    "Trial",
		Price:   0,
		Tandem:  true,
		Triaxis: false,
		Robots:  true,
		Period:  30,
		Users:   1,
	},
	{
		Name:    "Tandem",
		Price:   2500,
		Tandem:  true,
		Triaxis: false,
		Robots:  false,
		Period:  365,
		Users:   20,
	},
	{
		Name:    "Triaxis",
		Price:   2500,
		Tandem:  false,
		Triaxis: true,
		Robots:  false,
		Period:  365,
		Users:   20,
	},
}

var customers = []models.Customer{
	{
		ID:     1,
		Name:   "Andrei Oana",
		Email:  "aoana@destaco.com",
		Role:   "admin",
		Status: true,
	},
	{
		ID:     2,
		Name:   "Emanuel Hristea",
		Email:  "c-ehristea@destaco.com",
		Role:   "user",
		Status: true,
	},
}

var subscription = []models.Subscription{
	{
		CustomerID: 1,
		StripeID:   "cus_FEDaLVeqQoVy6m",
		TariffID:   1,
		Status:     true,
	},
	{
		CustomerID: 2,
		StripeID:   "cus_APBaLDeqQoVy8m",
		TariffID:   2,
		Status:     true,
	},
}

// Load import test data to database
func Load(db *gorm.DB) {
	db.DropTableIfExists(&models.License{})
	db.DropTableIfExists(&models.Subscription{})
	db.DropTableIfExists(&models.Customer{})
	db.DropTableIfExists(&models.Tariff{})

	result := db.Exec("SELECT 1 FROM pg_type WHERE typname = 'role';")

	switch {
	case result.RowsAffected == 0:
		err := db.Exec("CREATE TYPE role AS ENUM ('admin', 'user', 'guest');").Error
		if err != nil {
			log.Fatal("Error creating role ENUM")
			return
		}

	case result.Error != nil:
		log.Fatal("Error selecting role type")
		return
	}

	err := db.AutoMigrate(&models.Tariff{}, &models.Customer{}, &models.Subscription{}, &models.License{}).Error

	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i := range tariffs {
		err = db.Model(&models.Tariff{}).Create(&tariffs[i]).Error
		if err != nil {
			log.Fatalf("cannot seed tariff table: %v", err)
		}
	}

	for i := range customers {
		err = db.Model(&models.Customer{}).Create(&customers[i]).Error
		if err != nil {
			log.Fatalf("cannot seed customer table: %v", err)
		}
	}

	for i := range subscription {
		err = db.Model(&models.Subscription{}).Create(&subscription[i]).Error
		if err != nil {
			log.Fatalf("cannot seed subscription table: %v", err)
		}
	}
}
