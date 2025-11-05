<!--
SYNC IMPACT REPORT
==================
Version Change: 2.0.0 → 2.1.0
Rationale: MINOR version - Added dual-mode principle and vim-style interaction guidance

Modified Principles:
- Expanded: II. TUI-Native Design → now includes dual-mode architecture
- Added: III. Dual-Mode Architecture (TUI + CLI)

Added Sections:
- Mode interaction details (TUI for interactive, CLI for automation)
- Vim-style keyboard navigation guidance

Template Updates:
- ⚠️ .specify/templates/* - Should consider both modes in planning
- ⚠️ .claude/commands/* - May need adjustment for dual-mode development

Date: 2025-11-04
-->

# OnTop Constitution

## Core Principles

### I. Simplicity First

Keep the codebase simple and maintainable. Prefer straightforward solutions over clever ones.
Avoid premature optimization and unnecessary abstractions. If a feature adds complexity without
clear value for personal use, skip it.

**Rationale**: Personal projects should be maintainable by one person. Simple code is easier to
understand, debug, and extend when returning to the project after time away.

### II. Vim-Style Interaction

Embrace vim-style keyboard-driven workflows. Use familiar keybindings (j/k for navigation, etc.)
and modal interaction patterns. The TUI should feel natural to vim users with intuitive, efficient
keyboard shortcuts.

**Rationale**: Vim users expect certain interaction patterns. Leveraging muscle memory from vim
reduces cognitive load and makes the tool feel like a natural extension of the development
environment.

### III. Dual-Mode Architecture

The application MUST support two distinct modes of operation:

1. **TUI Mode** (Interactive): Full-screen terminal interface for browsing, viewing, and managing
   tasks interactively. Vim-style navigation and keybindings. Visual feedback and real-time
   updates.

2. **CLI Mode** (Unattended): Command-line interface for scripted/automated task management.
   Text-based input/output suitable for pipes, scripts, and automation. No interactive prompts.

Both modes operate on the same underlying data store and MUST maintain consistency.

**Rationale**: Interactive work requires rich visual feedback (TUI), while automation and scripting
require simple, scriptable commands (CLI). Supporting both modes makes the tool useful in different
contexts without forcing awkward compromises.

## Governance

### Amendment Procedure

This constitution can be amended at any time to reflect the evolving needs of the project. Update
the version number and date when making changes:
- **MAJOR**: Fundamental principle changes
- **MINOR**: New principles added
- **PATCH**: Clarifications or wording improvements

**Version**: 2.1.0 | **Ratified**: 2025-11-04 | **Last Amended**: 2025-11-04
