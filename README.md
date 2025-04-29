# Habit Tracker CLI

A stylish, terminal-based habit tracker with vim-like keybindings, built in Go.

![Habit Tracker CLI](https://you-can-add-a-screenshot-url-here.png)

## Features

- Track habits with colored checkmarks
- Weekly view (Monday to Sunday)
- Automatic weekly reset
- Vim-style navigation (hjkl)
- Persistent storage of habit data per week
- Color-coded checkmarks based on habit type
- Add and remove habits dynamically
- Simple and clean terminal interface

## Installation

### Using Go

```bash
go install github.com/yourusername/habit-tracker@latest
```

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/habit-tracker.git
cd habit-tracker

# Build and install
make install
```

## Usage

Run the habit tracker:

```bash
habit-tracker
```

### Keyboard Controls

- `h` - Move left
- `j` - Move down
- `k` - Move up
- `l` - Move right
- `Space` or `x` - Toggle habit status (checkmark/cross)
- `a` - Add a new habit
- `d` - Delete a habit
- `?` - Show help screen
- `q` or `Esc` - Quit the application

## Development

### Requirements

- Go 1.16 or higher
- [tcell](https://github.com/gdamore/tcell) for terminal handling
- [tview](https://github.com/rivo/tview) for UI components

### Build from Source

```bash
make build
```

The binary will be placed in the `build` directory.

## Data Storage

The application stores habit data in a JSON file in your home directory at:
`~/.habit-tracker/habits_YEAR_WEEK.json`

Data automatically resets each week, creating a new file for the new week.

## Default Habits

The application comes preconfigured with these habits:

1. water (blue)
2. exercise (red)
3. certification (yellow)
4. breath (white)
5. newsboat (orange)
6. recap (lightblue)
7. personal (green)
8. read (purple)

Each habit has a unique color for its checkmark when completed.

## License

MIT

## Credits

Built with:

- [tcell](https://github.com/gdamore/tcell)
- [tview](https://github.com/rivo/tview)
