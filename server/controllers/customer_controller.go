package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/emanuelhristea/lime/server/models"
	"github.com/gin-gonic/gin"
)

// GetSubscriptionList is a ...
func GetCustomerList(c *gin.Context) {
	preload, exists := c.GetQuery("load")
	customerList := &[]models.Customer{}
	if exists {
		customerList = models.CustomersList(strings.Split(preload, ",")...)
	} else {
		customerList = models.CustomersList()
	}
	respondJSON(c, http.StatusOK, customerList)
}

// GetSubscription is a ...
func GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	modelCustomer := models.Customer{}
	_customer, err := modelCustomer.FindCustomerByID(customerId)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(c, http.StatusOK, _customer)
}

func getCustomerFromForm(c *gin.Context) (*models.Customer, bool) {
	n := c.PostForm("name")
	if n == "" {
		respondJSON(c, http.StatusBadRequest, "Name is invalid")
		return nil, true
	}

	e := c.PostForm("email")
	if e == "" {
		respondJSON(c, http.StatusBadRequest, "Email is invalid")
		return nil, true
	}

	r := c.PostForm("role")
	if r == "" {
		respondJSON(c, http.StatusBadRequest, "Role is invalid")
		return nil, true
	}

	status := false
	if c.PostForm("status") != "" {
		status = true
	}

	modelCustomer := &models.Customer{
		Name:   n,
		Email:  e,
		Role:   models.Role(r),
		Status: status,
	}
	return modelCustomer, false
}

func CreateCustomer(c *gin.Context) {
	modelCustomer, shouldReturn := getCustomerFromForm(c)
	if shouldReturn {
		return
	}

	_customer, err := modelCustomer.SaveCustomer()
	if err != nil {
		respondJSON(c, http.StatusBadRequest, "Cannot save customer, duplicate name or email")
		return
	}

	respondJSON(c, http.StatusOK, _customer)
}

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	customerId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	modelCustomer, shouldReturn := getCustomerFromForm(c)
	if shouldReturn {
		return
	}

	_customer, err := modelCustomer.UpdateCustomer(customerId)

	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, _customer)
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
