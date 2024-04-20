package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/andreastihor/isd/isdsvc/backend/services/isd"
	"github.com/andreastihor/isd/isdsvc/backend/storage/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config struct to hold configuration data
type Config struct {
	Port int
	URL  string
}

var (
	serviceName = "ISD-backend"
	version     = "devel"
)

func main() {
	// Read config
	config := readConfig()
	logger, err := initLogger(config)
	if err != nil {
		log.Fatalf("error initializing logger: %v", err)
	}

	//initiate DB String
	dbConn := GetDBConnString(config)
	clubStorage, err := postgres.NewStorage(logger, dbConn)
	if err != nil {
		logger.Fatalf("Failed to initialize clubStorage storage: %v", err)
	}

	// Initialize your handler with the storage
	handler := isd.NewHandler(clubStorage)

	// Register routes
	isd.RegisterRoutes(handler)

	// Start server

	// Use "server.port" and "server.url" to access values under the [server] section
	addr := fmt.Sprintf("%s:%d", config.GetString("server.url"), config.GetInt("server.port"))
	logger.Infof("Server listening on %s\n", addr)
	logger.Fatal(http.ListenAndServe(addr, nil))
}

// readConfig reads configuration using viper
func readConfig() *viper.Viper {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)

	config.SetConfigFile("env/config")
	config.SetConfigType("ini") // Specify the config type as "ini"
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	return config
}

func initLogger(config *viper.Viper) (*logrus.Entry, error) {
	l := logrus.New()
	var logLevel logrus.Level

	llStr := config.GetString("server.logLevel")
	if llStr == "fromenv" {
		switch config.GetString("runtime.environment") {
		case "staging", "development":
			logLevel = logrus.DebugLevel // to simplify debugging
		default: // including production
			logLevel = logrus.InfoLevel
		}
	} else {
		var err error
		logLevel, err = logrus.ParseLevel(llStr)
		if err != nil {
			return nil, err
		}
	}

	l.SetLevel(logLevel)
	return l.WithFields(logrus.Fields{
		"service": serviceName,
		"version": version,
	}), nil
}

func GetDBConnString(config *viper.Viper) string {
	host := config.GetString("database.host")
	port := config.GetString("database.port")
	username := config.GetString("database.user")
	password := config.GetString("database.password")
	databaseName := config.GetString("database.dbname")
	sslMode := config.GetString("database.sslmode")
	// postgres: //user:password@localhost:5432/mydb?sslmode=disable
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", username, password, host, port, databaseName, sslMode)

}
