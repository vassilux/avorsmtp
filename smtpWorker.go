package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	//"labix.org/v2/mgo/bson"
	"log"
	"math/big"
	//"net"
	"net/mail"
	"net/smtp"

	//"strconv"
	"html/template"
	"strings"
)

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

type SmtpTemplateData struct {
	From    string
	Subject string
	Body    string
}

const (
	emailTemplate = `
Message : {{ .}}

Sincerely,
Barak Oubama
`

	jsonTemplate = `
Message : {{ .}}


Sincerely,
Barak Oubama

`
	EMAIL_SUBJECT = "AVOR event notification"
)

// Code below from http://golang.org/src/pkg/crypto/tls/handshake_server_test.go

func bigFromString(s string) *big.Int {
	ret := new(big.Int)
	ret.SetString(s, 10)
	return ret
}

func fromHex(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

var testRSACertificate = fromHex("308202b030820219a00302010202090085b0bba48a7fb8ca300d06092a864886f70d01010505003045310b3009060355040613024155311330110603550408130a536f6d652d53746174653121301f060355040a1318496e7465726e6574205769646769747320507479204c7464301e170d3130303432343039303933385a170d3131303432343039303933385a3045310b3009060355040613024155311330110603550408130a536f6d652d53746174653121301f060355040a1318496e7465726e6574205769646769747320507479204c746430819f300d06092a864886f70d010101050003818d0030818902818100bb79d6f517b5e5bf4610d0dc69bee62b07435ad0032d8a7a4385b71452e7a5654c2c78b8238cb5b482e5de1f953b7e62a52ca533d6fe125c7a56fcf506bffa587b263fb5cd04d3d0c921964ac7f4549f5abfef427100fe1899077f7e887d7df10439c4a22edb51c97ce3c04c3b326601cfafb11db8719a1ddbdb896baeda2d790203010001a381a73081a4301d0603551d0e04160414b1ade2855acfcb28db69ce2369ded3268e18883930750603551d23046e306c8014b1ade2855acfcb28db69ce2369ded3268e188839a149a4473045310b3009060355040613024155311330110603550408130a536f6d652d53746174653121301f060355040a1318496e7465726e6574205769646769747320507479204c746482090085b0bba48a7fb8ca300c0603551d13040530030101ff300d06092a864886f70d010105050003818100086c4524c76bb159ab0c52ccf2b014d7879d7a6475b55a9566e4c52b8eae12661feb4f38b36e60d392fdf74108b52513b1187a24fb301dbaed98b917ece7d73159db95d31d78ea50565cd5825a2d5a5f33c4b6d8c97590968c0f5298b5cd981f89205ff2a01ca31b9694dda9fd57e970e8266d71999b266e3850296c90a7bdd9")

var testRSAPrivateKey = &rsa.PrivateKey{
	PublicKey: rsa.PublicKey{
		N: bigFromString("131650079503776001033793877885499001334664249354723305978524647182322416328664556247316495448366990052837680518067798333412266673813370895702118944398081598789828837447552603077848001020611640547221687072142537202428102790818451901395596882588063427854225330436740647715202971973145151161964464812406232198521"),
		E: 65537,
	},
	D: bigFromString("29354450337804273969007277378287027274721892607543397931919078829901848876371746653677097639302788129485893852488285045793268732234230875671682624082413996177431586734171663258657462237320300610850244186316880055243099640544518318093544057213190320837094958164973959123058337475052510833916491060913053867729"),
	Primes: []*big.Int{
		bigFromString("11969277782311800166562047708379380720136961987713178380670422671426759650127150688426177829077494755200794297055316163155755835813760102405344560929062149"),
		bigFromString("10998999429884441391899182616418192492905073053684657075974935218461686523870125521822756579792315215543092255516093840728890783887287417039645833477273829"),
	},
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

type fakeAuth struct {
	smtp.Auth
}

func (a fakeAuth) Start(server *smtp.ServerInfo) (string, []byte,
	error) {
	server.TLS = true
	log.Printf("fakeAuth server name %s \n", server.Name)
	return a.Auth.Start(server)
}

type SmtpWorker struct {
}

func NewSmtpWorker() *SmtpWorker {
	smtpWorker := &SmtpWorker{}

	return smtpWorker
}

func (smtpWorker *SmtpWorker) buildEmailMessage(event Event) (message []byte, err error) {
	var docJson bytes.Buffer
	tj := template.New("jsonTemplate")
	tj, err = tj.Parse(jsonTemplate)
	if err != nil {
		log.Print("error trying to parse mail template")
		return nil, err
	}
	err = tj.ExecuteTemplate(&docJson, "jsonTemplate", template.HTML(event.Json()))
	if err != nil {
		log.Print("error trying to execute mail template")
		return nil, err
	}

	return docJson.Bytes(), nil
}

func (smtpWorker *SmtpWorker) send(emailUser *EmailUser, to []string, event Event) (err error) {
	if emailUser.Port == 465 {
		return smtpWorker.sendMailSSL(emailUser, to, event)
	} else {
		return smtpWorker.sendMail(emailUser, to, event)
	}
}

func (smtpWorker *SmtpWorker) sendMail(emailUser *EmailUser, to []string, event Event) (err error) {

	auth := smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer,
	)

	emailBody, err := smtpWorker.buildEmailMessage(event)
	if err != nil {
		log.Print("Failed build message from template.")
		return err
	}

	serverAddr := fmt.Sprintf("%s:%d", emailUser.EmailServer, emailUser.Port)

	header := make(map[string]string)
	header["Subject"] = encodeRFC2047(EMAIL_SUBJECT)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString(emailBody)

	err = smtp.SendMail(serverAddr, auth, emailUser.Username,
		to, []byte(message))
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (smtpWorker *SmtpWorker) sendMailSSL(emailUser *EmailUser, to []string, event Event) (err error) {
	//
	/*auth := fakeAuth{smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer,
	)}*/

	auth := smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer,
	)
	var emailBody []byte
	emailBody, err = smtpWorker.buildEmailMessage(event)
	if err != nil {
		log.Print("Failed build message from template.")
		return err
	}

	tlsconfig := new(tls.Config)
	tlsconfig.Certificates = make([]tls.Certificate, 1)
	tlsconfig.Certificates[0].Certificate = [][]byte{testRSACertificate}
	tlsconfig.Certificates[0].PrivateKey = testRSAPrivateKey
	tlsconfig.CipherSuites = []uint16{tls.TLS_RSA_WITH_RC4_128_SHA}
	tlsconfig.InsecureSkipVerify = false
	tlsconfig.MinVersion = tls.VersionSSL30
	tlsconfig.MaxVersion = tls.VersionTLS10

	serverAddr := fmt.Sprintf("%s:%d", emailUser.EmailServer, emailUser.Port)

	conn, err := tls.Dial("tcp", serverAddr, tlsconfig)
	if err != nil {
		log.Printf("Error Dialing %s\n", err)
		return err
	}

	client, err := smtp.NewClient(conn, emailUser.EmailServer)
	if err != nil {
		log.Printf("Error SMTP connection: %s\n", err)
		return err
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(tlsconfig); err != nil {
			log.Printf("Error performing StartTLS: %s\n", err)
			return err
		}
	}

	if ok, _ := client.Extension("AUTH"); ok {
		if err := client.Auth(auth); err != nil {
			log.Printf("Error during AUTH %s\n", err)
			return err
		}
	}

	if err = client.Mail(emailUser.Username); err != nil {
		log.Printf("Error: %s\n", err)
		return err
	}

	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			log.Printf("Error: %s\n", err)
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		log.Printf("Error: %s\n", err)
		return err
	}

	header := make(map[string]string)
	header["Subject"] = encodeRFC2047(EMAIL_SUBJECT)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString(emailBody)

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Printf("Error: %s\n", err)
		return err

	}

	err = w.Close()
	if err != nil {
		log.Printf("Error: %s\n", err)
		return err

	}

	client.Quit()

	return nil
}
