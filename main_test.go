package main

import (
	"testing"
)

func Test_SmtpWorker(t *testing.T) {
	smtpWorker := NewSmtpWorker()
	event := Event{
		AppId:      "vorimport",
		AsteriskId: "asterisk1",
		Type:       2,
		Name:       "EV_STOP",
		Data:       "Application stopped : 2",
	}

	to := []string{"v.gontcharov@gmail.com", "vassili.gontcharov@esifrance.net"}
	emailUser := &EmailUser{"v.gontcharov@gmail.com", "v@s184027_rep", "smtp.gmail.com", 587}

	//emailUser := &EmailUser{"vassili.gontcharov@esifrance.net", "vas184027", "mail.esifrance.net", 465}

	err := smtpWorker.send(emailUser, to, event)
	if err != nil {
		t.Fatalf(" Failed : %s", err)
	}
}

func Test_SmtpWorkerSSL(t *testing.T) {
	smtpWorker := NewSmtpWorker()
	event := Event{
		AppId:      "vorimport",
		AsteriskId: "asterisk1",
		Type:       2,
		Name:       "EV_STOP",
		Data:       "Application stopped : 2",
	}

	to := []string{"v.gontcharov@gmail.com", "vassili.gontcharov@esifrance.net"}
	emailUser := &EmailUser{"v.gontcharov@gmail.com", "v@s184027_rep", "smtp.gmail.com", 465}
	//emailUser := &EmailUser{"vassili.gontcharov@esifrance.net", "vas184027", "mail.esifrance.net", 465}

	err := smtpWorker.send(emailUser, to, event)
	if err != nil {
		t.Fatalf(" Failed : %s", err)
	}
}

func Test_Config(t *testing.T) {

	config, err := NewConfig()
	if err != nil {
		t.Fatalf(" Failed : %s", err)
	}

	if config.SmtpServer != "smtp.gmail.com" {
		t.Fatal(" Failed : Cannot find the username v.gontcharov@gmail.com")
	}

	if len(config.Recipients) == 0 {
		t.Fatal(" Failed : Recipients is empty")
	}
}

func Test_Mongo(t *testing.T) {

	mongoWorker := NewMongoWorker()

	err := mongoWorker.Open("127.0.0.1")

	if err != nil {
		t.Fatalf(" Failed : %s", err)
	}

	var events []Event

	events, err = mongoWorker.Fetch()
	if err != nil {
		t.Fatalf(" Failed : %s", err)
	}

	for _, event := range events {

		t.Logf("Event : %s", event.String())
		err = mongoWorker.Delete(event.Id)

	}

	defer mongoWorker.Close()
}
