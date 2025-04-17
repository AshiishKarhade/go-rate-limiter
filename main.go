package main

import (
	"fmt"
)

func main() {
	client := InitRedisConnection()
	fmt.Println(client)
}
