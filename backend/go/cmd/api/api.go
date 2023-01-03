package main

import (
	"cerberus-examples/internal/database"
	"cerberus-examples/internal/repositories"
	"cerberus-examples/internal/routes"
	"cerberus-examples/internal/server"
	"cerberus-examples/internal/services"
	"cerberus-examples/internal/utils"
	"context"
	"github.com/golang-migrate/migrate/v4"
	// Add cerberus imports here
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {

	var appPort, jwtSecret string
	var saltRounds int

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "appPort",
				Value:       "8081",
				Usage:       "Port of webserver to listen on",
				Destination: &appPort,
				EnvVars:     []string{"APP_PORT"},
			},
			&cli.StringFlag{
				Name:        "jwtSecret",
				Value:       "secret",
				Usage:       "Secret for signing JWT tokens",
				Destination: &jwtSecret,
				EnvVars:     []string{"JWT_SECRET"},
			},
			&cli.IntFlag{
				Name:        "saltRounds",
				Value:       10,
				Usage:       "Number of salt rounds for password encryption",
				Destination: &saltRounds,
				EnvVars:     []string{"SALT_ROUNDS"},
			},
			// Add cerberus config code here
		},
		Action: func(cCtx *cli.Context) error {

			// App context
			ctx := context.Background()

			// Add cerberus client and migration code here

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
				jwtSecret, saltRounds)

			publicRoutes := publicRoutes(userService)

			privateRoutes := privateRoutes(
				userService,
				services.NewProjectService(txProvider, projectRepo),
				services.NewSprintService(txProvider, sprintRepo),
				services.NewStoryService(txProvider, storyRepo))

			// Run server with context
			webserver := server.NewWebServer(ctx, appPort, jwtSecret, publicRoutes, privateRoutes)
			webserver.Start()

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
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
