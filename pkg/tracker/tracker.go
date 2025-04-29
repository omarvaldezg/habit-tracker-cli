package tracker

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Version can be set during build
var Version = "development"

// Habit represents a single habit to track
type Habit struct {
	Name       string          `json:"name"`
	Color      string          `json:"color"`
	DayResults map[string]bool `json:"days"` // Maps day names to completion status
}

// HabitTracker is the main application state
type HabitTracker struct {
	Habits  []Habit            `json:"habits"`
	Days    []string           `json:"days"`
	DataDir string             `json:"-"`
	App     *tview.Application `json:"-"`
	Table   *tview.Table       `json:"-"`
	Modal   *tview.Modal       `json:"-"`
}

// Colors for each habit
var colorMap = map[string]tcell.Color{
	"water":         tcell.ColorBlue,
	"exercise":      tcell.ColorRed,
	"certification": tcell.ColorYellow,
	"breath":        tcell.ColorWhite,
	"newsboat":      tcell.ColorOrange,
	"recap":         tcell.ColorLightBlue,
	"personal":      tcell.ColorGreen,
	"read":          tcell.ColorPurple,
}

// NewHabitTracker initializes a new habit tracker
func NewHabitTracker() *HabitTracker {
	// Get home directory for storing data
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	dataDir := filepath.Join(homeDir, ".habit-tracker")

	// Create data directory if it doesn't exist
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.Mkdir(dataDir, 0755)
	}

	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	// Default habits based on user's list - exactly matching the colorMap
	defaultHabits := []Habit{
		{Name: "water", Color: "blue", DayResults: make(map[string]bool)},
		{Name: "exercise", Color: "red", DayResults: make(map[string]bool)},
		{Name: "certification", Color: "yellow", DayResults: make(map[string]bool)},
		{Name: "breath", Color: "white", DayResults: make(map[string]bool)},
		{Name: "newsboat", Color: "orange", DayResults: make(map[string]bool)},
		{Name: "recap", Color: "lightblue", DayResults: make(map[string]bool)},
		{Name: "personal", Color: "green", DayResults: make(map[string]bool)},
		{Name: "read", Color: "purple", DayResults: make(map[string]bool)},
	}

	return &HabitTracker{
		Habits:  defaultHabits,
		Days:    days,
		DataDir: dataDir,
		App:     tview.NewApplication(),
	}
}

// LoadData loads saved habit data for the current week
func (ht *HabitTracker) LoadData() error {
	// Get the current week number and year
	year, week := time.Now().ISOWeek()
	filename := filepath.Join(ht.DataDir, fmt.Sprintf("habits_%d_%d.json", year, week))

	// Check if the file exists, if not, create a new one with default habits
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Initialize empty habit data
		for i := range ht.Habits {
			ht.Habits[i].DayResults = make(map[string]bool)
		}
		return ht.SaveData()
	}

	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Parse the habit data
	var loadedHabits []Habit
	if err := json.Unmarshal(data, &loadedHabits); err != nil {
		return err
	}

	// Update our habits with any loaded habits
	ht.Habits = loadedHabits

	// Ensure colors are up to date
	for i, habit := range ht.Habits {
		// Update color if it has changed in the colorMap
		for colorName := range colorMap {
			if habit.Name == colorName {
				switch colorName {
				case "water":
					ht.Habits[i].Color = "blue"
				case "exercise":
					ht.Habits[i].Color = "red"
				case "certification":
					ht.Habits[i].Color = "yellow"
				case "breath":
					ht.Habits[i].Color = "white"
				case "newsboat":
					ht.Habits[i].Color = "orange"
				case "recap":
					ht.Habits[i].Color = "lightblue"
				case "personal":
					ht.Habits[i].Color = "green"
				case "read":
					ht.Habits[i].Color = "purple"
				}
			}
		}
	}

	return nil
}

// SaveData saves the habit data for the current week
func (ht *HabitTracker) SaveData() error {
	year, week := time.Now().ISOWeek()
	filename := filepath.Join(ht.DataDir, fmt.Sprintf("habits_%d_%d.json", year, week))

	data, err := json.MarshalIndent(ht.Habits, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// AddHabit adds a new habit to track
func (ht *HabitTracker) AddHabit(name, color string) {
	newHabit := Habit{
		Name:       name,
		Color:      color,
		DayResults: make(map[string]bool),
	}
	ht.Habits = append(ht.Habits, newHabit)
	ht.SaveData()
	ht.BuildUI() // Rebuild UI to display new habit
}

// RemoveHabit removes a habit by name
func (ht *HabitTracker) RemoveHabit(name string) {
	for i, habit := range ht.Habits {
		if habit.Name == name {
			ht.Habits = append(ht.Habits[:i], ht.Habits[i+1:]...)
			break
		}
	}
	ht.SaveData()
	ht.BuildUI() // Rebuild UI to update after removal
}

// ShowAddHabitDialog displays a dialog to add a new habit
func (ht *HabitTracker) ShowAddHabitDialog() {
	// Create an input form
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Add New Habit")
	form.SetBackgroundColor(tcell.ColorDefault)

	// Add input fields
	var name, color string
	form.AddInputField("Habit Name", "", 20, nil, func(text string) {
		name = text
	})

	form.AddDropDown("Color", []string{"blue", "red", "green", "yellow", "white", "orange", "purple", "lightblue", "lightgreen"}, 0, func(option string, index int) {
		color = option
	})

	// Add buttons
	form.AddButton("Save", func() {
		if name != "" {
			ht.AddHabit(name, color)
			ht.App.SetRoot(ht.Table, true)
		}
	})

	form.AddButton("Cancel", func() {
		ht.App.SetRoot(ht.Table, true)
	})

	// Set the form as application root
	ht.App.SetRoot(form, true)
}

// ShowRemoveHabitDialog displays a dialog to remove a habit
func (ht *HabitTracker) ShowRemoveHabitDialog() {
	if len(ht.Habits) == 0 {
		return // Nothing to remove
	}

	// Create a list to select a habit to remove
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Select Habit to Remove")
	list.SetBackgroundColor(tcell.ColorDefault)

	// Add all habits to the list
	for _, habit := range ht.Habits {
		list.AddItem(habit.Name, "", 0, nil)
	}

	// Set selected callback
	list.SetSelectedFunc(func(index int, name string, secondName string, shortcut rune) {
		// Confirm deletion
		modal := tview.NewModal()
		modal.SetText(fmt.Sprintf("Are you sure you want to remove '%s'?", name))
		modal.SetBackgroundColor(tcell.ColorDefault)
		modal.AddButtons([]string{"Yes", "No"})
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				ht.RemoveHabit(name)
			}
			ht.App.SetRoot(ht.Table, true)
		})
		ht.App.SetRoot(modal, true)
	})

	// Add a cancel option
	list.AddItem("Cancel", "", 0, func() {
		ht.App.SetRoot(ht.Table, true)
	})

	// Set the list as application root
	ht.App.SetRoot(list, true)
}

// ToggleHabit toggles the completion status of a habit for a specific day
func (ht *HabitTracker) ToggleHabit(habitIndex int, day string) {
	if habitIndex < 0 || habitIndex >= len(ht.Habits) {
		return
	}

	current := ht.Habits[habitIndex].DayResults[day]
	ht.Habits[habitIndex].DayResults[day] = !current
	ht.SaveData()
}

// getColorByName converts a color name to tcell.Color
func getColorByName(colorName string) tcell.Color {
	switch colorName {
	case "blue":
		return tcell.ColorBlue
	case "red":
		return tcell.ColorRed
	case "yellow":
		return tcell.ColorYellow
	case "white":
		return tcell.ColorWhite
	case "orange":
		return tcell.ColorOrange
	case "purple":
		return tcell.ColorPurple
	case "green":
		return tcell.ColorGreen
	case "lightblue":
		return tcell.ColorLightBlue
	case "lightgreen":
		return tcell.ColorLightGreen
	default:
		return tcell.ColorWhite
	}
}

// BuildUI builds the terminal UI
func (ht *HabitTracker) BuildUI() {
	// Create a simple table without borders and transparent background
	ht.Table = tview.NewTable()
	ht.Table.SetBorders(false)
	ht.Table.SetEvaluateAllRows(true)
	ht.Table.SetFixed(1, 1)
	ht.Table.SetSelectable(true, true)
	ht.Table.SetBackgroundColor(tcell.ColorDefault) // Transparent background

	// Add top-left header cell
	headerCell := tview.NewTableCell("Habit")
	headerCell.SetTextColor(tcell.ColorWhite)
	headerCell.SetAlign(tview.AlignCenter)
	headerCell.SetSelectable(false)
	headerCell.SetExpansion(1)
	headerCell.SetBackgroundColor(tcell.ColorDefault) // Transparent background
	ht.Table.SetCell(0, 0, headerCell)

	// Add day headers
	for col, day := range ht.Days {
		dayCell := tview.NewTableCell(day)
		dayCell.SetTextColor(tcell.ColorWhite)
		dayCell.SetAlign(tview.AlignCenter)
		dayCell.SetSelectable(false)
		dayCell.SetExpansion(1)
		dayCell.SetBackgroundColor(tcell.ColorDefault) // Transparent background
		ht.Table.SetCell(0, col+1, dayCell)
	}

	// Add habit rows with extra spacing
	for row, habit := range ht.Habits {
		// Add habit name
		habitCell := tview.NewTableCell("  " + habit.Name + "  ")
		habitCell.SetTextColor(getColorByName(habit.Color))
		habitCell.SetAlign(tview.AlignLeft)
		habitCell.SetExpansion(1)
		habitCell.SetBackgroundColor(tcell.ColorDefault) // Transparent background
		ht.Table.SetCell(row*2+1, 0, habitCell)

		// Add status cells for each day
		for col, day := range ht.Days {
			cell := tview.NewTableCell("")
			if completed, exists := habit.DayResults[day]; exists && completed {
				cell.SetText("✓")
				cell.SetTextColor(getColorByName(habit.Color))
			} else {
				cell.SetText("✗")
				cell.SetTextColor(tcell.ColorGray)
			}
			cell.SetAlign(tview.AlignCenter)
			cell.SetExpansion(1)
			cell.SetBackgroundColor(tcell.ColorDefault) // Transparent background
			ht.Table.SetCell(row*2+1, col+1, cell)
		}

		// Add an empty row for spacing
		if row < len(ht.Habits)-1 {
			for col := 0; col <= len(ht.Days); col++ {
				spacerCell := tview.NewTableCell("")
				spacerCell.SetSelectable(false)
				spacerCell.SetBackgroundColor(tcell.ColorDefault) // Transparent background
				ht.Table.SetCell(row*2+2, col, spacerCell)
			}
		}
	}

	// Set vim-like keybindings for navigation with underline highlighting instead of background
	style := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Bold(true).
		Underline(true).
		Background(tcell.ColorDefault) // Keep background transparent
	ht.Table.SetSelectedStyle(style)

	ht.Table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			ht.App.Stop()
		}
	})

	ht.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, col := ht.Table.GetSelection()

		// Vim-like navigation
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'h': // Left
				// Always allow moving to the habit column
				if col > 0 {
					ht.Table.Select(row, col-1)
				}
				return nil
			case 'j': // Down
				// Check if we're on a content row and not at the end
				if row < ht.Table.GetRowCount()-1 {
					// Skip empty spacing rows
					if row%2 == 0 {
						ht.Table.Select(row+1, col)
					} else {
						ht.Table.Select(row+2, col)
					}
				}
				return nil
			case 'k': // Up
				// Skip empty spacing rows
				if row > 1 {
					if row%2 == 0 {
						ht.Table.Select(row-1, col)
					} else {
						ht.Table.Select(row-2, col)
					}
				}
				return nil
			case 'l': // Right
				if col < len(ht.Days) {
					ht.Table.Select(row, col+1)
				}
				return nil
			case ' ', 'x': // Toggle with space or x
				if row > 0 && col > 0 && row%2 == 1 { // Only toggle on content rows
					habitIndex := (row - 1) / 2
					day := ht.Days[col-1]
					ht.ToggleHabit(habitIndex, day)
					ht.UpdateUI()
				}
				return nil
			case 'q': // Quit
				ht.App.Stop()
				return nil
			case 'a': // Add new habit
				ht.ShowAddHabitDialog()
				return nil
			case 'd': // Delete a habit
				ht.ShowRemoveHabitDialog()
				return nil
			case '?': // Show help
				ht.ShowHelpDialog()
				return nil
			}
		}
		return event
	})
}

// ShowHelpDialog displays a help dialog with keybindings
func (ht *HabitTracker) ShowHelpDialog() {
	modal := tview.NewModal()
	modal.SetText(`
Habit Tracker Keybindings:

h, j, k, l - Navigate (vim-style)
Space or x - Toggle habit status
a          - Add new habit
d          - Delete a habit
?          - Show this help dialog
q or Esc   - Quit application
`)
	modal.SetBackgroundColor(tcell.ColorDefault)
	modal.AddButtons([]string{"Close"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		ht.App.SetRoot(ht.Table, true)
	})
	ht.App.SetRoot(modal, true)
}

// UpdateUI updates the UI with current data
func (ht *HabitTracker) UpdateUI() {
	// Update habit rows
	for row, habit := range ht.Habits {
		for col, day := range ht.Days {
			cell := ht.Table.GetCell(row*2+1, col+1)
			if completed, exists := habit.DayResults[day]; exists && completed {
				cell.SetText("✓")
				cell.SetTextColor(getColorByName(habit.Color))
			} else {
				cell.SetText("✗")
				cell.SetTextColor(tcell.ColorGray)
			}
		}
	}
}

// Run starts the application
func (ht *HabitTracker) Run() error {
	// Load data
	if err := ht.LoadData(); err != nil {
		return err
	}

	// Build UI
	ht.BuildUI()

	// Create application layout with title
	mainFlex := tview.NewFlex()
	mainFlex.SetDirection(tview.FlexRow)
	mainFlex.SetBackgroundColor(tcell.ColorDefault) // Transparent background

	// Create ASCII art title with transparent background
	title := tview.NewTextView()
	title.SetTextAlign(tview.AlignCenter)
	title.SetBackgroundColor(tcell.ColorDefault) // Transparent background
	title.SetText(`
 _   _       _     _ _      _____               _             
| | | |     | |   (_) |    |_   _|             | |            
| |_| | __ _| |__  _| |_     | |_ __ __ _  ___| | _____ _ __ 
|  _  |/ _` + "`" + ` | '_ \| | __|    | | '__/ _` + "`" + ` |/ __| |/ / _ \ '__|
| | | | (_| | |_) | | |_     | | | | (_| | (__|   <  __/ |   
\_| |_/\__,_|_.__/|_|\__|    \_/_|  \__,_|\___|_|\_\___|_|   
`)
	title.SetTextColor(tcell.ColorGreen)

	// Add title and some padding
	mainFlex.AddItem(title, 7, 0, false)

	// Transparent padding box
	paddingBox := tview.NewBox()
	paddingBox.SetBackgroundColor(tcell.ColorDefault)
	mainFlex.AddItem(paddingBox, 1, 0, false) // Small padding

	// Create a horizontal flex to center the table
	tableFlex := tview.NewFlex()
	tableFlex.SetBackgroundColor(tcell.ColorDefault) // Transparent background
	tableFlex.AddItem(nil, 0, 1, false)
	tableFlex.AddItem(ht.Table, 0, 3, true)
	tableFlex.AddItem(nil, 0, 1, false)

	// Add centered table to the main layout
	mainFlex.AddItem(tableFlex, 0, 1, true)

	// Set the application root
	ht.App.SetRoot(mainFlex, true)
	ht.App.SetFocus(ht.Table)

	// Run the application
	return ht.App.Run()
}

// GetVersion returns the current version
func GetVersion() string {
	return Version
}
