package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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
		licenseList = models.LicensesList(subscriptionId, strings.Split(preload, ",")...)
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

func CreateLicense(c *gin.Context) {
	sId := c.Param("subscripionId")
	subscriptionId, err := strconv.ParseUint(sId, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	mac := c.PostForm("mac")
	modelSubscription := models.Subscription{}
	_subscription, err := modelSubscription.FindSubscriptionByID(subscriptionId, "Customer", "Tariff", "Licenses")
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	if _subscription.ID == 0 {
		respondJSON(c, http.StatusNotFound, "Subscripton not found!")
		return
	}

	status := false
	if c.PostForm("status") != "" {
		status = true
	}

	key, response := addLicenseToSubscription(_subscription, mac, status)
	if response != "" {
		respondJSON(c, http.StatusMethodNotAllowed, response)
		return
	}

	license := license.Key{
		Key:   base64.StdEncoding.EncodeToString([]byte(key.License)),
		Hash:  key.Hash,
		Mac:   key.Mac,
		Valid: key.Status,
	}

	respondJSON(c, http.StatusOK, license)
}

func UpdateLicense(c *gin.Context) {
	id := c.Param("id")
	licenseId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	status := false
	if c.PostForm("status") != "" {
		status = true
	}

	_license := &models.License{}
	_license, err = _license.FindLicenseByID(licenseId, "Subscription", "Subscription.Tariff", "Subscription.Licenses")
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	nb := numberOfActiveLicenses(&_license.Subscription)
	if status && nb >= _license.Subscription.Tariff.Users {
		respondJSON(c, http.StatusNotFound, "Your have reached the maximum number of users for your subscription!")
		return
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

	match, err := regexp.MatchString("^([0-9A-F]{2}[:-]){5}([0-9A-F]{2})$", request.Mac)
	if !match || err != nil {
		respondJSON(c, http.StatusNotFound, "The MAC address is invalid!")
		return
	}

	_license, err := modelLicense.FindLicense(licenseKey, request.Mac, "Subscription")
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	l, err := license.Decode([]byte(licenseKey), license.GetPublicKey())
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	sub := &_license.Subscription
	if sub.Expired() {
		if sub.Status {
			sub.Status = false
			sub.UpdateSubscription(sub.ID)
		}

		models.SetLicenseStatusBySubID(_license.ID, false)
		respondJSON(c, http.StatusNotFound, "License expired")
		return
	}

	if !_license.Status {
		respondJSON(c, http.StatusNotFound, "License inactive")
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
	modelSubscription := models.Subscription{}
	request := &requestLicense{}
	c.BindJSON(&request)

	_subscription, err := modelSubscription.FindSubscriptionByStripeID(request.StripeID, "Customer", "Tariff", "Licenses")
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	if _subscription.ID == 0 {
		respondJSON(c, http.StatusNotFound, "Subscription not found!")
		return
	}

	key, response := addLicenseToSubscription(_subscription, request.Mac, true)
	if response != "" {
		if strings.Contains(response, "licenses_mac_key") {
			response = "The mac address: " + request.Mac + " is already registered. Please realease it first!"
		}
		respondJSON(c, http.StatusMethodNotAllowed, response)
		return
	}
	license := license.Key{
		Key:   base64.StdEncoding.EncodeToString([]byte(key.License)),
		Hash:  key.Hash,
		Mac:   key.Mac,
		Valid: key.Status,
	}

	respondJSON(c, http.StatusOK, license)
}

func ReleaseKey(c *gin.Context) {
	modelLicense := models.License{}

	request := &requestLicense{}
	c.BindJSON(&request)

	licenseKey, err := base64.StdEncoding.DecodeString(request.License)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	match, err := regexp.MatchString("^([0-9A-F]{2}[:-]){5}([0-9A-F]{2})$", request.Mac)
	if !match || err != nil {
		respondJSON(c, http.StatusNotFound, "The MAC address is invalid!")
		return
	}

	_license, err := modelLicense.FindLicense(licenseKey, request.Mac)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	nb, err := models.DeleteLicense(_license.ID)
	if err != nil {
		respondJSON(c, http.StatusConflict, err.Error())
	}
	respondJSON(c, http.StatusOK, nb)
}

func numberOfActiveLicenses(_subscription *models.Subscription) int {
	num := 0
	if _subscription.Expired() {
		return num
	}

	for _, lic := range _subscription.Licenses {
		if !lic.Status {
			continue
		}
		num++
	}
	return num
}

func addLicenseToSubscription(_subscription *models.Subscription, mac string, status bool) (*models.License, string) {
	match, err := regexp.MatchString("^([0-9A-F]{2}[:-]){5}([0-9A-F]{2})$", mac)
	if !match || err != nil {
		return &models.License{}, "The MAC address is invalid!"
	}

	if !_subscription.Customer.Status || !_subscription.Status {
		return &models.License{}, "Your subscription was deactivated!"
	}
	if numberOfActiveLicenses(_subscription) >= _subscription.Tariff.Users {
		return &models.License{}, "Your have reached the maximum number of users for your subscription!"
	}

	limit := license.Limits{
		Tandem:  _subscription.Tariff.Tandem,
		Triaxis: _subscription.Tariff.Triaxis,
		Robots:  _subscription.Tariff.Robots,
		Period:  _subscription.Tariff.Period,
		Devices: _subscription.Tariff.Users,
	}

	metadata := []byte(`{"message": "test message"}`)
	expiry := time.Duration(limit.Period) * 24 * time.Hour

	//Activate subscription on first license use
	if len(_subscription.Licenses) == 0 && _subscription.CreatedAt.After(_subscription.IssuedAt) {
		_subscription.IssuedAt = time.Now().UTC()
		_subscription.ExpiresAt = time.Now().UTC().Add(expiry)
		_subscription.UpdateSubscription(_subscription.ID)
	}
	_license := &license.License{
		Iss: _subscription.Customer.Name,
		Cus: _subscription.StripeID,
		Sub: _subscription.TariffID,
		Typ: _subscription.Tariff.Name,
		Lim: limit,
		Dat: metadata,
		Exp: _subscription.ExpiresAt,
		Iat: time.Now().UTC(),
	}

	encoded, err := _license.Encode(license.GetPrivateKey())
	if err != nil {
		return &models.License{}, err.Error()
	}

	hash := md5.Sum([]byte(encoded))
	licenseHash := hex.EncodeToString(hash[:])

	key := &models.License{
		SubscriptionID: _subscription.ID,
		License:        encoded,
		Mac:            mac,
		Hash:           licenseHash,
		Status:         status,
	}

	_, err = key.SaveLicense()
	if err != nil {
		return &models.License{}, err.Error()
	}
	return key, ""
}

// GetUserSubscriptions is a ...
func GetUserSubscriptions(c *gin.Context) {
	request := &requestSubscriptions{}
	c.BindJSON(request)

	match, err := regexp.MatchString("^([0-9A-F]{2}[:-]){5}([0-9A-F]{2})$", request.Mac)
	if !match || err != nil {
		respondJSON(c, http.StatusNotFound, "The MAC address is invalid!")
		return
	}

	modelCustomer := models.Customer{}
	_customer, err := modelCustomer.FindCustomerByEmail(request.Email, "Subscriptions", "Subscriptions.Tariff", "Subscriptions.Licenses")

	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	response := []license.Subscription{}
	for _, sub := range _customer.Subscriptions {
		if !sub.Status {
			continue
		}

		if sub.Expired() && sub.Status {
			sub.Status = false
			sub.UpdateSubscription(sub.ID)
		}

		active := numberOfActiveLicenses(&sub)
		expiresIn := int(sub.ExpiresAt.Sub(time.Now().UTC()).Hours() / 24)

		temp := license.Subscription{
			Plan:       sub.Tariff.Name,
			PurchaseID: sub.StripeID,
			Limits: license.Limits{
				Tandem:  sub.Tariff.Tandem,
				Triaxis: sub.Tariff.Triaxis,
				Robots:  sub.Tariff.Robots,
				Period:  sub.Tariff.Period,
				Devices: sub.Tariff.Users,
			},

			InUse:     active,
			Status:    sub.Status,
			ExpiresIn: expiresIn,
			Role:      _customer.Role,
		}

		for _, lic := range sub.Licenses {
			if lic.Mac == request.Mac {
				if err != nil {
					respondJSON(c, http.StatusBadRequest, err.Error())
					return
				}
				temp.LicenseKey = license.Key{
					Key:   base64.StdEncoding.EncodeToString(lic.License),
					Hash:  lic.Hash,
					Valid: lic.Status,
					Mac:   lic.Mac,
				}
				break
			}
		}

		response = append(response, temp)
	}

	respondJSON(c, http.StatusOK, response)
}
