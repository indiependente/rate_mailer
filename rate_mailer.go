package main

import (
	"os"
	"log"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"io/ioutil"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strconv"
)

func main() {

	if len(os.Args) != 4 {
		log.Fatal("rate_mailer expects 3 arguments!\nUsage: rate_mailer GBP EUR /var/log/log.txt")
	}

	logFile, err := os.OpenFile(os.Args[3], os.O_WRONLY, 0666)

	if err != nil {
		panic(err)
	}

	defer logFile.Close()

	// direct all log messages to log.txt
	log.SetOutput(logFile)

	send_mail(request(os.Args[1:3]), os.Args[1:3], "")
	log.Println("email sent.")
}

func request(symbols []string) (float64){
	ip := os.Getenv("HOSTIP")
	resp, err := http.Get("http://" + ip  + ":9090/rates/"+symbols[0]+"_"+symbols[1])
	if err != nil {
		log.Println("GET HTTP request failed")
		return 0.0
	} 
	defer resp.Body.Close()
	
	var d struct {
		Rate float64 `json:"rate"`
		Took string `json:"took"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		log.Println("JSON decoding failed")
		return 0.0
	}
	return d.Rate
	
}

func send_mail(rate float64, sym []string, message string) {
	dat, err := ioutil.ReadFile("credentials")
	if err != nil {
		log.Println("credentials reading failed")
		return
	}
	symbols := [2]string{strings.ToUpper(sym[0]), strings.ToUpper(sym[1])}
	creds := strings.Split(string(dat), "\n")[:3]
	sender := creds[0]
	password := creds[1]
	receiver := creds[2]
	t := time.Now()
	timestamp := strings.Join(strings.Split(t.String(), " ")[:2], ",")
	e := email.NewEmail()
	e.From = sender
	e.To = strings.Split(receiver, " ")
	e.Subject = "Rate " + symbols[0] + " " + symbols[1]
	e.HTML = []byte("<h1>"+symbols[0] +" ➡️ "+ symbols[1] +"<br/>"+strconv.FormatFloat(rate, 'f', -1, 64)+"</h1><br/><h2>"+message+"</h2><br/><span>"+timestamp+"</span>")
	err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", sender, password, "smtp.gmail.com"))
	if err != nil {
		log.Fatal(err)
	}
	return
}
