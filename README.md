# Exchango
Read calendar information from an Exchange Server using Go.

[![Build Status](https://travis-ci.org/sgoertzen/xchango.svg?branch=master)](https://travis-ci.org/sgoertzen/xchango)
[![Build Status](https://semaphoreci.com/api/v1/projects/fb260216-9410-469d-a2fd-6920887e3461/431376/badge.svg)](https://semaphoreci.com/sgoertzen/xchango)

## Install
go get github.com/sgoertzen/xchango

## Usage
```sh
import "github.com/sgoertzen/xchango"

func main() {
	xchango.SetExchangeConfig(/* your class that implements ExchangeConfig interface */)
	
	user := ExchangeUser { Username: "sally", Password: "123" }
	cal, err := xchango.GetExchangeCalendar(user)
	if err != nil {
		// handle error
	}
	
	appointments, er := xchango.GetExchangeAppointments(user, cal)
	if er != nil {
		// handle error
	}
	
	for _, app := range appointments {
		// Do something with each appointment
	}
}
```
