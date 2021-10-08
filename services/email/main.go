package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mongmx/simple-microservice/services/email/events"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, ec, err := newNatsConn()
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	defer ec.Close()

	event := events.New(ec)
	_ = event.CreateEmailCreatedListener()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.Fatal(e.Start(":3000"))
}

func newNatsConn() (*nats.Conn, *nats.EncodedConn, error) {
	nc, err := nats.Connect("nats-srv:4222")
	if err != nil {
		return nil, nil, err
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, nil, err
	}
	return nc, ec, nil
}
