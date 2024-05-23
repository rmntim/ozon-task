package main

import (
	"fmt"

	"github.com/rmntim/ozon-task/intenal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg, dbCfg := config.MustLoad()

	fmt.Println(cfg)
	fmt.Println(dbCfg)
}
