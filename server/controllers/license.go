package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/emanuelhristea/lime/license"
	"github.com/emanuelhristea/lime/server/models"
	"github.com/gin-gonic/gin"
)

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
	log.Print(licenseKey)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	_license, err := modelLicense.FindLicense(licenseKey)
	log.Print(_license)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	if _license.ID == 0 {
		respondJSON(c, http.StatusNotFound, "License not found!")
		return
	}

	l, err := license.Decode([]byte(licenseKey), license.GetPublicKey())
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
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
	modelTariff := models.Tariff{}
	modelCustomer := models.Customer{}

	request := &requestLicense{}
	c.BindJSON(&request)

	_subscription, err := modelSubscription.FindSubscriptionByStripeID(request.StripeID)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	if _subscription.ID == 0 {
		respondJSON(c, http.StatusNotFound, "Customers not found!")
		return
	}

	_customer, _ := modelCustomer.FindCustomerByID(_subscription.CustomerID)

	_tariff, err := modelTariff.FindTariffByID(_subscription.TariffID)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	if _tariff.ID == 0 {
		respondJSON(c, http.StatusNotFound, "Tariff not found!")
		return
	}

	limit := license.Limits{
		Tandem:  _tariff.Tandem,
		Triaxis: _tariff.Triaxis,
		Robots:  _tariff.Robots,
		Users:   _tariff.Users,
	}
	metadata := []byte(`{"message": "test message"}`)
	_license := &license.License{
		Iss: _customer.Name,
		Cus: _subscription.StripeID,
		Sub: _subscription.TariffID,
		Typ: _tariff.Name,
		Lim: limit,
		Dat: metadata,
		Exp: time.Now().UTC().Add(month),
		Iat: time.Now().UTC(),
	}

	encoded, err := _license.Encode(license.GetPrivateKey())
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	models.DeactivateLicenseBySubID(_subscription.ID)

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
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, base64.StdEncoding.EncodeToString([]byte(encoded)))
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

// UpdateKey is a ...
// @accept application/json
// @Produce application/json
// @Param
// @Success 200 {string} string "{"status":"200", "msg":""}"
// @Router /key/:customer_id [PATCH]
func UpdateKey(c *gin.Context) {
	respondJSON(c, http.StatusOK, "UpdateKey")
}

func GetTariffList(c *gin.Context) {
	preload, exists := c.GetQuery("load")
	log.Print(preload)
	tariffsList := &[]models.Tariff{}
	if exists {
		tariffsList = models.TariffsList(preload)
	} else {
		tariffsList = models.TariffsList()
	}
	c.JSON(http.StatusOK, tariffsList)
}

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
	c.JSON(http.StatusOK, _tariff)
}

func CreateTariff(c *gin.Context) {
	n := c.PostForm("name")
	if n == "" {
		respondJSON(c, http.StatusBadRequest, "Name is invalid")
		return
	}

	p := c.PostForm("price")
	pr, err := strconv.ParseFloat(p, 64)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Price format is invalid")
		return
	}

	price := int(pr * 100)
	if price < 0 {
		respondJSON(c, http.StatusBadRequest, "Price cannot be negative")
		return
	}

	tandem := false
	if c.PostForm("tandem") == "on" {
		tandem = true
	}

	triaxis := false
	if c.PostForm("triaxis") == "on" {
		triaxis = true
	}

	robots := false
	if c.PostForm("robots") == "on" {
		robots = true
	}

	u := c.PostForm("users")
	users, err := strconv.ParseInt(u, 10, 64)
	if err != nil || users < 1 || users > 100 {
		respondJSON(c, http.StatusBadRequest, "Number of users is invalid")
		return
	}

	modelTariff := &models.Tariff{
		Name:    n,
		Price:   price,
		Tandem:  tandem,
		Triaxis: triaxis,
		Robots:  robots,
		Users:   int(users),
	}

	_tariff, err := modelTariff.SaveTariff()
	if err != nil {
		respondJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, _tariff.Name)
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

func CreateCustomer(c *gin.Context) {
	n := c.PostForm("name")
	if n == "" {
		respondJSON(c, http.StatusBadRequest, "Name is invalid")
		return
	}

	status := false
	if c.PostForm("status") == "on" {
		status = true
	}

	modelCustomer := &models.Customer{
		Name:   n,
		Status: status,
	}

	_customer, err := modelCustomer.SaveCustomer()
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Cannot save customer, probably duplicate name")
		return
	}

	respondJSON(c, http.StatusOK, _customer.Name)
}

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	customerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	n := c.PostForm("name")
	if n == "" {
		respondJSON(c, http.StatusBadRequest, "Customer not found")
		return
	}

	status := false
	if c.PostForm("status") == "on" {
		status = true
	}
	_customer := models.Customer{
		Name:   n,
		Status: status,
	}

	_, err = _customer.UpdateCustomer(customerId)

	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, _customer.Name)
}

func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	customerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	rows, err := models.DeleteCustomer(customerId)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Cannot delete customer that has active subscriptions")
		return
	}

	respondJSON(c, http.StatusOK, fmt.Sprintf("%d", rows))
}

func CreateSubscription(c *gin.Context) {
	stripe := c.PostForm("stripe_id")
	if stripe == "" {
		respondJSON(c, http.StatusBadRequest, "Name is invalid")
		return
	}

	tariff, err := strconv.ParseUint(c.PostForm("tariff_id"), 10, 64)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Invalid plan selected")
		return
	}

	status := false
	if c.PostForm("status") == "on" {
		status = true
	}

	id := c.Param("id")
	customerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	modelSubscription := &models.Subscription{
		StripeID:   stripe,
		TariffID:   tariff,
		CustomerID: customerId,
		Status:     status,
	}

	_subscription, err := modelSubscription.SaveSubscription()
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Cannot save subscription. Duplicate payment id or plan")
		return
	}

	respondJSON(c, http.StatusOK, _subscription.StripeID)
}

func UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	customerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	n := c.PostForm("name")
	if n == "" {
		respondJSON(c, http.StatusBadRequest, "Customer not found")
		return
	}

	status := false
	if c.PostForm("status") == "on" {
		status = true
	}
	_customer := models.Customer{
		Name:   n,
		Status: status,
	}

	_, err = _customer.UpdateCustomer(customerId)

	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, _customer.Name)
}

func DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	customerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	rows, err := models.DeleteCustomer(customerId)
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Cannot delete customer that has active subscriptions")
		return
	}

	respondJSON(c, http.StatusOK, fmt.Sprintf("%d", rows))
}
