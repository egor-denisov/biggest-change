package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/go-playground/assert"
)

const testConfigStr = `
app:
  name: "test-app"
  version: "1.0.0"
  countOfBlocks: 100
  maxGoroutines: 50
  averageAddressesInBlock: 200
  cacheSize: 100

api:
  url: test-URL
  rps: 60
  timewindow: 1s
  timeout: 5s
  maxRetries: 5
  timeBetweenRetries: 500ms

http:
  port: ":8080"
  timeout: 5s

logger:
  logLevel: "info"
`

const testEnvStr = `
APP_NAME=test-app
APP_VERSION=1.0.0
APP_COUNT_OF_BLOCKS=100
APP_MAX_GOROUTINES=50
APP_AVG_ADDRS=200
APP_CACHE_SIZE=100
API_URL=test-URL
API_RPS=60
API_TIME_WINDOW_RPS=1s
API_TIMEOUT=5s
API_TIME_BETWEEN_RETRIES=500ms
HTTP_PORT=:8080
HTTP_TIMEOUT=5s
LOG_LEVEL=info`

func Test_MustLoadPath_ExistentPath(t *testing.T) {
	for _, test := range testsMustLoadPath {
		t.Run(test.name, func(t *testing.T) {
			// temp config.yml file
			tempFileConfig, err := os.CreateTemp("", "config-*.yml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tempFileConfig.Name())

			_, err = tempFileConfig.WriteString(test.configFile)
			if err != nil {
				t.Fatal(err)
			}
			// temp .env file
			tempFileEnv, err := os.CreateTemp("", "*.env")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tempFileEnv.Name())

			_, err = tempFileEnv.WriteString(test.envFile)
			if err != nil {
				t.Fatal(err)
			}

			config := MustLoadPath(tempFileConfig.Name(), tempFileEnv.Name())
			assert.Equal(t, test.expectedConfig, config)
		})
	}
}

var testsMustLoadPath = []struct {
	name           string
	configFile     string
	envFile        string
	expectedConfig *Config
}{
	{
		name:       "Getting default value (w/o config.yml & .env)",
		configFile: "something: else",
		envFile:    "",
		expectedConfig: &Config{
			App: App{
				Name:                    "biggest-change",
				Version:                 "1.0.0",
				CountOfBlocks:           100,
				MaxGoroutines:           50,
				AverageAddressesInBlock: 200,
				CacheSize:               100,
			},
			API: API{
				URL:                "",
				Rps:                60,
				TimeWindowRPS:      time.Second,
				Timeout:            5 * time.Second,
				MaxRetries:         5,
				TimeBetweenRetries: 500 * time.Millisecond,
			},
			HTTP: HTTP{
				Port:    ":8080",
				Timeout: 5 * time.Second,
			},
			Log: Log{
				Level: "debug",
			},
		},
	},
	{
		name:       "Only config.yml",
		configFile: testConfigStr,
		envFile:    "",
		expectedConfig: &Config{
			App: App{
				Name:                    "test-app",
				Version:                 "1.0.0",
				CountOfBlocks:           100,
				MaxGoroutines:           50,
				AverageAddressesInBlock: 200,
				CacheSize:               100,
			},
			API: API{
				URL:                "test-URL",
				Rps:                60,
				TimeWindowRPS:      time.Second,
				Timeout:            5 * time.Second,
				MaxRetries:         5,
				TimeBetweenRetries: 500 * time.Millisecond,
			},
			HTTP: HTTP{
				Port:    ":8080",
				Timeout: 5 * time.Second,
			},
			Log: Log{
				Level: "info",
			},
		},
	},
	{
		name:       "Change something with .env",
		configFile: testConfigStr,
		envFile:    "APP_NAME:else-name",
		expectedConfig: &Config{
			App: App{
				Name:                    "else-name",
				Version:                 "1.0.0",
				CountOfBlocks:           100,
				MaxGoroutines:           50,
				AverageAddressesInBlock: 200,
				CacheSize:               100,
			},
			API: API{
				URL:                "test-URL",
				Rps:                60,
				TimeWindowRPS:      time.Second,
				Timeout:            5 * time.Second,
				MaxRetries:         5,
				TimeBetweenRetries: 500 * time.Millisecond,
			},
			HTTP: HTTP{
				Port:    ":8080",
				Timeout: 5 * time.Second,
			},
			Log: Log{
				Level: "info",
			},
		},
	},
	{
		name:       "Only .env",
		configFile: "something: else",
		envFile:    testEnvStr,
		expectedConfig: &Config{
			App: App{
				Name:                    "else-name",
				Version:                 "1.0.0",
				CountOfBlocks:           100,
				MaxGoroutines:           50,
				AverageAddressesInBlock: 200,
				CacheSize:               100,
			},
			API: API{
				URL:                "test-URL",
				Rps:                60,
				TimeWindowRPS:      time.Second,
				Timeout:            5 * time.Second,
				MaxRetries:         5,
				TimeBetweenRetries: 500 * time.Millisecond,
			},
			HTTP: HTTP{
				Port:    ":8080",
				Timeout: 5 * time.Second,
			},
			Log: Log{
				Level: "info",
			},
		},
	},
}

func Test_MustLoadPath_NonExistentPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	MustLoadPath("non_existent_config.yml", "non_existent_env.env")
}

func Test_fetchConfigPath(t *testing.T) {
	for _, test := range testsFetchConfigPath {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.argsValue
			t.Setenv("CONFIG_PATH", test.envValue)

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			configPath := fetchConfigPath()

			assert.Equal(t, test.expected, configPath)
		})
	}
}

var testsFetchConfigPath = []struct {
	name      string
	argsValue []string
	envValue  string
	expected  string
}{
	{
		name:      "Not field",
		argsValue: []string{"cmd", ""},
		envValue:  "",
		expected:  _defaultConfigPath,
	},
	{
		name:      "Ok - from environment",
		argsValue: []string{"cmd", ""},
		envValue:  "./test_config2.yml",
		expected:  "./test_config2.yml",
	},
	{
		name:      "Ok - from flag",
		argsValue: []string{"cmd", "-config", "./test_config.yml"},
		envValue:  "",
		expected:  "./test_config.yml",
	},
}
