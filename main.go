package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const familyTreeFile = "family_tree.json"

// Person represents an individual in the family tree.
type Person struct {
	Name      string   `json:"name"`
	Relations []string `json:"relations"`
}

func main() {
	createFamilyTreeFile()

	if len(os.Args) < 2 {
		fmt.Println("Usage: family-tree <command> [options]")
		fmt.Println("\nCommands:")
		fmt.Println("  add person       Add a person to the family tree")
		fmt.Println("  add relationship Add a relationship to a person in the family tree")
		fmt.Println("  connect          Connect two people in the family tree")
		fmt.Println("  countsons        Count the number of sons for an individual")
		fmt.Println("  countdaughters   Count the number of daughters for an individual")
		fmt.Println("  countwives       Count the number of wives for an individual")
		fmt.Println("  father           Find the father of an individual")
		fmt.Println("  help             Show available commands")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Command 'add' requires an additional argument (person or relationship).")
			os.Exit(1)
		}
		subcommand := os.Args[2]
		switch subcommand {
		case "person":
			if len(os.Args) < 4 {
				fmt.Println("Usage: family-tree add person <name>")
				os.Exit(1)
			}
			name := os.Args[3]
			addPerson(name)
		case "relationship":
			if len(os.Args) < 4 {
				fmt.Println("Usage: family-tree add relationship <name>")
				os.Exit(1)
			}
			name := os.Args[3]
			addRelationship(name)
		default:
			fmt.Println("Unknown subcommand for 'add'. Use 'person' or 'relationship'.")
			os.Exit(1)
		}
	case "connect":
		if len(os.Args) < 7 || os.Args[4] != "as" || os.Args[6] != "of" {
			fmt.Println("Usage: family-tree connect <name1> as <relationship> of <name2>")
			os.Exit(1)
		}
		name1 := os.Args[2]
		relationship := os.Args[5]
		name2 := os.Args[7]
		connectPeople(name1, relationship, name2)
	case "countsons":
		if len(os.Args) < 3 {
			fmt.Println("Usage: family-tree countsons <name>")
			os.Exit(1)
		}
		name := os.Args[2]
		count := countSons(name)
		fmt.Printf("%s has %d sons.\n", name, count)
	case "countdaughters":
		if len(os.Args) < 3 {
			fmt.Println("Usage: family-tree countdaughters <name>")
			os.Exit(1)
		}
		name := os.Args[2]
		count := countDaughters(name)
		fmt.Printf("%s has %d daughters.\n", name, count)
	case "countwives":
		if len(os.Args) < 3 {
			fmt.Println("Usage: family-tree countwives <name>")
			os.Exit(1)
		}
		name := os.Args[2]
		count := countWives(name)
		fmt.Printf("%s has %d wives.\n", name, count)
	case "father":
		if len(os.Args) < 4 || os.Args[2] != "of" {
			fmt.Println("Usage: family-tree father of <name>")
			os.Exit(1)
		}
		name := os.Args[3]
		fatherName := findFather(name)
		if fatherName != "" {
			fmt.Printf("Father of %s is %s.\n", name, fatherName)
		} else {
			fmt.Printf("Father of %s is not in the family tree.\n", name)
		}
	case "help":
		fmt.Println("Available commands:")
		fmt.Println("  add person       Add a person to the family tree")
		fmt.Println("  add relationship Add a relationship to a person in the family tree")
		fmt.Println("  connect          Connect two people in the family tree")
		fmt.Println("  countsons        Count the number of sons for an individual")
		fmt.Println("  countdaughters   Count the number of daughters for an individual")
		fmt.Println("  countwives       Count the number of wives for an individual")
		fmt.Println("  father           Find the father of an individual")
		fmt.Println("  help             Show available commands")
	default:
		fmt.Println("Unknown command. Use 'help' to see available commands.")
		os.Exit(1)
	}
}

func createFamilyTreeFile() {
	if _, err := os.Stat(familyTreeFile); os.IsNotExist(err) {
		// Family tree file does not exist, create an empty one
		initialData := make(map[string]Person)
		data, err := json.Marshal(initialData)
		if err != nil {
			fmt.Printf("Error encoding family tree data: %v\n", err)
			os.Exit(1)
		}
		err = writeFamilyTreeFile(data)
		if err != nil {
			fmt.Printf("Error creating family tree file: %v\n", err)
			os.Exit(1)
		}
	}
}

func addPerson(name string) {
	data, err := readFamilyTreeFile()
	if err != nil {
		fmt.Printf("Error reading family tree file: %v\n", err)
		os.Exit(1)
	}

	var familyTree map[string]Person
	err = json.Unmarshal(data, &familyTree)
	if err != nil {
		fmt.Printf("Error decoding family tree data: %v\n", err)
		os.Exit(1)
	}

	if _, exists := familyTree[name]; exists {
		fmt.Printf("%s is already in the family tree.\n", name)
	} else {
		familyTree[name] = Person{Name: name, Relations: []string{}}
		newData, err := json.MarshalIndent(familyTree, "", "  ")
		if err != nil {
			fmt.Printf("Error encoding family tree data: %v\n", err)
			os.Exit(1)
		}

		err = writeFamilyTreeFile(newData)
		if err != nil {
			fmt.Printf("Error writing family tree file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added %s to the family tree.\n", name)
	}
}

func addRelationship(name string) {
	data, err := readFamilyTreeFile()
	if err != nil {
		fmt.Printf("Error reading family tree file: %v\n", err)
		os.Exit(1)
	}

	var familyTree map[string]Person
	err = json.Unmarshal(data, &familyTree)
	if err != nil {
		fmt.Printf("Error decoding family tree data: %v\n", err)
		os.Exit(1)
	}

	if _, exists := familyTree[name]; exists {
		var relation string
		if len(os.Args) >= 5 {
			relation = os.Args[4]
		} else {
			fmt.Printf("Please provide a relationship (e.g., father, son).\n")
			os.Exit(1)
		}

		person := familyTree[name]
		person.Relations = append(person.Relations, relation)
		familyTree[name] = person

		newData, err := json.MarshalIndent(familyTree, "", "  ")
		if err != nil {
			fmt.Printf("Error encoding family tree data: %v\n", err)
			os.Exit(1)
		}

		err = writeFamilyTreeFile(newData)
		if err != nil {
			fmt.Printf("Error writing family tree file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added %s as %s's %s.\n", relation, name, relation)
	} else {
		fmt.Printf("%s is not in the family tree. You can add the person using 'add person' first.\n", name)
	}
}

func connectPeople(name1, relationship, name2 string) {
	data, err := readFamilyTreeFile()
	if err != nil {
		fmt.Printf("Error reading family tree file: %v\n", err)
		os.Exit(1)
	}

	var familyTree map[string]Person
	err = json.Unmarshal(data, &familyTree)
	if err != nil {
		fmt.Printf("Error decoding family tree data: %v\n", err)
		os.Exit(1)
	}

	if person1, exists := familyTree[name1]; exists {
		if person2, exists := familyTree[name2]; exists {
			person1.Relations = append(person1.Relations, relationship)
			familyTree[name1] = person1

			// Add reverse relationship
			// For example, if Amit Dhakad is a son of KK Dhakad, then KK Dhakad is a parent of Amit Dhakad
			person2.Relations = append(person2.Relations, "parent")
			familyTree[name2] = person2

			newData, err := json.MarshalIndent(familyTree, "", "  ")
			if err != nil {
				fmt.Printf("Error encoding family tree data: %v\n", err)
				os.Exit(1)
			}

			err = writeFamilyTreeFile(newData)
			if err != nil {
				fmt.Printf("Error writing family tree file: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Connected %s as %s of %s.\n", name1, relationship, name2)
		} else {
			fmt.Printf("%s is not in the family tree. You can add the person using 'add person' first.\n", name2)
		}
	} else {
		fmt.Printf("%s is not in the family tree. You can add the person using 'add person' first.\n", name1)
	}
}

func countSons(name string) int {
	data, err := readFamilyTreeFile()
	if err != nil {
		fmt.Printf("Error reading family tree file: %v\n", err)
		os.Exit(1)
	}

	var familyTree map[string]Person
	err = json.Unmarshal(data, &familyTree)
	if err != nil {
		fmt.Printf("Error decoding family tree data: %v\n", err)
		os.Exit(1)
	}

	if person, exists := familyTree[name]; exists {
		count := 0
		for _, relation := range person.Relations {
			if relation == "son" {
				count++
			}
		}
		return count
	}

	fmt.Printf("%s is not in the family tree.\n", name)
	os.Exit(1)
	return 0
}

func countDaughters(name string) int {
	data, err := readFamilyTreeFile()
	if err != nil {
		fmt.Printf("Error reading family tree file: %v\n", err)
		os.Exit(1)
	}

	var familyTree map[string]Person
	err = json.Unmarshal(data, &familyTree)
	if err != nil {
		fmt.Printf("Error decoding family tree data: %v\n", err)
		os.Exit(1)
	}

	if person, exists := familyTree[name]; exists {
		count := 0
		for _, relation := range person.Relations {
			if relation == "daughter" {
				count++
			}
		}
		return count
	}

	fmt.Printf("%s is not in the family tree.\n", name)
	os.Exit(1)
	return 0
}

func countWives(name string) int {
	data, err := readFamilyTreeFile()
	if err != nil {
		fmt.Printf("Error reading family tree file: %v\n", err)
		os.Exit(1)
	}

	var familyTree map[string]Person
	err = json.Unmarshal(data, &familyTree)
	if err != nil {
		fmt.Printf("Error decoding family tree data: %v\n", err)
		os.Exit(1)
	}

	if person, exists := familyTree[name]; exists {
		count := 0
		for _, relation := range person.Relations {
			if relation == "wife" {
				count++
			}
		}
		return count
	}

	fmt.Printf("%s is not in the family tree.\n", name)
	os.Exit(1)
	return 0
}

func findFather(name string) string {
	data, err := readFamilyTreeFile()
	if err != nil {
		fmt.Printf("Error reading family tree file: %v\n", err)
		os.Exit(1)
	}

	var familyTree map[string]Person
	err = json.Unmarshal(data, &familyTree)
	if err != nil {
		fmt.Printf("Error decoding family tree data: %v\n", err)
		os.Exit(1)
	}

	if person, exists := familyTree[name]; exists {
		for _, relation := range person.Relations {
			if relation == "father" {
				// Search for the father's name in the family tree
				for key, value := range familyTree {
					if key != name && value.Name == person.Name {
						return key
					}
				}
			}
		}
	}

	// If no father is found, return an empty string
	return ""
}

func readFamilyTreeFile() ([]byte, error) {
	file, err := os.Open(familyTreeFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []byte
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		data = append(data, buffer[:n]...)
	}
	return data, nil
}

func writeFamilyTreeFile(data []byte) error {
	file, err := os.Create(familyTreeFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

