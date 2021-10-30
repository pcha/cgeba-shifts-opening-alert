package dependencies

import "os"

type Parameters struct {
	MongoURL    string
	MongoDBName string
	BotToken    string
}

func BuildParameters() Parameters {
	return Parameters{
		MongoURL:    os.Getenv("MONGO_URL"),
		MongoDBName: os.Getenv("MONGO_DB"),
		BotToken:    os.Getenv("BOT_TOKEN"),
	}
}
