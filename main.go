package main

import (
	"fmt"
	"github.com/0x4D5352/bootdev_blog_aggregator/internal/config"
)

func main() {
	cfg := config.Read()
	fmt.Println("starting:")
	fmt.Println(cfg)
	cfg.SetUser("mussar")
	cfg = config.Read()
	fmt.Println("ending:")
	fmt.Println(cfg)
}
