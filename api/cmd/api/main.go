//	@title			Heart of Yours API
//	@version		1.0
//	@description	A simple fitness tracker API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Kit
//	@contact.url	https://github.com/kit-g

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

// @host		localhost:8080
// @BasePath	/
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	_ "heart/docs"
	"heart/internal/awsx"
	"heart/internal/config"
	"heart/internal/dbx"
	"heart/internal/firebasex"
	"heart/internal/models"
	"heart/internal/routerx"

	"log"
	"os"
)

var cfg *config.AppConfig

func Init() {
	var err error
	cfg, err = config.NewAppConfig()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	if cfg.FirebaseConfig.Credentials != "" {
		if err := firebasex.Init(cfg.FirebaseConfig.Credentials); err != nil {
			log.Fatal("Failed to initialize Firebase client:", err)
		}
	}

	if err := dbx.Connect(&cfg.DBConfig); err != nil {
		log.Fatal("Failed to connect to DB:", err)
		return
	}

	if err := awsx.InitS3(cfg.UploadBucket, cfg.AwsRegion); err != nil {
		log.Fatal("Failed to initialize S3 client:", err)
		return
	}

	_ = dbx.DB.AutoMigrate(
		&models.User{},
		&models.Exercise{},
		&models.Workout{},
		&models.WorkoutExercise{},
		&models.Set{},
		&models.Template{},
	)
}

func main() {
	Init()

	r := routerx.Router(cfg.CORSOrigins)

	mode := os.Getenv("MODE")

	if mode == "lambda" {
		// Route Gin router with the Lambda adapter
		log.Println("Running in Lambda mode...")
		lambda.Start(ginadapter.New(r).ProxyWithContext)
	} else {
		// Default: local dev mode
		log.Println("Running in local mode on :8080...")
		if err := r.Run(":8080"); err != nil {
			log.Fatal("Server failed to start:", err)
		}
	}
}
