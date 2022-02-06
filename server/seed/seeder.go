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
		Users:   1,
	},
	{
		Name:    "Tandem",
		Price:   2500,
		Tandem:  true,
		Triaxis: false,
		Robots:  false,
		Users:   20,
	},
	{
		Name:    "Triaxis",
		Price:   2500,
		Tandem:  false,
		Triaxis: true,
		Robots:  false,
		Users:   20,
	},
}

var customers = []models.Customer{
	{
		ID:     1,
		Name:   "Andrei Oana",
		Status: true,
	},
	{
		ID:     2,
		Name:   "Emanuel Hristea",
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
	// err := db.DropTableIfExists(&models.Tariff{}, &models.Customer{}, &models.Subscription{}, &models.License{}).Error
	// if err != nil {
	// 	log.Fatalf("cannot drop table: %v", err)
	// }
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
