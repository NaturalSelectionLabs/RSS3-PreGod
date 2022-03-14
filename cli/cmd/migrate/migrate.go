package migrate

import (
	"context"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/handler"
	mongomodel "github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Timeout     int
	MongoDSN    string
	PostgresDSN string
}

var (
	_ cmd.Command = &Migrate{}
)

type Migrate struct {
	config Config

	mongoClient    *mongo.Client
	postgresClient *gorm.DB
}

func (m *Migrate) Initialize() error {
	// Initialize MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(m.config.MongoDSN))
	if err != nil {
		return err
	}

	m.mongoClient = mongoClient

	logrus.Infoln("Connected to Mongo")

	// Initialize Postgres
	postgresClient, err := gorm.Open(postgres.New(postgres.Config{
		DSN: m.config.PostgresDSN,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return err
	}

	m.postgresClient = postgresClient

	logrus.Infoln("Connected to Postgres")

	return nil
}

func (m *Migrate) Run() error {
	ctx := context.Background()

	fileCollection := m.mongoClient.Database("rss3").Collection("files")

	logrus.Infoln("Begin pulling files")

	cursor, err := fileCollection.Find(ctx, bson.D{})
	if err != nil {
		return err
	}

	var files []mongomodel.File
	for cursor.Next(ctx) { // nolint:wsl // This should be a bug with lint
		var file mongomodel.File
		if err = cursor.Decode(&file); err != nil {
			logrus.Errorln(err)

			continue
		}

		files = append(files, file)
	}

	logrus.Infoln("Begin importing files")

	for _, file := range files {
		// Deprecated
		if strings.Contains(file.Path, "-list-backlinks.following") {
			continue
		}

		if strings.Contains(file.Path, "-list-links.following") {
			if err := handler.MigrateLinkList(m.postgresClient, file); err != nil {
				logrus.Errorln(err)
			}

			continue
		}

		if err := handler.MigrateIndex(m.postgresClient, file); err != nil {
			logrus.Errorln(err)

			continue
		}
	}

	return nil
}
