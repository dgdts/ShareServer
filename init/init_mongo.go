package init

import (
	"github.com/dgdts/UniversalServer/pkg/config"
	"github.com/dgdts/UniversalServer/pkg/mongo"
)

func InitMongo(config *config.GlobalConfig) {
	mongoConfig := &mongo.MongoClient{
		Path:        config.Mongo.Path,
		Username:    config.Mongo.Username,
		Password:    config.Mongo.Password,
		MaxPoolSize: config.Mongo.MaxPoolSize,
		MinPoolSize: config.Mongo.MinPoolSize,
		Database:    config.Mongo.Database,
	}
	mongo.RegisterConnection(mongoConfig)
}
