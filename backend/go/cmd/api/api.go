package main

import (
	"cerberus-examples/internal/common"
	"cerberus-examples/internal/database"
	"cerberus-examples/internal/repositories"
	"cerberus-examples/internal/routes"
	"cerberus-examples/internal/server"
	"cerberus-examples/internal/services"
	"cerberus-examples/internal/utils"
	"context"
	cerberus "github.com/a11n-io/go-cerberus"
	"github.com/golang-migrate/migrate/v4"
	cerberusmigrate "github.com/golang-migrate/migrate/v4/database/cerberus"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	var appPort, jwtSecret, cerberusHost, cerberusApiKey, cerberusApiSecret string
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
			&cli.StringFlag{
				Name:        "cerberusHost",
				Value:       "https://cerberus-api.a11n.io:443",
				Usage:       "Cerberus API host",
				Destination: &cerberusHost,
				EnvVars:     []string{"CERBERUS_HOST"},
			},
			&cli.StringFlag{
				Name:        "cerberusApiKey",
				Value:       "",
				Usage:       "Cerberus API key for this app",
				Destination: &cerberusApiKey,
				EnvVars:     []string{"CERBERUS_API_KEY"},
			},
			&cli.StringFlag{
				Name:        "cerberusApiSecret",
				Value:       "",
				Usage:       "Cerberus API secret for this app",
				Destination: &cerberusApiSecret,
				EnvVars:     []string{"CERBERUS_API_SECRET"},
			},
		},
		Action: func(cCtx *cli.Context) error {

			// App context
			ctx := context.Background()

			cerberusClient := cerberus.NewClient(cerberusHost, cerberusApiKey, cerberusApiSecret)

			cdriver, err := cerberusmigrate.WithInstance(cerberusClient, &cerberusmigrate.Config{})
			if err != nil {
				log.Fatalf("could not get cerberus driver: %v", err.Error())
			}
			cm, err := migrate.NewWithDatabaseInstance(
				"file://cerberusmigrations", "cerberus", cdriver)
			if err != nil {
				log.Fatalf("could not get cerberus migrate: %v", err.Error())
			} else {
				if err := cm.Up(); err != nil {
					log.Println(err)
				}
				log.Println("cerberus migration done")
			}

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
				jwtSecret, saltRounds, cerberusClient)

			// migrate existing data to cerberus

			accounts, err := accountRepo.FindAll()
			if err != nil {
				log.Fatal(err)
			}

			for _, account := range accounts {

				// admin user
				adminUser, err := userRepo.FindOneByEmail("admin")
				if err != nil {
					log.Fatal(err)
				}

				// Get token
				tokenPair, err := cerberusClient.GetUserToken(account.Id, adminUser.Id)
				if err != nil {
					log.Fatal(err)
				}
				mctx := context.WithValue(ctx, "cerberusTokenPair", tokenPair)

				err = cerberusClient.ExecuteWithCtx(mctx,
					cerberusClient.CreateAccountCmd(account.Id),
					cerberusClient.CreateSuperRoleCmd(common.AccountAdministrator_R),
					cerberusClient.CreateUserCmd(adminUser.Id, adminUser.Email, adminUser.Name),
					cerberusClient.AssignRoleCmd(common.AccountAdministrator_R, adminUser.Id),
					cerberusClient.CreateResourceCmd(account.Id, "", common.Account_RT),
					cerberusClient.CreateRolePermissionCmd(common.AccountAdministrator_R, account.Id, []string{common.CanManageAccount_P}))
				if err != nil {
					log.Fatal(err)
				}

				// all users
				users, err := userRepo.FindAll(account.Id)
				if err != nil {
					log.Fatal(err)
				}

				for _, user := range users {

					err = cerberusClient.ExecuteWithCtx(mctx,
						cerberusClient.CreateUserCmd(user.Id, user.Email, user.Name))
					if err != nil {
						log.Fatal(err)
					}
				}

				// projects
				projects, err := projectRepo.FindByAccount(account.Id)
				if err != nil {
					log.Fatal(err)
				}

				for _, project := range projects {
					err = cerberusClient.ExecuteWithCtx(mctx,
						cerberusClient.CreateResourceCmd(project.Id, account.Id, common.Project_RT))
					if err != nil {
						log.Fatal(err)
					}

					// sprints
					sprints, err := sprintRepo.FindByProject(project.Id)
					if err != nil {
						log.Fatal(err)
					}

					for _, sprint := range sprints {
						err = cerberusClient.ExecuteWithCtx(mctx,
							cerberusClient.CreateResourceCmd(sprint.Id, project.Id, common.Sprint_RT))
						if err != nil {
							log.Fatal(err)
						}

						// stories
						stories, err := storyRepo.FindBySprint(sprint.Id)
						if err != nil {
							log.Fatal(err)
						}

						for _, story := range stories {
							err = cerberusClient.ExecuteWithCtx(mctx,
								cerberusClient.CreateResourceCmd(story.Id, sprint.Id, common.Story_RT))
							if err != nil {
								log.Fatal(err)
							}
						}
					}
				}
			}

			publicRoutes := publicRoutes(userService)

			privateRoutes := privateRoutes(
				cerberusClient,
				userService,
				services.NewProjectService(txProvider, projectRepo, cerberusClient),
				services.NewSprintService(txProvider, sprintRepo, cerberusClient),
				services.NewStoryService(txProvider, storyRepo, cerberusClient))

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
	cerberusClient cerberus.CerberusClient,
	userService services.UserService,
	projectService services.ProjectService,
	sprintService services.SprintService,
	storyService services.StoryService) []routes.Routable {
	return []routes.Routable{
		routes.NewUserRoutes(userService, cerberusClient),
		routes.NewProjectRoutes(projectService, cerberusClient),
		routes.NewSprintRoutes(sprintService, cerberusClient),
		routes.NewStoryRoutes(storyService, cerberusClient),
	}
}
