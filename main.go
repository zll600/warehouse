package main

import (
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

func main() {
	log.Println("Enter main function")
	if _, err := os.Stat("./warehouse.db"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			initdb()
		} else {
			log.Fatal(err)
		}
	}

	if len(os.Args) != 2 {
		log.Fatalln("usage warehouse user-id")
	}

	userID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	printAllData(userID)
	log.Println(userID)
}
