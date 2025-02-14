package config

import (
	"database/sql"
	"html/template"
	"log"
	"os"

	model "github.com/demkowo/booking/models"
	sqlclient "github.com/demkowo/booking/utils/sql-client"
)

var (
	Values       ci = &conf{}
	dbConnection    = os.Getenv("DB_CONNECTION")
)

type ci interface {
	Get() *conf
	Set(conf)
}

type conf struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	AuthCache     map[string]*model.Account
	UserCache     map[string]*model.User
	InProduction  bool
	Session       string
	Toast         toast
	JWTSecret     []byte
	DBcli         sqlclient.SqlClient
	DB            *sql.DB
}

type toast struct {
	Active  bool
	Message string
	Success bool
}

func (m *conf) Get() *conf {
	m.JWTSecret = []byte(os.Getenv("JWT_SECRET"))

	return m
}

func (m *conf) Set(c conf) {
	log.Println("=== model Values Set ===")
	m.UseCache = c.UseCache
	m.InProduction = c.InProduction
	m.Session = c.Session
	m.TemplateCache = c.TemplateCache
	m.Toast = c.Toast
}

func getDBclient() sqlclient.SqlClient {
	dbcli, err := sqlclient.Open("postgres", dbConnection)
	if err != nil {
		log.Panic(err)
	}
	defer dbcli.Close()

	return dbcli
}

func getDB() *sql.DB {
	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	return db
}
