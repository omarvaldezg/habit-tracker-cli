package main

import (
	"fmt"
	"os"

	"github.com/valdezdata/habit-tracker/pkg/tracker"
)

// version is set during build using ldflags
var version = "development"

func main() {
	// Pass version to the tracker package
	tracker.Version = version

	// Initialize the habit tracker
	habitTracker := tracker.NewHabitTracker()

	// Run the application
	if err := habitTracker.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
