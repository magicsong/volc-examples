package display

import (
	"encoding/json"
	"fmt"
	"os"
)

// PrintAsJSON formats any object as JSON and prints it to stdout
// It uses indentation for better readability
func PrintAsJSON[T any](v T) {
	PrintAsJSONWithIndent(v, "  ")
}

// PrintAsJSONWithIndent formats any object as JSON with custom indentation and prints it to stdout
func PrintAsJSONWithIndent[T any](v T, indent string) {
	data, err := json.MarshalIndent(v, "", indent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling to JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

// PrintAsJSONOrError formats any object as JSON and prints it to stdout
// If the marshalling fails, it prints the error and returns it
func PrintAsJSONOrError[T any](v T) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling to JSON: %v\n", err)
		return err
	}
	fmt.Println(string(data))
	return nil
}