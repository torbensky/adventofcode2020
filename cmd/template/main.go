package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("this command accepts only 1 argument - the number of the new day")
	}

	dayNum, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("day needs to be a valid number")
	}

	currentDay := time.Now().Day()
	if dayNum-currentDay > 1 {
		log.Fatal("that day is too far in the future")
	}

	for {
		if time.Now().Day() < dayNum {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	fmt.Println(getDayHTML(dayNum))
}

func getDayHTML(day int) []byte {
	url := fmt.Sprintf("https://adventofcode.com/2020/day/%d", day)
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		log.Fatal(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return html
}
