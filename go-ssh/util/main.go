package util

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

type ScriptConfig struct {
	Name          string             `yaml:"name"`
	Script        string             `yaml:"script"`
	Timeout       string             `yaml:"timeout"`
	Pattern       string             `yaml:"pattern"`
	Credentials   []CredentialConfig `yaml:"credentials"`
	ParsedTimeout time.Duration      // For internal use only
	Ignored       bool               // For internal use only
}

type CredentialConfig struct {
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	User               string `yaml:"user"`
	KeyFile            string `yaml:"keyfile"`
	ScriptResult       string // For internal use only
	ScriptReturnCode   int    // For internal use only
	ScriptError        string // For internal use only
	ResultPatternMatch int8   // For internal use only
}

type Config struct {
	Version string         `yaml:"version"`
	Scripts []ScriptConfig `yaml:"scripts"`
}

func ParseFlags(c, p *string) (*string, *string) {

	flag.StringVar(c, "config", "config.yml", "Path to your ssh_exporter config file")
	flag.StringVar(p, "port", "9428", "Port probed metrics are served on.")

	flag.Parse()

	return c, p
}

func ParseConfig(c string) (Config, error) {

	raw, err := ioutil.ReadFile(c)
	FatalCheck(err)

	initialConfig := Config{}
	err = yaml.Unmarshal([]byte(raw), &initialConfig)
	SoftCheck(err)

	finalConfig, err := adjustConfig(initialConfig)
	SoftCheck(err)

	return finalConfig, err
}

func FatalCheck(e error) {

	if e != nil {
		log.Fatal("error: ", e)
	}
}

func SoftCheck(e error) bool {

	if e != nil {
		LogMsg(fmt.Sprintf("%v", e))
		return true
	} else {
		return false
	}
}

func LogMsg(s string) {

	log.Printf("ssh_exporter :: %s", fmt.Sprintf("%s", s))
}

func adjustConfig(c Config) (Config, error) {

	for c_i, v_i := range c.Scripts {
		for c_j, v_j := range v_i.Credentials {
			if v_j.Port == "" {
				c.Scripts[c_i].Credentials[c_j].Port = "22"
			}
		}

		tmp, err := time.ParseDuration(c.Scripts[c_i].Timeout)
		if !SoftCheck(err) {
			c.Scripts[c_i].ParsedTimeout = tmp
		} else {
			LogMsg(fmt.Sprintf("Failed to parse `timeout` for %s. Default to 10s", c.Scripts[c_i].Name))
			c.Scripts[c_i].ParsedTimeout, _ = time.ParseDuration("10s")
		}
	}

	return c, nil
}
