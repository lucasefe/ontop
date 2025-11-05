# Quickstart Guide: Row-Based UI Mode

**Feature**: 004-row-ui | **For**: End Users | **Date**: 2025-11-05

## What's New

OnTop now supports two ways to visualize your tasks in the TUI:

- **Column Mode** (default): Traditional vertical columns for inbox, in_progress, and done
- **Row Mode** (new): Horizontal rows showing each workflow stage stacked vertically

Toggle between modes instantly with a single keypress. Your preference is saved automatically.

## Quick Start

### Toggle View Mode

Press `v` (for "view") at any time in the Kanban board to switch layouts:

```
Column Mode:          →  Press 'v'  →      Row Mode:

┌─────────────┐                           ╔═══════════════════╗
│   INBOX     │                           ║ INBOX             ║
│             │                           ║ [◀] task1 task2 [▶]║
│  • task1    │                           ╚═══════════════════╝
│  • task2    │                           ╔═══════════════════╗
│  • task3    │                           ║ IN PROGRESS       ║
└─────────────┘                           ║ [◀] task4 task5 [▶]║
                                          ╚═══════════════════╝
┌─────────────┐                           ╔═══════════════════╗
│ IN PROGRESS │                           ║ DONE              ║
│             │                           ║ [◀] task6 task7 [▶]║
│  • task4    │                           ╚═══════════════════╝
│  • task5    │
└─────────────┘
```

Press `v` again to return to column mode.

## Navigation

Navigation adapts intelligently to your current layout:

### Column Mode (Existing)

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down to next task in current column |
| `k` / `↑` | Move up to previous task in current column |
| `h` / `←` | Switch to column on the left |
| `l` / `→` | Switch to column on the right |

### Row Mode (New)

| Key | Action |
|-----|--------|
| `j` / `↓` | Switch to row below (inbox → in_progress → done) |
| `k` / `↑` | Switch to row above (done → in_progress → inbox) |
| `h` / `←` | Move to previous task in current row |
| `l` / `→` | Move to next task in current row |

**Navigation Principle**: Arrow keys always move in the primary visual dimension:
- Column mode: Primary dimension is vertical → up/down navigate tasks
- Row mode: Primary dimension is horizontal → left/right navigate tasks

### Scrolling (Row Mode Only)

When a row has more tasks than fit on screen, scroll indicators appear:

```
╔═══════════════════════════════════╗
║ INBOX                             ║
║ [◀] task2 | task3 | task4 [▶]     ║  ← Indicators show more tasks
╚═══════════════════════════════════╝
    ↑                         ↑
  More left               More right
```

- `[◀]`: Tasks exist to the left (scroll left with `h` / `←`)
- `[▶]`: Tasks exist to the right (scroll right with `l` / `→`)
- Scroll offset is independent per row (each row remembers its position)
- Scroll resets when you toggle to column mode

## Task Operations

All task operations work identically in both modes:

| Key | Action | Both Modes |
|-----|--------|------------|
| `Enter` | View task details | ✓ |
| `n` | Create new task | ✓ |
| `e` | Edit selected task | ✓ |
| `m` | Move task to different stage | ✓ |
| `d` | Delete task (with confirmation) | ✓ |
| `a` | Toggle show archived tasks | ✓ |
| `s` | Cycle sort mode | ✓ |
| `v` | Toggle view mode (column/row) | ✓ |
| `?` | Show help | ✓ |
| `q` | Quit application | ✓ |

## Moving Tasks

### In Column Mode

1. Press `m` on selected task
2. Use `h`/`l` to choose destination column
3. Press `Enter` to confirm

### In Row Mode

1. Press `m` on selected task
2. Use `j`/`k` to choose destination row
3. Press `Enter` to confirm

The move UI adapts to match your current layout for consistency.

## Configuration

Your view mode preference is saved automatically to:

```
~/.config/ontop/ontop.toml
```

Example config file:
```toml
[ui]
view_mode = "row"
```

### How Persistence Works

1. **First launch**: App starts in column mode (default)
2. **Toggle to row mode**: Press `v` → preference saved immediately
3. **Next launch**: App starts in row mode (your saved preference)
4. **Toggle back**: Press `v` → preference updated to column mode

**Note**: If config file doesn't exist or can't be written, app continues working normally (preference just won't persist between sessions).

## Examples

### Scenario 1: Switching to Row Mode

You find the column layout doesn't work well in your terminal:

1. Launch OnTop (starts in column mode by default)
2. Press `v` to switch to row mode
3. Navigate with `j`/`k` between rows, `h`/`l` within rows
4. Your preference is saved automatically
5. Next launch: OnTop starts in row mode

### Scenario 2: Working with Many Tasks

You have 10 tasks in your inbox, but your terminal is 100 characters wide:

1. Switch to row mode with `v`
2. Navigate to inbox row (if not already there)
3. Use `l` / `→` to scroll right through tasks
4. Indicators `[◀]` and `[▶]` show there are more tasks
5. Select the task you want with `l` / `h`
6. Perform any operation (view, edit, move, etc.)

### Scenario 3: Comparing Layouts

You're curious which layout works better for your workflow:

1. Start in column mode (default)
2. Work with tasks using existing navigation
3. Press `v` to try row mode
4. Work with tasks using adapted navigation
5. Press `v` again to return to column mode
6. Toggle freely - your final choice is saved

## Tips

- **Try both modes**: Different terminal sizes and task counts may favor different layouts
- **Muscle memory**: vim keybindings work in both modes, just adapted to layout
- **Scroll indicators**: In row mode, watch for `[◀]` `[▶]` to know if more tasks exist
- **Preference persistence**: Toggle is immediate and persistent - no need to "save settings"
- **No data changes**: Switching layouts only affects visualization, never task data

## Troubleshooting

### Config File Issues

**Problem**: Preference doesn't persist between launches

**Solution**: Check that `~/.config/ontop/` directory is writable:
```bash
ls -ld ~/.config/ontop/
# Should show: drwxr-xr-x (user has write permission)
```

If directory doesn't exist:
```bash
mkdir -p ~/.config/ontop
```

### Terminal Width Issues

**Problem**: Row mode looks cramped or tasks are truncated

**Solution**: Row mode works best with terminal width ≥ 100 characters. Check your terminal width:
```bash
tput cols
```

If less than 80, consider:
- Increasing terminal window size
- Using column mode instead
- Adjusting terminal font size

### Navigation Confusion

**Problem**: Arrows don't do what I expect

**Solution**: Remember that navigation adapts to layout:
- Column mode: Arrows move in 2D grid (task-by-task and column-by-column)
- Row mode: Arrows move between rows and within rows
- Press `?` to see context-sensitive help for current layout

## Technical Details

### Config File Location

OnTop follows the XDG Base Directory specification:

- **Linux/macOS**: `~/.config/ontop/ontop.toml`
- **Windows**: `%APPDATA%\ontop\ontop.toml` (if XDG not set)

### Config Format

TOML format, human-readable and editable:

```toml
# OnTop Configuration

[ui]
view_mode = "column"  # Valid values: "column" or "row"
```

You can manually edit this file if needed. Invalid values default to "column".

### Performance

- **Toggle speed**: Instant (<1 second per requirement)
- **Config save**: ~1ms (non-blocking)
- **Rendering**: Same performance as column mode
- **Memory**: Negligible overhead (<500 bytes)

## Related Documentation

- **TUI Navigation**: See main README for comprehensive TUI keybindings
- **CLI Usage**: Row mode only affects TUI, CLI commands unchanged
- **Task Management**: Core task operations (create, edit, move) work identically in both modes

## Feedback

If you find issues or have suggestions for improving row mode:

1. Check that you're using the latest version
2. Verify config file is writable (see Troubleshooting above)
3. Report issues with details about terminal size and OS

## Summary

- **Toggle**: Press `v` anytime to switch between column and row layouts
- **Navigation**: Adapts intelligently - arrows move in primary visual dimension
- **Persistence**: Preference saved automatically, effective on next launch
- **Compatibility**: All task operations work in both modes
- **Flexibility**: Choose the layout that works best for your terminal and workflow
