package main

import (
	"fmt"

	"github.com/zekroTJA/configoration"
)

func main() {
	c, err := configoration.NewBuilder().
		SetBasePath("./testdata").
		AddJsonFile("test1.json", true).
		AddJsonFile("test2.json", true).
		AddYamlFile("test3.yaml", true).
		AddEnvironmentVariables("TEST_", false).
		Build()

	if err != nil {
		panic(err)
	}

	fmt.Println(
		c.GetSection("g:e").GetBool("f"),
	)
}
