package config

import (
	"runtime"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

// Build information. Populated at build-time.
var (
	Version   string
	Branch    string
	Commit    string
	GoVersion = runtime.Version()
)

//ZeroUUID is a zero-value of uuid.UUID
var ZeroUUID = uuid.UUID{}

const (

	//DateLayout is the formatting string for all date-time objects
	DateLayout = time.RFC3339Nano

	//EnvPrefix is a prefix to all ENV variables used in this app
	EnvPrefix = "BID"
	//APIPrefixV1 URL prefix in API version 1
	APIPrefixV1 = "/api/v1"
	//DefaultPort default port the service is served on
	DefaultPort = "9000"
)

// ErrorMessage defines the type for the errors channel
type ErrorMessage struct {
	Message string
	Err     error
}

func bindEnvVariable(name string, fallback interface{}) {
	viper.SetDefault(name, fallback)
	viper.BindEnv(name)
}

//SetupEnv configures app to read ENV variables
func SetupEnv() {
	viper.SetEnvPrefix(EnvPrefix)
	// General
	bindEnvVariable("PORT", DefaultPort)
}
