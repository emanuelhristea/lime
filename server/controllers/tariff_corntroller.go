package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/emanuelhristea/lime/server/models"
	"github.com/gin-gonic/gin"
)

// GetTariffList is a ...
func GetTariffList(c *gin.Context) {
	preload, exists := c.GetQuery("load")

	tariffsList := &[]models.Tariff{}
	if exists {
		tariffsList = models.TariffsList(strings.Split(preload, ",")...)
	} else {
		tariffsList = models.TariffsList()
	}
	respondJSON(c, http.StatusOK, tariffsList)
}

// GetTariff is a ...
func GetTariff(c *gin.Context) {
	id := c.Param("id")

	tariffId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	modelTariff := models.Tariff{}
	_tariff, err := modelTariff.FindTariffByID(tariffId)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, _tariff)
}

func getTariffFromForm(c *gin.Context) (*models.Tariff, bool) {
	n := c.PostForm("name")
	if n == "" {
		respondJSON(c, http.StatusBadRequest, "Name is invalid")
		return nil, true
	}

	p := c.PostForm("price")
	pr, err := strconv.ParseFloat(p, 64)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Price format is invalid")
		return nil, true
	}

	price := int(pr * 100)
	if price < 0 {
		respondJSON(c, http.StatusBadRequest, "Price cannot be negative")
		return nil, true
	}

	tandem := false
	if c.PostForm("tandem") != "" {
		tandem = true
	}

	triaxis := false
	if c.PostForm("triaxis") != "" {
		triaxis = true
	}

	robots := false
	if c.PostForm("robots") != "" {
		robots = true
	}

	u := c.PostForm("users")
	users, err := strconv.ParseInt(u, 10, 64)
	if err != nil || users < 1 || users > 100 {
		respondJSON(c, http.StatusBadRequest, "Number of users is invalid")
		return nil, true
	}

	e := c.PostForm("period")
	period, err := strconv.ParseInt(e, 10, 64)
	if err != nil || period < 1 || period > 1000 {
		respondJSON(c, http.StatusBadRequest, "License period is invalid")
		return nil, true
	}

	modelTariff := &models.Tariff{
		Name:    n,
		Price:   price,
		Tandem:  tandem,
		Triaxis: triaxis,
		Robots:  robots,
		Users:   int(users),
		Period:  int(period),
	}
	return modelTariff, false
}

// CreateTariff is a ...
func CreateTariff(c *gin.Context) {
	modelTariff, shouldReturn := getTariffFromForm(c)
	if shouldReturn {
		return
	}

	_tariff, err := modelTariff.SaveTariff()
	if err != nil {
		respondJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, _tariff)
}

func UpdateTariff(c *gin.Context) {
	id := c.Param("id")
	tariffId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	_existing := &models.Tariff{}
	_existing, err = _existing.FindTariffByID(tariffId, "Subscriptions", "Subscriptions.Licenses")
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	modelTariff, shouldReturn := getTariffFromForm(c)
	if shouldReturn {
		return
	}

	_tariff, err := modelTariff.UpdateTariff(tariffId)
	if err != nil {
		respondJSON(c, http.StatusConflict, err.Error())
		return
	}

	//Update expiry of all existing subscriptions if the period is modified
	if _existing.Period != modelTariff.Period {
		for _, sub := range _existing.Subscriptions {
			expiry := time.Duration(_tariff.Period) * 24 * time.Hour
			sub.ExpiresAt = sub.IssuedAt.Add(expiry)
			sub, err := sub.UpdateSubscription(sub.ID)
			if err != nil {
				continue
			}

			for _, lic := range sub.Licenses {
				_, err := lic.UpdateLicenseExpiry(sub.ExpiresAt)
				if err != nil {
					log.Print(err.Error())
					continue
				}
			}
		}
	}

	respondJSON(c, http.StatusOK, _tariff)
}

func DeleteTariff(c *gin.Context) {
	id := c.Param("id")
	tariffId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	rows, err := models.DeleteTariff(tariffId)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Cannot delete plan that is in use")
		return
	}

	respondJSON(c, http.StatusOK, fmt.Sprintf("%d", rows))
}
