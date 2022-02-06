package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/emanuelhristea/lime/license"
	"github.com/emanuelhristea/lime/server/middleware"
	"github.com/emanuelhristea/lime/server/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// MainHandler is a ...
func MainHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(middleware.IdentityKey)

	if user == nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Authorization",
		})
	} else {
		action := c.Param("action")

		switch action {
		case "/new":
			c.HTML(http.StatusOK, "new_customer.html", gin.H{
				"title": "Create new customer",
			})
		default:
			log.Print("serve main")
			customersList := models.CustomersList()
			c.HTML(http.StatusOK, "customers.html", gin.H{
				"title":     "Customers",
				"customers": customersList,
			})
		}
	}
}

// CustomerSubscriptionList is a ...
func CustomerSubscriptionList(c *gin.Context) {
	customerID := c.Param("id")
	action := c.Param("action")
	subscriptionsList := models.SubscriptionsByCustomerID(customerID)
	licensesList := models.LicensesListBySubscriptionID(customerID)

	switch action {
	case "/":
	case "/new":

		if len(*subscriptionsList) > 0 {
			month := time.Hour * 24 * 31
			modelTariff := models.Tariff{}

			_tariff, _ := modelTariff.FindTariffByID((*subscriptionsList)[0].ID)

			limit := license.Limits{
				Tandem:  _tariff.Tandem,
				Triaxis: _tariff.Triaxis,
				Robots:  _tariff.Robots,
				Users:   _tariff.Users,
			}
			metadata := []byte(`{"message": "test message"}`)
			_license := &license.License{
				Iss: (*subscriptionsList)[0].CustomerName,
				Cus: (*subscriptionsList)[0].StripeID,
				Sub: (*subscriptionsList)[0].TariffID,
				Typ: _tariff.Name,
				Lim: limit,
				Dat: metadata,
				Exp: time.Now().UTC().Add(month),
				Iat: time.Now().UTC(),
			}
			encoded, _ := _license.Encode(license.GetPrivateKey())

			hash := md5.Sum([]byte(encoded))
			licenseHash := hex.EncodeToString(hash[:])

			models.DeactivateLicenseBySubID((*subscriptionsList)[0].ID)
			key := &models.License{
				SubscriptionID: (*subscriptionsList)[0].ID,
				License:        encoded,
				Hash:           licenseHash,
				Status:         true,
			}
			key.SaveLicense()

			c.Redirect(http.StatusFound, "/admin/subscription/"+customerID+"/")
		}
	default:
		c.Redirect(http.StatusFound, "/admin/subscription/"+customerID+"/")
	}

	name := ""
	if len(*subscriptionsList) > 0 {
		name = (*subscriptionsList)[0].CustomerName
	} else {
		customer := &models.Customer{}
		_id, err := strconv.ParseUint(customerID, 10, 32)
		if err == nil {
			_customer, err := customer.FindCustomerByID(uint32(_id))
			if err == nil {
				name = _customer.Name
			}
		}
	}

	c.HTML(http.StatusOK, "subscriptions.html", gin.H{
		"title":         "Subscription and Licenses for " + name,
		"customerID":    customerID,
		"Subscriptions": subscriptionsList,
		"Licenses":      licensesList,
	})

}

// DownloadLicense is a ...
func DownloadLicense(c *gin.Context) {
	licenseID := c.Param("id")
	license := models.License{}

	uid, err := strconv.ParseUint(licenseID, 10, 32)
	if err != nil {
		panic(err)
	}

	license.FindLicenseByID(uint32(uid))

	body := string(license.License)
	reader := strings.NewReader(body)
	contentLength := int64(len(body))
	contentType := "application/octet-stream"
	extraHeaders := map[string]string{"Content-Disposition": `attachment; filename="` + license.Hash + `"`}
	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

func TariffsList(c *gin.Context) {
	action := c.Param("action")

	switch action {
	case "/":
		tariffsList := models.TariffsList()
		c.HTML(http.StatusOK, "tariffs.html", gin.H{
			"title":   "Pricing",
			"Tariffs": tariffsList,
		})
	case "/new":
		c.HTML(http.StatusOK, "new_tariff.html", gin.H{
			"title": "Create new pricing",
		})
	}
}

func TariffAction(c *gin.Context) {
	id := c.Param("id")
	action := c.Param("action")
	switch action {
	case "/":
	case "/delete":
		tariffId, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			panic(err)
		}
		models.DeleteTariff(uint32(tariffId))
		tariffsList := models.TariffsList()
		c.HTML(http.StatusOK, "tariffs.html", gin.H{
			"title":   "Pricing",
			"Tariffs": tariffsList,
		})
	}
}
