package pkg

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func GetWorkingDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func SplitID(s string) int {
	res := strings.Split(s, "/")
	num, _ := strconv.Atoi(res[len(res)-1])
	return num
}
