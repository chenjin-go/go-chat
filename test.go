package main

import (
	"fmt"
	"strconv"
)

func main() {

	int, err := strconv.Atoi("1")
	fmt.Println(int, err)
}
