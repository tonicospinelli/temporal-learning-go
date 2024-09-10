package main

import (
	"fmt"
	"learning.temporal/greeting"
	"os"
)

func main() {
	name := os.Args[1]
	greeting := greeting.GreetSomeone(name)
	fmt.Println(greeting)
}
