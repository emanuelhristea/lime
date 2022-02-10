package controllers

import (
	"github.com/gin-gonic/gin"
)

type requestLicense struct {
	StripeID string `json:"customer"`
	License  string `json:"license"`
}

// ResponseData is a ...
type ResponseData struct {
	Status int         `json:"status"`
	Msg    interface{} `json:"msg"`
}

func respondJSON(g *gin.Context, status int, msg interface{}) {
	res := &ResponseData{
		Status: status,
		Msg:    msg,
	}
	g.JSON(status, res)
}
