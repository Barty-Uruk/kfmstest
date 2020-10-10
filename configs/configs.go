package configs

import (
	"encoding/json"
	"flag"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const appName = "kas"

var options = []option{
	// config section
	{"config", "string", "", "config file"},

	// server config section
	{"servers.main.http.port", "int", 8080, "Server http port"},
	{"servers.public.http.port", "int", 8081, "Server http port"},
	{"server.protection.apikey", "string", "1cf1c47908fa1cf19a38e5d963d84958", "Key for protected api"},

	// jwt section
	{"jwt.secret", "string", "s", "The secret key"},
	{"jwt.aeskey", "string", "LKHlhb899Y09olUi", "The aes key"},
	{"jwt.ttlsec", "int", 1800, "Jwt ttl"},

	// aerospike config section
	{"aerospike.enabled", "bool", false, "aerospike enabled"},
	{"aerospike.host", "string", appName, "aerospike host"},
	{"aerospike.port", "int", 3000, "The Aerospike port"},
	{"aerospike.namespace", "string", appName + "-jwt-tokens", "The Aerospike namespace"},
	{"aerospike.connectTimeout", "int", 30, "Initial host connection timeout duration. The timeout when opening a connection to the server host for the first time (in sec)"},
	{"aerospike.connectIdleTimeout", "int", 14, "Connection idle timeout. Every time a connection is used, its idle deadline will be extended by this duration. When this deadline is reached, the connection will be closed and discarded from the connection pool (in sec)"},
	{"aerospike.connectAttempts", "int", 20, "Number of attempts to connect to the database"},
	{"aerospike.connectAttemptsPause", "int", 10, "Pause in seconds between attempts to connect to the database (in sec)"},

	// database config section
	{"database.driver", "string", "postgresql", "database driver"},
	{"database.host", "string", "localhost", "database host"},
	{"database.port", "int", 5432, "database port"},
	{"database.user", "string", "root", "database user"},
	{"database.password", "string", "empty", "database password"},
	{"database.databasename", "string", appName, "database name"},
	{"database.secure", "string", "disable", "database SSL support"},

	// nats config section
	{"nats.enabled", "bool", false, "deviceinfo enabled"},
	{"nats.host", "string", "localhost", "host addr of nats"},
	{"nats.port", "int", 4222, "port addr of nats"},
	{"nats.username", "string", "nats_client", "username from authentification"},
	{"nats.password", "string", "N7KySil1ES", "password from authentification"},

	// opentracing config section
	{"opentracing.enabled", "bool", false, "opentracing enabled"},
	{"opentracing.servicename", "string", appName, "opentracing service name"},
	{"opentracing.reporter.logspans", "bool", false, "opentracing whether the reporter should also log the spans"},
	{"opentracing.agent.host", "string", "jaeger-agent", "opentracing host jaeger agent"},
	{"opentracing.agent.port", "int", 6831, "opentracing port jaeger agent"},
	{"opentracing.sampler.type", "string", "const", "opentracing type tracing"},
	{"opentracing.sampler.param", "float64", 1, "opentracing param sampling"},

	// logger config section
	{"logger.level", "string", "info", "LogLevel is global log level:  EMERG(0), ALERT(1), CRIT(2), ERR(3), WARNING(4), NOTICE(5), INFO(6), DEBUG(7)"},
	{"logger.timeformat", "string", "2006-01-02T15:04:05.999999999Z07:00", "LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00"},

	// app config section
	{"app.cert", "string", "", "Crypto certificate for Notebooks"},
	{"app.applicationID", "string", "default", "Default applicationID"},
	{"app.auth.sendPush", "bool", false, "Send checkcode over push service"},
	{"app.salt.password", "string", "lSLW3QF6d6Jcwmtptdeg8zxFk", "Password salt"},
	{"app.salt.phone", "string", "lID7VO3q4WTbU9KLv7lH", "Phone salt"},
}

type Config struct {
	Servers     Servers `yaml:"Servers"`
	Database    Database
	Nats        Nats
	Opentracing Opentracing
	Logger      Logger
	Aerospike   Aerospike
	Jwt         Jwt
	App         App
}

type Jwt struct {
	Secret string
	AesKey string
	TtlSec int
}

type Aerospike struct {
	Enabled              bool
	Host                 string
	Port                 int
	Namespace            string
	ConnectTimeout       int
	ConnectIdleTimeout   int
	ConnectAttempts      int
	ConnectAttemptsPause int
}

type Protection struct {
	Apikey string
}

type Servers struct {
	Main   Server
	Public Server
}

type Server struct {
	Http Http
}

type Http struct {
	Port int
}

type App struct {
	Cert          string
	Salt          Salt
	ApplicationID string
	Auth          struct {
		SendPush bool `mapstructure:"sendPush"`
	}
}

type Salt struct {
	Password string
	Phone    string
}

type Nats struct {
	Enabled  bool
	Host     string
	Port     int
	Username string
	Password string
}

type Database struct {
	Driver       string
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	Secure       string
}

type Opentracing struct {
	Enabled     bool
	ServiceName string
	Agent       OpentracingAgent
	Reporter    OpentracingReporter
	Sampler     OpentracingSampler
}

type OpentracingSampler struct {
	Type  string
	Param float64
}

type OpentracingAgent struct {
	Host string
	Port int
}

type OpentracingReporter struct {
	LogSpans bool
}

type Logger struct {
	Level      string
	TimeFormat string
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}

func NewConfig() *Config {
	return &Config{}
}

// Read read parameters for config.
// Read from environment variables, flags or file.
func (c *Config) Read() error {
	viper.SetEnvPrefix(appName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, o := range options {
		switch o.typing {
		case "string":
			pflag.String(o.name, o.value.(string), o.description)
		case "int":
			pflag.Int(o.name, o.value.(int), o.description)
		case "bool":
			pflag.Bool(o.name, o.value.(bool), o.description)
		default:
			viper.SetDefault(o.name, o.value)
		}
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	if fileName := viper.GetString("config"); fileName != "" {
		viper.SetConfigFile(fileName)
		viper.SetConfigType("toml")

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}

// Print print config structure
func (c *Config) Print() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	log.Println(string(b))
	return nil
}
