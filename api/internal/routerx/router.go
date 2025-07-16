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

	exercisesGroup := r.Group("/exercises")
	exercisesGroup.Use(middleware.AuthenticationMiddleware())
	exercisesGroup.GET("", Authenticated(handlers.GetExercises))

	workoutsGroup := r.Group("/workouts")
	workoutsGroup.Use(middleware.AuthenticationMiddleware())
	workoutsGroup.GET("", Authenticated(handlers.GetWorkouts))
	workoutsGroup.POST("", Authenticated(handlers.MakeWorkout))
	workoutsGroup.GET(":workoutId", Authenticated(handlers.GetWorkout))
	workoutsGroup.DELETE(":workoutId", Authenticated(handlers.DeleteWorkout))

	templatesGroup := r.Group("/templates")
	templatesGroup.Use(middleware.AuthenticationMiddleware())
	templatesGroup.GET("", Authenticated(handlers.GetTemplates))
	templatesGroup.POST("", Authenticated(handlers.MakeTemplate))
	templatesGroup.GET(":templateId", Authenticated(handlers.GetTemplate))
	templatesGroup.DELETE(":templateId", Authenticated(handlers.DeleteTemplate))

	accountGroup := r.Group("/accounts")
	accountGroup.Use(middleware.AuthenticationMiddleware())
	accountGroup.POST("", Authenticated(handlers.RegisterAccount))
	accountGroup.DELETE("", Authenticated(handlers.DeleteAccount))
	accountGroup.PUT(":accountId", Authenticated(handlers.EditAccount))
	accountGroup.GET(":accountId", Authenticated(handlers.GetAccount))

	feedbackGroup := r.Group("/feedback")
	feedbackGroup.Use(middleware.AuthenticationMiddleware())
	feedbackGroup.POST("", Authenticated(handlers.LeaveFeedback))

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
	c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}

const allowedHeaders = `Content-Type,Authorization,Accept,Accept-Language,X-Timezone,X-App-Version,Referer,User-Agent,`
