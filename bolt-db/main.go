package main

import (
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	// here use absolute path
	db, err := bolt.Open("D:\\go\\src\\BlockChain-Learning\\bolt-db\\my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
