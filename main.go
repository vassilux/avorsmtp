package main

import (
	//"labix.org/v2/mgo/bson"
	log "github.com/cihub/seelog"
	"os"
	"time"
)

var (
	smtpWorker       *SmtpWorker
	mongoWorker      *MongoWorker
	config           *Config
	stopPocessEvents chan bool
)

func loadLogger() {
	logger, err := log.LoggerFromConfigAsFile("logger.xml")

	if err != nil {
		log.Error("Can not load the logger configuration file, Please check if the file logger.xml exists on current directory", err)
		os.Exit(1)
	} else {
		log.ReplaceLogger(logger)
		logger.Flush()
	}

}

func doSendMail(event Event) (err error) {
	to := config.Recipients
	emailUser := &EmailUser{config.SmtpUsername, config.SmtpPassword, config.SmtpServer, config.SmtpPort}
	err = smtpWorker.send(emailUser, to, event)
	return err
}

func processEvents() {

	if config == nil {
		log.Critical("config is nil. ")
		cleanup()
		return
	}

	err := mongoWorker.Open(config.MongoHost)

	if err != nil {
		log.Critical(" Error : %s.", err)
		cleanup()
	}

	var events []Event

	events, err = mongoWorker.Fetch()
	if err != nil {
		log.Critical(" Error : %s.", err)
		cleanup()
	}

	log.Debugf("%d events to send.", len(events))
	for _, event := range events {
		log.Tracef("Sending event %s.", event.String())
		err = doSendMail(event)
		if err != nil {
			log.Criticalf(" Error : %s.", err)
			cleanup()
		}
		err = mongoWorker.Delete(event.Id)
		if err != nil {
			log.Critical(" Error : %s.", err)
			cleanup()
		}
		log.Tracef("event %s sent and deleted.", event.String())

	}

	defer mongoWorker.Close()
}

func cleanup() {
	stopPocessEvents <- true
	os.Exit(1)
}

func init() {
	//
	loadLogger()
}

func main() {
	var err error
	config, err = NewConfig()
	if err != nil {
		log.Criticalf("Error : %s.", err)
		return
	}

	smtpWorker = NewSmtpWorker()
	if smtpWorker == nil {
		log.Criticalf("smtpWorker is nil.")
		return
	}
	mongoWorker = NewMongoWorker()

	if mongoWorker == nil {
		log.Criticalf("mongoWorker is nil.")
		return
	}

	durationTestCall := time.Duration(config.Schedule) * time.Second
	//
	stopPocessEvents = schedule(processEvents, durationTestCall)

	for {
		log.Info("Sleeping...")
		time.Sleep(5000 * time.Second) //
	}
}
