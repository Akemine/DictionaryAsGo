package main

import (
	"flag"
	"fmt"
	"os"

	"dictionnary.com/dictionary"
)

func main() {

	action := flag.String("action", "list", "Action to perform on the dictionary")
	d, err := dictionary.New("./badger")
	handleError(err)
	defer d.Close()

	flag.Parse()
	switch *action {
	case "list":
		actionList(d)
	case "add":
		actionAdd(d, flag.Args())
	case "define":
		actionDefine(d, flag.Args())
	case "remove":
		actionRemove(d, flag.Args())
	case "tag":
		actionTag(d, flag.Args())
	default:
		fmt.Println("Unknown command : ", *action)
	}

}

// actionList retrieves and prints all entries from the dictionary.
func actionList(d *dictionary.Dictionary) {
	sortedList, mapList, err := d.List()
	handleError(err)
	for _, word := range sortedList {
		fmt.Println(mapList[word])
	}
}

// actionAdd adds a new word with its definition and tag to the dictionary.
func actionAdd(d *dictionary.Dictionary, args []string) {
	word := args[0]
	definition := args[1]
	// Check if the tag provided is valid
	if !dictionary.IsValidTag(dictionary.Tag(args[2])) {
		fmt.Printf("Invalid tag: %v\n", dictionary.Tag(args[2]))
	} else {
		tag := dictionary.Tag(args[2])
		err := d.Add(word, definition, tag) // Add the word, definition, and tag to the dictionary
		handleError(err)
		fmt.Printf("Word %s added to the dictionary\n Definition: %s\n Tag: %s\n", word, definition, tag)
	}
}

// actionDefine retrieves and prints the definition of a word from the dictionary.
func actionDefine(d *dictionary.Dictionary, args []string) {
	word := args[0]
	definition, err := d.Get(word) // Get the definition of the word from the dictionary
	handleError(err)
	fmt.Printf("Definition of %s is: %s\n", word, definition.Definition)
}

// actionTag retrieves and prints all entries associated with a specific tag from the dictionary.
func actionTag(d *dictionary.Dictionary, args []string) {
	// Check if the tag provided is valid
	if !dictionary.IsValidTag(dictionary.Tag(args[0])) {
		fmt.Printf("Invalid tag: %v\n", dictionary.Tag(args[0]))
	} else {
		tagChoosed := dictionary.Tag(args[0])
		sortedList, mapList, err := d.ListTag(tagChoosed)
		handleError(err)
		for _, word := range sortedList {
			fmt.Printf("%s tag: %s\n", tagChoosed, mapList[word])
		}
	}
}

// actionRemove removes a word from the dictionary.
func actionRemove(d *dictionary.Dictionary, args []string) {
	word := args[0]
	err := d.Remove(word) // Remove the word from the dictionary
	handleError(err)
	fmt.Printf("Word %s has been removed from the dictionary\n", word)
}

// handleError prints an error message and exits the program if an error occurs.
func handleError(err error) {
	if err != nil {
		fmt.Printf("Dictionary error: %v\n", err)
		os.Exit(1) // Exit the program with status code 1
	}
}
