package migrate

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewMigrateCommand() *cobra.Command {
	migrate := Migrate{
		config: Config{
			Timeout: 3, // Default is 3 seconds
		},
	}

	command := &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			if err := migrate.Initialize(); err != nil {
				logrus.Fatalln(err)
			}
			if err := migrate.Run(); err != nil {
				logrus.Fatalln(err)
			}
		},
	}

	command.Flags().StringVar(&migrate.config.MongoDSN, "mongo-dsn", "mongodb://rss3:password@localhost/rss3-prod", "")
	command.Flags().StringVar(&migrate.config.PostgresDSN, "postgres-dsn", "postgresql://rss3:password@localhost:5432/pregod", "")

	return command
}
