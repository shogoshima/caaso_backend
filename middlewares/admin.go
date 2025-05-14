package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var allowedIpAdresses = []string{
	"3.77.79.249",
	"3.77.79.250",
	"35.90.103.132",
	"35.90.103.133",
	"35.90.103.134",
	"35.90.103.135",
	"44.208.168.68",
	"44.208.168.69",
	"44.208.168.70",
	"44.208.168.71",
}

func OriginWhitelist(c *gin.Context) {

	ip := c.ClientIP()
	println(ip)

	allowed := false
	for _, allowedIp := range allowedIpAdresses {
		if ip == allowedIp {
			allowed = true
			break
		}
	}

	if !allowed {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden: Unauthorized Origin"})
		return
	}

	c.Next()
}

// Middleware to validate admin token
func CheckAdmin(c *gin.Context) {

	token := c.GetHeader("Authorization")

	if token != "Bearer "+os.Getenv("RETOOL_ADMIN_TOKEN") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: Invalid Token"})
		return
	}

	c.Next()
}
