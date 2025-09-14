package main

import (
	"fmt"
	"log"

	"easyblog/config"
	"easyblog/database"
	"easyblog/handlers/common"
	"easyblog/handlers/guest"

	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func main() {
	log.SetPrefix("[EasyBlog] ") // set logging prefix

	// load configs
	cfg := config.New()
	if err := cfg.LoadConfig(); err != nil {
		log.Fatalln(err)
	}

	// connect to postgres
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(db *gorm.DB) {
		if err := database.CloseDB(db); err != nil {
			log.Printf("Failed to close postgres: %v", err)
		}
	}(db)

	// init db (create if not exist)
	if err := database.InitSchema(db, "database/sql/schema.sql"); err != nil {
		log.Panicln(err)
	}

	// add routers and load routes
	router := httprouter.New()
	router.HandleMethodNotAllowed = false
	router.HandleOPTIONS = false
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	common.Routes(router, db)
	guest.Routes(router, db)

	// start server
	addr := fmt.Sprintf(":%d", cfg.ServerCfg.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Panicln(err)
	}
}
