package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

	// Define routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

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
