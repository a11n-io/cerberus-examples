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

	db, err := database.NewDB()
	utils.PanicOnError(err)
	defer func() {
		utils.PanicOnError(db.Close())
	}()
	_, err = db.Exec("PRAGMA foreign_keys=ON")
	utils.PanicOnError(err)

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

	userRepo := repositories.NewUserRepo(db)
	accountRepo := repositories.NewAccountRepo(db)
	projectRepo := repositories.NewProjectRepo(db)
	sprintRepo := repositories.NewSprintRepo(db)
	storyRepo := repositories.NewStoryRepo(db)

	userService := services.NewUserService(
		txProvider,
		userRepo,
		accountRepo,
		_env.JWT_SECRET, _env.SALT_ROUNDS)

	publicRoutes := publicRoutes(userService)

	privateRoutes := privateRoutes(
		userService,
		services.NewProjectService(txProvider, projectRepo),
		services.NewSprintService(txProvider, sprintRepo),
		services.NewStoryService(txProvider, storyRepo))

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
	userService services.UserService,
	projectService services.ProjectService,
	sprintService services.SprintService,
	storyService services.StoryService) []routes.Routable {
	return []routes.Routable{
		routes.NewUserRoutes(userService),
		routes.NewProjectRoutes(projectService),
		routes.NewSprintRoutes(sprintService),
		routes.NewStoryRoutes(storyService),
	}
}
