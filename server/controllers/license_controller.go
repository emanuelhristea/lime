package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/emanuelhristea/lime/license"
	"github.com/emanuelhristea/lime/server/models"
	"github.com/gin-gonic/gin"
)

// GetSubscriptionList is a ...
func GetLicensesList(c *gin.Context) {
	subscriptionId := c.Param("subscriptionId")
	preload, exists := c.GetQuery("load")
	licenseList := &[]models.License{}
	if exists {
		licenseList = models.LicensesList(subscriptionId, preload)
	} else {
		licenseList = models.LicensesList(subscriptionId)
	}
	respondJSON(c, http.StatusOK, licenseList)
}

// GetSubscription is a ...
func GetLicense(c *gin.Context) {
	id := c.Param("id")
	licenseId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	modelLicense := models.License{}
	_license, err := modelLicense.FindLicenseByID(licenseId)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, _license)
}

func UpdateLicense(c *gin.Context) {
	id := c.Param("id")
	licenseId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	_license := &models.License{}
	_license, err = _license.FindLicenseByID(licenseId)

	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	status := false
	if c.PostForm("status") != "" {
		status = true
	}

	err = models.SetLicenseStatusByID(licenseId, status)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, _license)
}

func DeleteLicense(c *gin.Context) {
	id := c.Param("id")
	licenseId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	rows, err := models.DeleteLicense(licenseId)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Cannot delete customer that has active subscriptions")
		return
	}

	respondJSON(c, http.StatusOK, fmt.Sprintf("%d", rows))
}

// VerifyKey is a ...
// @Accept application/json
// @Produce application/json
// @Param
// @Success 200 {string} string "{"status":"200", "msg":""}"
// @Router /verify [post]
func VerifyKey(c *gin.Context) {
	modelLicense := models.License{}

	request := &requestLicense{}
	c.BindJSON(&request)

	licenseKey, err := base64.StdEncoding.DecodeString(request.License)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	_license, err := modelLicense.FindLicense(licenseKey)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	l, err := license.Decode([]byte(licenseKey), license.GetPublicKey())
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	if l.Expired() {
		models.SetLicenseStatusBySubID(_license.ID, false)
	}

	if !_license.Status {
		respondJSON(c, http.StatusNotFound, "License expired!")
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(l)

	respondJSON(c, http.StatusOK, "Active")
}

// CreateKey is a ...
// @Accept application/json
// @Produce application/json
// @Param
// @Success 200 {string} string "{"status":"200", "msg":""}"
// @Router /key [post]
func CreateKey(c *gin.Context) {
	month := time.Hour * 24 * 31
	modelSubscription := models.Subscription{}

	request := &requestLicense{}
	c.BindJSON(&request)

	_subscription, err := modelSubscription.FindSubscriptionByStripeID(request.StripeID, "Customer", "Tariff", "Licenses")
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	if _subscription.ID == 0 {
		respondJSON(c, http.StatusNotFound, "Customers not found!")
		return
	}

	encoded, response := addLicenseToSubscription(_subscription, month)
	if response != "" {
		respondJSON(c, http.StatusMethodNotAllowed, response)
		return
	}

	respondJSON(c, http.StatusOK, base64.StdEncoding.EncodeToString([]byte(encoded)))
}

func numberOfActiveLicenses(_subscription *models.Subscription) int {
	num := 0
	for _, lic := range _subscription.Licenses {
		l, err := license.Decode([]byte(lic.License), license.GetPublicKey())
		if err != nil || !lic.Status || l.Expired() {
			continue
		}
		num++
	}
	return num
}

func addLicenseToSubscription(_subscription *models.Subscription, expiry time.Duration) ([]byte, string) {
	if !_subscription.Customer.Status || !_subscription.Status {
		return nil, "Your subscription was deactivated!"
	}
	if numberOfActiveLicenses(_subscription) >= _subscription.Tariff.Users {
		return nil, "Your have reached the maximum number of users for your subscription!"
	}

	limit := license.Limits{
		Tandem:  _subscription.Tariff.Tandem,
		Triaxis: _subscription.Tariff.Triaxis,
		Robots:  _subscription.Tariff.Robots,
		Users:   _subscription.Tariff.Users,
	}

	metadata := []byte(`{"message": "test message"}`)
	_license := &license.License{
		Iss: _subscription.Customer.Name,
		Cus: _subscription.StripeID,
		Sub: _subscription.TariffID,
		Typ: _subscription.Tariff.Name,
		Lim: limit,
		Dat: metadata,
		Exp: time.Now().UTC().Add(expiry),
		Iat: time.Now().UTC(),
	}

	encoded, err := _license.Encode(license.GetPrivateKey())
	if err != nil {
		return nil, err.Error()
	}

	hash := md5.Sum([]byte(encoded))
	licenseHash := hex.EncodeToString(hash[:])

	key := &models.License{
		SubscriptionID: _subscription.ID,
		License:        encoded,
		Hash:           licenseHash,
		Status:         true,
	}

	_, err = key.SaveLicense()
	if err != nil {
		return nil, err.Error()
	}
	return encoded, ""
}

// GetKey is a ...
// @Accept application/json
// @Produce application/json
// @Param
// @Success 200 {string} string "{"status":"200", "msg":""}"
// @Router /key/:customer_id [get]
func GetKey(c *gin.Context) {
	respondJSON(c, http.StatusOK, "GetKey")
}

// GetUserSubscriptions is a ...
func GetUserSubscriptions(c *gin.Context) {

}

// UpdateKey is a ...
// @accept application/json
// @Produce application/json
// @Param
// @Success 200 {string} string "{"status":"200", "msg":""}"
// @Router /key/:customer_id [PATCH]
func UpdateKey(c *gin.Context) {
	respondJSON(c, http.StatusOK, "UpdateKey")
}

// UpdateKey is a ...
// @accept application/json
// @Produce application/json
// @Param
// @Success 200 {string} string "{"status":"200", "msg":""}"
// @Router /key/:customer_id [DELETE]
func DeleteKey(c *gin.Context) {
	respondJSON(c, http.StatusOK, "UpdateKey")
}
