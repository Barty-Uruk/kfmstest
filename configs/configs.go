package configs

import (
	"encoding/json"
	"flag"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const appName = "KFMS"

var options = []option{
	// config section
	{"config", "string", "", "config file"},

	// server config section
	{"servers.main.http.port", "int", 8080, "Server http port"},
	{"servers.public.http.port", "int", 8081, "Server http port"},
	{"server.protection.apikey", "string", "1cf1c47908fa1cf19a38e5d963d84958", "Key for protected api"},

	// opentracing config section
	{"opentracing.enabled", "bool", false, "opentracing enabled"},
	{"opentracing.servicename", "string", appName, "opentracing service name"},
	{"opentracing.reporter.logspans", "bool", false, "opentracing whether the reporter should also log the spans"},
	{"opentracing.agent.host", "string", "jaeger-agent", "opentracing host jaeger agent"},
	{"opentracing.agent.port", "int", 6831, "opentracing port jaeger agent"},
	{"opentracing.sampler.type", "string", "const", "opentracing type tracing"},
	{"opentracing.sampler.param", "float64", 1, "opentracing param sampling"},

	{"s3.address", "string", "", "addres of s3 storage"},
	{"s3.bucket", "string", "storage", "s3 bucket name"},
	{"s3.rootfoldername", "string", "", "root folder name with /"},

	// logger config section
	{"logger.level", "string", "info", "LogLevel is global log level:  EMERG(0), ALERT(1), CRIT(2), ERR(3), WARNING(4), NOTICE(5), INFO(6), DEBUG(7)"},
	{"logger.timeformat", "string", "2006-01-02T15:04:05.999999999Z07:00", "LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00"},
}

type Config struct {
	Servers     Servers `yaml:"Servers"`
	Opentracing Opentracing
	S3          AmazonS3
	Logger      Logger
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
type AmazonS3 struct {
	Address        string
	Bucket         string
	RootFolderName string
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
