<!--
SYNC IMPACT REPORT
==================
Version Change: 2.1.0 → 2.2.0
Rationale: MINOR version - Added Documentation Currency principle (new principle added)

Modified Principles:
- None

Added Sections:
- IV. Documentation Currency - Requirement to keep README synchronized with features

Template Updates:
- ✅ .specify/templates/plan-template.md - Constitution Check section already flexible
- ✅ .specify/templates/spec-template.md - No changes needed (requirements-focused)
- ✅ .specify/templates/tasks-template.md - Polish phase includes documentation updates
- ⚠️ README.md - NEEDS UPDATE to include subtask features from spec 002

Follow-up TODOs:
- Update README.md to document subtask hierarchy features (in progress)
- Consider adding documentation validation to CI/CD pipeline in future

Date: 2025-11-05
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

### IV. Documentation Currency

The README MUST remain synchronized with feature implementations and project capabilities. When new
features are added or existing features are modified, the documentation MUST be updated in the same
development cycle to reflect the current state of the application.

**Rationale**: Stale documentation misleads users and degrades trust in the project. For personal
projects maintained intermittently, accurate documentation serves as memory aid when returning
after time away. Documentation drift creates confusion about what the tool actually does versus
what it's documented to do.

**Requirements**:
- README MUST document all user-facing features (TUI and CLI)
- Feature additions MUST include corresponding README updates
- Breaking changes MUST update relevant documentation sections
- Examples and usage patterns MUST reflect actual current behavior

## Governance

### Amendment Procedure

This constitution can be amended at any time to reflect the evolving needs of the project. Update
the version number and date when making changes:
- **MAJOR**: Fundamental principle changes
- **MINOR**: New principles added
- **PATCH**: Clarifications or wording improvements

**Version**: 2.2.0 | **Ratified**: 2025-11-04 | **Last Amended**: 2025-11-05
