package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Config struct {
	SmtpServer         string
	SmtpPort           int
	SmtpUsername       string
	SmtpPassword       string
	Subject            string
	MongoHost          string
	Schedule           int
	Recipients         []string
	InsecureSkipVerify bool
	UnencryptedAuth    bool
	SkipAuth           bool
	TemplateFile       string
	EmailFormat        string
	Notificationsdb    string
}

func NewConfig() (config *Config, err error) {
	var file []byte
	file, err = ioutil.ReadFile("config.json")

	if err != nil {
		return nil, err
	}

	config = new(Config)
	if err = json.Unmarshal(file, config); err != nil {
		return nil, err
	}

	if config.EmailFormat == "" {
		config.EmailFormat = "text/plain"
	}

	if config.Subject == "" {
		config.Subject = "avorsmtp notification"
	}

	if config.Notificationsdb == "" {
		config.Notificationsdb = "notificationsdb"
	}

	return config, nil
}

func schedule(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
