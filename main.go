package main

import (
	"databases/moss"
	"fmt"

	"github.com/philippgille/gokv"
)

type insight struct {
	SenID string
}

func main() {
	//Create client
	client, err := moss.NewStore(nil)
	if err != nil {
		panic(err)
	}
	interactWithStore(client)

}

// interactWithStore stores, retrieves, prints and deletes a value.
func interactWithStore(store gokv.Store) {
	// Store value
	val := insight{
		SenID: "abc@123",
	}

	if err := store.Set("key", val); err != nil {
		fmt.Println(err)
	}

	// Retrieve value
	retrievedVal := new(insight)
	if found, err := store.Get("key", retrievedVal); err != nil {
		panic(err)
	} else if !found {
		fmt.Println("Value not found")
	}

	fmt.Printf("key: %+v\n", *retrievedVal) // Prints `key: {SenID: abc@123}`

	// Delete value
	if err := store.Delete("key"); err != nil {
		fmt.Println(err)
	}

	// Retrieve value again
	newdata := new(insight)
	if found, err := store.Get("key", newdata); err != nil {
		panic(err)
	} else if !found {
		fmt.Println("Value not found")
	}
	fmt.Printf("After delete ")
	fmt.Printf("key: %+v\n", *newdata) // Prints `key: {SenID:}`

	if err := store.Close(); err != nil {
		fmt.Println(err)
	}
}
