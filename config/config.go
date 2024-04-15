package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	_defaultConfigPath = "./config/config.yml"
	_defaultEnvPath    = ".env"
)

type (
	Config struct {
		App  `yaml:"app"`
		API  `yaml:"api"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
	}

	App struct {
		Name                    string `env:"APP_NAME"            env-default:"biggest-change" yaml:"name"`
		Version                 string `env:"APP_VERSION"         env-default:"1.0.0"          yaml:"version"`
		DefaultCountOfBlocks    uint   `env:"APP_COUNT_OF_BLOCKS" env-default:"100"            yaml:"countOfBlocks"`
		MaxGoroutines           int    `env:"APP_MAX_GOROUTINES"  env-default:"50"             yaml:"maxGoroutines"`
		AverageAddressesInBlock int    `env:"APP_AVG_ADDRS"       env-default:"200"            yaml:"averageAddressesInBlock"`
		CacheSize               int    `env:"APP_CACHE_SIZE"      env-default:"100"            yaml:"cacheSize"`
	}

	API struct {
		Url                string        `env:"API_URL"                  env-default:""      yaml:"url"`
		Rps                int           `env:"API_RPS"                  env-default:"60"    yaml:"rps"`
		TimeWindowRPS      time.Duration `env:"API_TIME_WINDOW_RPS"      env-default:"1s"    yaml:"timewindow"`
		Timeout            time.Duration `env:"API_TIMEOUT"              env-default:"5s"    yaml:"timeout"`
		MaxRetries         int           `env:"API_MAX_RETRIES"          env-default:"5"     yaml:"maxRetries"`
		TimeBetweenRetries time.Duration `env:"API_TIME_BETWEEN_RETRIES" env-default:"500ms" yaml:"timeBetweenRetries"`
	}

	HTTP struct {
		Port    string        `env:"HTTP_PORT"    env-default:":8080" yaml:"port"`
		Timeout time.Duration `env:"HTTP_TIMEOUT" env-default:"5s"    yaml:"timeout"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"debug" yaml:"logLevel"`
	}
)

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath, _defaultEnvPath)
}

func MustLoadPath(configPath, envPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	// try loading .env file
	godotenv.Load(envPath)

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		b, _ := os.ReadFile(configPath)
		fmt.Println(string(b))
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	fmt.Println(res)

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		res = _defaultConfigPath
	}

	return res
}
