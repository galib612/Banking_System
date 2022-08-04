package main

import (
	"database/sql"
	"log"

	"github.com/galib612/simplebank/api"
	authdb "github.com/galib612/simplebank/authdb/sqlc"
	db "github.com/galib612/simplebank/db/sqlc"
	"github.com/galib612/simplebank/util"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

/*
func authMariadbSchemaMigration(migrateType int) error {
	db, err := sql.Open("mysql", "root:secret@tcp(localhost:3306)/simple_bank?multiStatements=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal("withInstance:", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://authdb/migration",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("NewWithDatabaseInstance:", err.Error())
	}
	if migrateType > 0 {
		if err = m.Up(); err != nil {
			fmt.Print(err, ",trying to down schema...:")
			m.Down()
			err = m.Up()
		}
	} else {
		err = m.Down()

	}

	return err
}
*/

// @title           Simple Bank Api
// @version         1.0
// @description     This is sample simple bank client server. Its developed in golang environment and used Postgres Db. For Routing gin web frame work is used in this Rest Api.
// @termsOfService  www.nokia.com

// @contact.name   API Support
// @contact.url    https://github.com/galib612/simplebank
// @contact.email  mohammd.galib@nokia.com

// @license.name  Nokia
// @license.url   www.nokia.com

// @securityDefinitions.basic  BasicAuth
// @host      localhost:8080

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load the config file:", err)
	}

	postgresconn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to postgres db:", err)
	}

	authdbconn, err := sql.Open(config.AuthDBDriver, config.AuthDBSource)
	if err != nil {
		log.Fatal("cannot connect to postgres db:", err)
	}

	postgresstore := db.NewStore(postgresconn)
	authdbstore := authdb.NewStore(authdbconn)
	server := api.NewServer(postgresstore, authdbstore)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start the server:", err)
	}
}
