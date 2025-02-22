package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Auth struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type NumPolicy struct {
	Num int `yaml:"number"`
}

// RegexPolicy removes all images that match the given regex.
type RegexPolicy struct {
	// Repos contains list of regex to match repo name
	Repos []string `yaml:"repos"`
	/// Tags contains list of regex to match tag name
	Tags []string `yaml:"tags"`
}

// NotTouchedPolicy cleans images that are recently not touched within given period
type NotTouchedPolicy struct {
	// Time is time period in second.
	Time int64 `yaml:"time"`
}

type Policy struct {
	// Type of the policy, e.g. "number", "regex", "recentlyNotTouched"
	Type string `yaml:"type"`
	// NumPolicy configures policy to retain given number tags in repo
	NumPolicy *NumPolicy `yaml:"numberPolicy,omitempty"`
	// RegexPolicy configures policy to clean images that match the regex patterns
	RegexPolicy *RegexPolicy `yaml:"regexPolicy,omitempty"`
	// TouchPolicy configures policy to clean images that are recently not touched within given period
	NotTouchedPolicy *NotTouchedPolicy `yaml:"notTouchedPolicy,omitempty"`
	// RetainTags is tag patterns to be retained
	RetainTags []string `yaml:"retainTags"`
}

type Trigger struct {
	// Cron expression to regularly trigger the cleanup
	Cron string `yaml:"cron"`
}

type XSRF struct {
	Enabled bool   `yaml:"enabled"`
	Key     string `yaml:"key"`
}

type C struct {
	Host     string   `yaml:"host"`
	Insecure bool     `yaml:"insecure"`
	Version  string   `yaml:"version"`
	Auth     Auth     `yaml:"auth"`
	Projects []string `yaml:"projects"`
	Policy   Policy   `yaml:"policy"`
	Trigger  *Trigger `yaml:"trigger"`
	XSRF     XSRF     `yaml:"xsrf"`
}

var Config = C{}

func Load(configFile string) error {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		logrus.WithField("f", configFile).Error("Read config file error: ", err)
		return err
	}

	err = yaml.Unmarshal(b, &Config)
	if err != nil {
		logrus.WithField("f", configFile).Error("Unmarshal config file error: ", err)
		return err
	}

	return Normalize(&Config)
}

func Normalize(c *C) error {
	trimed := strings.TrimSpace(c.Version)
	trimed = strings.TrimPrefix(trimed, "v")
	if len(trimed) < 3 {
		return fmt.Errorf("unrecoganized version %s, please provide version like 1.4, 1.7.5", c.Version)
	}
	c.Version = trimed[:3]

	if HasCronSchedule() {
		_, err := cron.ParseStandard(c.Trigger.Cron)
		if err != nil {
			return err
		}
	}

	return nil
}

func HasCronSchedule() bool {
	return Config.Trigger != nil && len(Config.Trigger.Cron) > 0
}
