package main

import (
	"fmt"

	"github.com/rmntim/ozon-task/intenal/config"
)

func main() {
	cfg, dbCfg := config.MustLoad()

	fmt.Println(cfg)
	fmt.Println(dbCfg)
}
