package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/go-playground/assert"
)

var testConfigStr = `
app:
  name: "test-app"
  version: "1.0.0"

api:
  url: test-url
  rps: 60

http:
  port: ":8080"
  timeout: "5s"

logger:
  logLevel: "info"
`

var testEnvStr = `
APP_NAME=test-app
APP_VERSION=1.0.0
API_URL=test-url
API_RPS=60
HTTP_PORT=:8080
HTTP_TIMEOUT=5s
LOG_LEVEL=info`

func Test_MustLoadPath_ExistentPath(t *testing.T) {
	for _, test := range tests_MustLoadPath {
		t.Run(test.name, func(t *testing.T) {
			// temp config.yml file
			tempFileConfig, err := os.CreateTemp("", "config-*.yml")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tempFileConfig.Name())

			_, err = tempFileConfig.Write([]byte(test.configFile))
			if err != nil {
				t.Fatal(err)
			}
			// temp .env file
			tempFileEnv, err := os.CreateTemp("", "*.env")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tempFileEnv.Name())

			_, err = tempFileEnv.Write([]byte(test.envFile))
			if err != nil {
				t.Fatal(err)
			}

			config := MustLoadPath(tempFileConfig.Name(), tempFileEnv.Name())
			assert.Equal(t, test.expectedConfig, config)
		})
	}
}

var tests_MustLoadPath = []struct {
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
				Name:    "biggest-change",
				Version: "1.0.0",
			},
			API: API{
				Url: "",
				Rps: 60,
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
				Name:    "test-app",
				Version: "1.0.0",
			},
			API: API{
				Url: "test-url",
				Rps: 60,
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
				Name:    "else-name",
				Version: "1.0.0",
			},
			API: API{
				Url: "test-url",
				Rps: 60,
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
				Name:    "else-name",
				Version: "1.0.0",
			},
			API: API{
				Url: "test-url",
				Rps: 60,
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

var tests_fetchConfigPath = []struct {
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

func Test_fetchConfigPath(t *testing.T) {
	for _, test := range tests_fetchConfigPath {
		t.Run(test.name, func(t *testing.T) {
			os.Args = test.argsValue
			os.Setenv("CONFIG_PATH", test.envValue)

			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			configPath := fetchConfigPath()

			assert.Equal(t, test.expected, configPath)
		})
	}
}
