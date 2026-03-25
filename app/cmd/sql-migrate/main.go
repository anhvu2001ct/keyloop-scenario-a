package main

import (
	"context"
	"log"
	"os"
	"scenario-a/internal/config/env"
	"scenario-a/internal/sqlmigrate"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
)

var logger = log.New(os.Stderr, "[sql-migrate] ", log.LstdFlags)

func newMigrator(mode, engine string) *sqlmigrate.SqlMigration {
	envMode := env.EnvMode(mode)
	if !envMode.Valid() {
		logger.Fatalf("unexpected value for `env`: %s", envMode)
	}

	if !slices.Contains([]string{"postgres"}, engine) {
		logger.Fatalf("unexpected value for `db`: %s", engine)
	}

	return sqlmigrate.New(envMode, engine)
}

func main() {
	cmd := &cli.Command{
		Usage:     "sql migration tool",
		UsageText: "Perform migration on the sql database. The env variables must contains `SQL_DBNAME`, `SQL_HOST_PORT` `SQL_USER`, `SQL_PASSWORD`",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "default",
				Usage: "set the environment. If 'test', will attempt to load '.env.test'",
			},
			&cli.StringFlag{
				Name:  "db",
				Value: "postgres",
				Usage: "set the sql engine. Only support 'postgres'",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "new",
				Usage:     "create new migration file",
				ArgsUsage: "name=<string>",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "name",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					migrationName := strings.TrimSpace(c.StringArg("name"))
					if migrationName == "" {
						logger.Fatal("name must be specified")
					}

					logger.Printf("Using env=%s, db=%s\n", c.String("env"), c.String("db"))

					migrator := newMigrator(c.String("env"), c.String("db"))
					return migrator.NewMigration(migrationName)
				},
			},
			{
				Name:  "up",
				Usage: "apply migrations",
				Action: func(ctx context.Context, c *cli.Command) error {
					logger.Printf("Using env=%s, db=%s\n", c.String("env"), c.String("db"))

					migrator := newMigrator(c.String("env"), c.String("db"))
					return migrator.Up()
				},
			},
			{
				Name:  "down",
				Usage: "de-apply latest migrations",
				Action: func(ctx context.Context, c *cli.Command) error {
					logger.Printf("Using env=%s, db=%s\n", c.String("env"), c.String("db"))

					migrator := newMigrator(c.String("env"), c.String("db"))
					return migrator.Down()
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Fatal(err)
	}
}
