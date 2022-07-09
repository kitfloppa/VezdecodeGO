package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	timefile, err := ioutil.ReadFile("input.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lines := strings.Split(string(timefile), "\n")

	for _, str := range lines {
		str = strings.Trim(str, "\r") // For windows OS
		t, err := time.ParseDuration(str)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		time.Sleep(t)
		fmt.Println(t.String())
	}

}
