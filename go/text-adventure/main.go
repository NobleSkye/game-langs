package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Room struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Exits       map[string]int `json:"exits"`
	NPCs        []int          `json:"npcs"`
	Items       []int          `json:"items"`
}

type NPC struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Dialogue    string `json:"dialogue"`
}

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func loadRooms(filename string) ([]Room, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var rooms []Room
	err = json.Unmarshal(data, &rooms)
	return rooms, err
}

func loadNPCs(filename string) ([]NPC, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var npcs []NPC
	err = json.Unmarshal(data, &npcs)
	return npcs, err
}

func loadItems(filename string) ([]Item, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var items []Item
	err = json.Unmarshal(data, &items)
	return items, err
}

func main() {
	rooms, err := loadRooms("rooms.json")
	if err != nil {
		fmt.Println("Error loading rooms:", err)
		os.Exit(1)
	}

	npcs, err := loadNPCs("npcs.json")
	if err != nil {
		fmt.Println("Error loading NPCs:", err)
		os.Exit(1)
	}

	items, err := loadItems("items.json")
	if err != nil {
		fmt.Println("Error loading items:", err)
		os.Exit(1)
	}

	inventory := make(map[int]Item) // Simple inventory system

	// Game loop
	currentRoom := rooms[0] // Start in the first room
	for {
		fmt.Printf("\nYou are in %s.\n%s\n", currentRoom.Name, currentRoom.Description)
		fmt.Println("Exits:")
		for direction, roomID := range currentRoom.Exits {
			fmt.Printf("- %s to %d\n", direction, roomID)
		}

		// Display NPCs
		fmt.Println("You see:")
		for _, npcID := range currentRoom.NPCs {
			for _, npc := range npcs {
				if npc.ID == npcID {
					fmt.Printf("- %s: %s\n", npc.Name, npc.Description)
				}
			}
		}

		// Display items
		fmt.Println("Items available:")
		for _, itemID := range currentRoom.Items {
			for _, item := range items {
				if item.ID == itemID {
					fmt.Printf("- %s: %s\n", item.Name, item.Description)
				}
			}
		}

		// Player input
		var input string
		fmt.Print("What do you want to do? ")
		fmt.Scanln(&input)

		// Handle movement
		if exitRoomID, exists := currentRoom.Exits[input]; exists {
			for _, room := range rooms {
				if room.ID == exitRoomID {
					currentRoom = room
					break
				}
			}
		} else if input == "talk" {
			for _, npcID := range currentRoom.NPCs {
				for _, npc := range npcs {
					if npc.ID == npcID {
						fmt.Println(npc.Dialogue)
					}
				}
			}
		} else if input == "take lantern" {
			for _, itemID := range currentRoom.Items {
				for _, item := range items {
					if item.ID == itemID && item.Name == "Lantern" {
						inventory[item.ID] = item
						fmt.Println("You have taken the lantern.")
						currentRoom.Items = removeItem(currentRoom.Items, item.ID) // Remove item from room
					}
				}
			}
		} else {
			fmt.Println("You can't go that way or perform that action.")
		}
	}
}

// Utility function to remove an item from a slice
func removeItem(items []int, itemID int) []int {
	for i, id := range items {
		if id == itemID {
			return append(items[:i], items[i+1:]...)
		}
	}
	return items
}
