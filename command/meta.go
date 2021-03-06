package main

import (
	"flag"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Meta struct {
	Authtoken string
	Subdomain string
	Loglevel  string
}

type FlagSetFlags uint

func (m *Meta) FlagSet(n string) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)
	f.StringVar(&m.Authtoken, "authtoken", "", "PagerDuty API authentication token")
	f.StringVar(&m.Subdomain, "subdomain", "", "PagerDuty account name (subdomain)")
	f.StringVar(&m.Loglevel, "loglevel", "", "Logging level")
	return f
}

func (m *Meta) Client() *pagerduty.Client {
	return pagerduty.NewClient(m.Subdomain, m.Authtoken)
}

func (m *Meta) Help() string {
	helpText := `
	Generral options:

	-authtoken PagerDuty API authentication token
	-subdomain PagerDuty account ID
	-loglevel Logging level
`
	return strings.TrimSpace(helpText)
}

func (m *Meta) validate() error {
	if m.Authtoken == "" {
		return fmt.Errorf("Authtoken can not be blank")
	}
	if m.Subdomain == "" {
		return fmt.Errorf("Subdomain can not be blank")
	}
	return nil
}

func (m *Meta) Setup() error {
	m.setupLogging()
	if err := m.loadConfig(); err != nil {
		log.Warn(err)
	}
	return m.validate()
}

func (m *Meta) setupLogging() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	switch m.Loglevel {
	case "info", "":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Fatal("Unknown log level", m.Loglevel)
	}
}

func (m *Meta) loadConfig() error {
	path, err := homedir.Dir()
	if err != nil {
		return err
	}
	configFile := filepath.Join(path, ".pd.yml")
	if _, err := os.Stat(configFile); err != nil {
		return err
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	other := &Meta{}
	if err := yaml.Unmarshal(data, other); err != nil {
		return err
	}
	if m.Authtoken == "" {
		m.Authtoken = other.Authtoken
	}
	if m.Loglevel == "" {
		m.Loglevel = other.Loglevel
	}
	if m.Subdomain == "" {
		m.Subdomain = other.Subdomain
	}
	return nil
}
