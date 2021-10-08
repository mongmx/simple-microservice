package main

import (
	"errors"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mongmx/simple-microservice/services/auth/handlers"
	"github.com/mongmx/simple-microservice/services/auth/models"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	nc, ec, err := newNatsConn()
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	defer ec.Close()

	db, err := newPostgresConn()
	if err != nil {
		log.Fatal(err)
	}
	err = migrateDB(db)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("COOKIE_SECRET")))))
	logger, _ := zap.NewProduction()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))

	h := handlers.New(db, ec)
	e.POST("/api/auth/signin", h.SignIn)
	e.POST("/api/auth/signup", h.SignUp)

	e.Use(sessionJwt)
	e.GET("/api/auth/current-user", h.CurrentUser)
	e.POST("/api/auth/signout", h.SignOut)
	e.POST("/api/auth/update-password", h.UpdatePassword)

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

func newPostgresConn() (*gorm.DB, error) {
	postgresServer := "auth-postgres-srv"
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	dsn := "host=" + postgresServer + " user=postgres password=" + postgresPassword + " dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	return nil
}

func sessionJwt(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("pay9.sess", c)
		if err != nil {
			return next(c)
		}
		tokenString, ok := sess.Values["jwt"].(string)
		if !ok {
			return next(c)
		}
		token, err := jwt.ParseWithClaims(
			tokenString,
			&handlers.JwtCustomClaims{},
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(os.Getenv("JWT_KEY")), nil
			},
		)
		if err != nil {
			return next(c)
		}
		claims, ok := token.Claims.(*handlers.JwtCustomClaims)
		if !(ok && token.Valid) {
			return next(c)
		}
		c.Set("currentUser", claims)
		return next(c)
	}
}
