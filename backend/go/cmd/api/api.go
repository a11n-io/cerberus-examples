package main

import (
	"cerberus-examples/env"
	"cerberus-examples/internal/database"
	"cerberus-examples/internal/repositories"
	"cerberus-examples/internal/routes"
	"cerberus-examples/internal/server"
	"cerberus-examples/internal/services"
	"cerberus-examples/internal/utils"
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {
	// App context
	ctx := context.Background()

	// env config
	_env := env.GetEnv(".env.dev")

	// CerberusCode
	//cerberusClient := cerberus.NewClient(_env.CERBERUS_HOST, _env.CERBERUS_API_KEY, _env.CERBERUS_API_SECRET)

	db, err := database.NewDB()
	utils.PanicOnError(err)
	defer func() {
		utils.PanicOnError(db.Close())
	}()
	_, err = db.Exec("PRAGMA foreign_keys=ON")
	utils.PanicOnError(err)

	//cdriver, err := cerberusmigrate.WithInstance(cerberusClient, &cerberusmigrate.Config{})
	//if err != nil {
	//	log.Fatalf("could not get cerberus driver: %v", err.Error())
	//}
	//cm, err := migrate.NewWithDatabaseInstance(
	//	"file://cerberusmigrations", "cerberus", cdriver)
	//if err != nil {
	//	log.Fatalf("could not get cerberus migrate: %v", err.Error())
	//} else {
	//	if err := cm.Up(); err != nil {
	//		log.Println(err)
	//	}
	//	log.Println("cerberus migration done")
	//}

	// migrate
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", "sqlite3", driver)
	if err != nil {
		log.Println(err)
	} else {
		if err := m.Up(); err != nil {
			log.Println(err)
		}
		log.Println("sqlite migration done")
	}

	txProvider := database.NewTxProvider(db)

	userService := services.NewUserService(
		txProvider,
		repositories.NewUserRepo(db),
		repositories.NewAccountRepo(db),
		_env.JWT_SECRET, _env.SALT_ROUNDS /*, cerberusClient*/)

	publicRoutes := publicRoutes(userService)

	privateRoutes := privateRoutes(
		/*cerberusClient,*/
		userService,
		services.NewProjectService(txProvider, repositories.NewProjectRepo(db) /*, cerberusClient*/),
		services.NewSprintService(txProvider, repositories.NewSprintRepo(db) /*, cerberusClient*/),
		services.NewStoryService(txProvider, repositories.NewStoryRepo(db) /*, cerberusClient*/))

	// Run server with context
	webserver := server.NewWebServer(ctx, _env.APP_PORT, _env.JWT_SECRET, publicRoutes, privateRoutes)
	webserver.Start()
}

func publicRoutes(
	authService services.UserService) []routes.Routable {
	return []routes.Routable{
		routes.NewAuthRoutes(authService),
	}
}

func privateRoutes(
	//cerberusClient cerberus.CerberusClient,
	userService services.UserService,
	projectService services.ProjectService,
	sprintService services.SprintService,
	storyService services.StoryService) []routes.Routable {
	return []routes.Routable{
		routes.NewUserRoutes(userService /*, cerberusClient*/),
		routes.NewProjectRoutes(projectService /*, cerberusClient*/),
		routes.NewSprintRoutes(sprintService /*, cerberusClient*/),
		routes.NewStoryRoutes(storyService /*, cerberusClient*/),
	}
}
