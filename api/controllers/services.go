package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sano/database"
)

// List all services, just says if last one check was ok or not

func ServicesIndex(c *gin.Context) {
	lookups := database.GetLookups()
	c.JSON(http.StatusAccepted, lookups)
}

func ServicesShow(c *gin.Context) {
	var serviceName = c.Param("service")
	serviceLookups := database.GetLookup(serviceName)
	c.JSON(http.StatusAccepted, serviceLookups)
}
