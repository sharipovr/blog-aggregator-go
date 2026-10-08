package main

import (
	"fmt"
	"log"

	"github.com/sharipovr/blog-aggregator-go/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	err = c.SetUser("Rustem Sharipov")
	if err != nil {
		log.Fatal(err)
	}
	c, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.DbUrl)
	fmt.Println(c.CurrentUserName)
}
