package routerx

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"heart/internal/handlers"
	"heart/internal/middleware"
	"net/http"
	"strings"
)

func Router(origins string) *gin.Engine {
	r := gin.Default()

	r.Use(CORSMiddleware(origins))

	r.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "up"})
		},
	)

	r.GET(
		"/version",
		func(c *gin.Context) {
			c.JSON(http.StatusOK, Info())
		},
	)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	//r.POST("/refresh", Route(handlers.Refresh))

	// Protected routes
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthenticationMiddleware())
	//authGroup.GET("/me", Authenticated(handlers.Me))

	exercisesGroup := r.Group("/exercises")
	exercisesGroup.Use(middleware.AuthenticationMiddleware())
	exercisesGroup.GET("", Authenticated(handlers.GetExercises))

	return r
}

// CORSMiddleware returns a Gin middleware that handles CORS requests.
func CORSMiddleware(origins string) gin.HandlerFunc {
	parsed := allowedOrigins(origins)
	originSet := make(map[string]struct{}, len(parsed))
	for _, o := range parsed {
		originSet[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if _, ok := originSet[origin]; ok {
			corsHeaders(c, origin)
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func allowedOrigins(origins string) []string {
	var result []string
	for _, o := range strings.Split(origins, ",") {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func corsHeaders(c *gin.Context, origin string) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}
