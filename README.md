# OnTop - Dual-Mode Task Manager

A modern, flexible task management tool built in Go that works both as an interactive TUI (Terminal User Interface) and a traditional CLI.

## Features

- **Dual-Mode Interface**: Use as an interactive Kanban board TUI or run commands from the terminal
- **Kanban Workflow**: Organize tasks across columns (Todo, In Progress, Done)
- **SQLite Storage**: Local database storage with automatic migrations
- **Rich CLI Commands**: Add, list, show, move, and update tasks from the command line
- **Terminal UI**: Beautiful, interactive Kanban board powered by Bubble Tea

## Installation

### Build from source

```bash
# Clone the repository
git clone https://github.com/lucasefe/ontop.git
cd ontop

# Build the binary
make build

# Or install to $GOPATH/bin
make install
```

### Requirements

- Go 1.25.3 or later

## Usage

### Interactive TUI Mode

Launch the interactive Kanban board by running `ontop` without any arguments:

```bash
./ontop
```

In TUI mode, you can:
- Navigate tasks with arrow keys
- View task details
- Move tasks between columns
- Interactive task management

### CLI Mode

Use traditional command-line commands for task management:

```bash
# Add a new task
./ontop add "Implement new feature" --priority high --description "Add user authentication"

# List all tasks
./ontop list

# Show task details
./ontop show <task-id>

# Move task to a different column
./ontop move <task-id> "In Progress"

# Update task attributes
./ontop update <task-id> --title "New title" --priority medium
```

### Available Commands

- `add` - Create a new task
- `list`, `ls` - List all tasks
- `show` - Show detailed information about a task
- `move`, `mv` - Move a task to a different column
- `update`, `edit` - Update task attributes
- `help` - Show help message

### Global Options

- `--db-path` - Specify custom database path (default: `~/.ontop/tasks.db`)

Example:
```bash
./ontop --db-path /path/to/custom.db list
```

## Development

### Project Structure

```
cmd/ontop/          # Main application entry point
internal/
  cli/              # CLI command implementations
  models/           # Data models (Task, Column)
  service/          # Business logic and services
  storage/          # Database layer (SQLite)
  tui/              # Terminal UI components
tests/              # Test files
```

### Building and Testing

```bash
# Build the binary
make build

# Run tests
make test

# Run directly (TUI mode)
make run

# Clean build artifacts and database
make clean

# Show all available make targets
make help
```

### Technologies

- **Go 1.25.3** - Primary language
- **Bubble Tea** - Terminal UI framework
- **SQLite** - Local database
- **Lipgloss** - Terminal styling

## Database

OnTop uses SQLite for local storage. By default, the database is stored at `~/.ontop/tasks.db`.

The schema is automatically created and migrated on first run.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

[Add your license here]

## Author

Lucas Efe
