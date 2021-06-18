package main

import (
	"github.com/gotway/service-examples/cmd/stock/api"
	"github.com/gotway/service-examples/cmd/stock/redis"
)

func main() {
	redis.Init()
	api.NewAPI()
}
