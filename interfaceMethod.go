package main

import "fmt"

// 1. Define the interface (the "contract")
type Speaker interface {
	Speak() string
}

type Human struct {
	Name string
}

type Dog struct {
	Breed string
}

func main() {
	bayo := Human{"Bayo"}
	rex := Dog{"German Shepherd"}

	Announce(bayo) // Prints: Announcement: Hello, my name is Bayo
	Announce(rex)  // Prints: Announcement: Woof! I am a German Shepherd
}

func Announce(s Speaker) {
	// This calls the specific Speak() version for whoever was passed in
	fmt.Println("Announcement:", s.Speak())
}

// 2. Human implements Speaker by having a Speak() method
func (h Human) Speak() string {
	return "Hello, my name is " + h.Name
}

// 3. Dog also implements Speaker
func (d Dog) Speak() string {
	return "Woof! I am a " + d.Breed
}

