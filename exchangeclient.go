package xchango

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"
)

type ExchangeConfig interface {
	ExchangeURL() string
	MaxFetchSize() int
	UserDomain() string
	ExchangeVersion() string
	LookAheadDays() int
}

type exchangeVersion interface {
	FolderRequest() string
	CalendarRequest() string
	CalendarDetailRequest() string
}

type ExchangeUser struct {
	Username string
	Password string
}

type ExchangeCalendar struct {
	Folderid  string
	Changekey string
}

var exchangeConfig ExchangeConfig
var version exchangeVersion
var host string

func SetExchangeConfig(config ExchangeConfig) {
	exchangeConfig = config

	switch config.ExchangeVersion() {
	case "2006":
		version = exchange2006{}
	default:
		panic(errors.New(fmt.Sprintf("Unsupported exchange version %s.  Current supported versions are: 2006", config.ExchangeVersion())))
	}

	u, err := url.Parse(config.ExchangeURL())
	if err != nil {
		panic(err)
	}
	host = u.Host
}

func GetExchangeCalendar(user *ExchangeUser) (*ExchangeCalendar, error) {
	soapReq := version.FolderRequest()
	results, err := postContents([]byte(soapReq), user)

	if err != nil {
		return &ExchangeCalendar{}, err
	}

	log.Printf("Exchange server return results of: %s", results)
	item := parseCalendarFolder(string(results))

	cal := ExchangeCalendar{
		Folderid:  item.Id,
		Changekey: item.ChangeKey,
	}
	return &cal, nil
}

func GetExchangeAppointments(user *ExchangeUser, cal *ExchangeCalendar) (*[]Appointment, error) {
	// This first call will just get ids for each appt
	calRequest := buildCalendarItemRequest(cal.Folderid, cal.Changekey)
	calResults, err := postContents(calRequest, user)
	if err != nil {
		log.Println("Error while getting soap response:", err)
		return nil, err
	}
	itemIds := parseAppointments(calResults)

	// This call will get all the fields given the ids
	appRequest := buildCalendarDetailRequest(itemIds)
	appResults, err := postContents(appRequest, user)
	if err != nil {
		return nil, err
	}

	appointments := parseAppointments(appResults)
	if err != nil {
		return nil, err
	}

	return &appointments, nil
}

func buildCalendarItemRequest(folderid string, changekey string) []byte {

	days := exchangeConfig.LookAheadDays()
	if days < 1 {
		days = 14
	}
	startDate := time.Now().UTC().Format(time.RFC3339)
	endDate := time.Now().UTC().AddDate(0, 0, days).Format(time.RFC3339)

	data := struct {
		StartDate    string
		EndDate      string
		FolderId     string
		ChangeKey    string
		MaxFetchSize int
	}{
		startDate,
		endDate,
		folderid,
		changekey,
		exchangeConfig.MaxFetchSize(),
	}

	t, err := template.New("cal").Parse(version.CalendarRequest())
	if err != nil {
		log.Println("Error while parsing template for item request", err)
	}
	var doc bytes.Buffer
	t.Execute(&doc, data)
	if err != nil {
		log.Println("Error while building contents ", err)
	}

	return doc.Bytes()
}

func buildCalendarDetailRequest(itemIds []Appointment) []byte {

	data := struct {
		Appointments []Appointment
	}{
		itemIds,
	}

	t, err := template.New("detail").Parse(version.CalendarDetailRequest())
	if err != nil {
		log.Println("Error while parsing for detail request", err)
	}
	var doc bytes.Buffer
	t.Execute(&doc, data)
	return doc.Bytes()
}

func postContents(contents []byte, user *ExchangeUser) (string, error) {
	req2, err := http.NewRequest("POST", exchangeConfig.ExchangeURL(), bytes.NewBuffer(contents))

	req2.Header.Set("Host", user.Username+"@"+host)
	req2.Header.Set("Content-Type", "text/xml")
	req2.SetBasicAuth(exchangeConfig.UserDomain()+"/"+user.Username, user.Password)

	// TODO, allow client to be injected by tests!
	client := &http.Client{}
	response, err := client.Do(req2)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(content), nil
}
