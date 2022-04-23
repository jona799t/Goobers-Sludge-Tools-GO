package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	accounts, _ := ioutil.ReadFile("data/accounts.txt")
	for _, account := range strings.Split(string(accounts), "\n") {
		fmt.Println(account)
	}
}
