package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"sano/api/controllers"
	"strconv"
)

func Start(configHttpPort int) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	servicesApiGroup := r.Group("/services")
	servicesApiGroup.GET("/", controllers.ServicesIndex)
	servicesApiGroup.GET("/:service", controllers.ServicesShow)

	httpPort := captureAvailablePortForHttp(configHttpPort)
	if httpPort != configHttpPort {
		log.Printf("[%d] is used by another service, first available port after this is [%d], so I choose this for API", configHttpPort, httpPort)
	}

	log.Println("APIs are ready!")
	err := r.Run(":" + strconv.Itoa(httpPort))
	if err != nil {
		log.Panicln(err)
	}
}

func captureAvailablePortForHttp(startFrom int) int {
	availablePort := startFrom

	for {
		availablePort = startFrom + 1
		_, err := net.Dial("tcp", ":"+strconv.Itoa(availablePort))
		if err != nil {
			break
		}
	}

	return availablePort
}
