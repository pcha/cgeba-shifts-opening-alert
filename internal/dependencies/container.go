package dependencies

import (
	"cgeba-shift-opening-alerter/internal/bot"
	"cgeba-shift-opening-alerter/internal/bot/tgusers"
)

var container *Container

type Container struct {
	registry map[string]interface{}
	params   Parameters
}

func BuildContainer() *Container {
	if container != nil {
		return container
	}
	container = &Container{
		registry: make(map[string]interface{}),
		params:   BuildParameters(),
	}
	return container
}

//func (c Container) set(key string, dep interface{}) {
//	c[key] = dep
//}
//
//func (c Container) get(key string) (interface{}, error) {
//	dep, ok := c[key]
//	if !ok {
//		return nil, errors.New("dependency not found")
//	}
//	return dep, nil
//}

func (c *Container) GetUsersRepository() (*tgusers.Repository, error) {
	return tgusers.NewRepository(c.params.MongoURL, c.params.MongoDBName)
}

func (c *Container) GetBot() (*bot.Bot, error) {
	repo, err := c.GetUsersRepository()
	if err != nil {
		return nil, err
	}
	return bot.NewBot(c.params.BotToken, repo)
}
