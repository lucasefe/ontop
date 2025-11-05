# Implementation Plan: README Visual Enhancement with Screenshots

**Branch**: `003-readme-visual-enhancement` | **Date**: 2025-11-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-readme-visual-enhancement/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Enhance the README documentation with high-quality visual assets demonstrating OnTop's TUI and CLI capabilities. Create comprehensive documentation including screenshots of the kanban board, CLI command examples with actual output, and a feature showcase section. This is a documentation-only feature requiring sample data creation, screenshot capture, and README restructuring to provide immediate visual context for potential users evaluating the project.

## Technical Context

**Language/Version**: Go 1.25.3 (existing project)
**Primary Dependencies**: Markdown (README format), PNG image format, shell scripting for sample data generation
**Storage**: Git repository (`docs/images/` for screenshots), filesystem for markdown files
**Testing**: Manual validation of README rendering on GitHub, visual inspection of screenshots, accessibility check for alt text
**Target Platform**: GitHub web interface (markdown rendering), cross-platform documentation viewers
**Project Type**: Documentation (non-code feature)
**Performance Goals**: README loads within 3 seconds on GitHub despite multiple images
**Constraints**: Total image size under 2MB, PNG format for screenshots, descriptive alt text for accessibility
**Scale/Scope**: 3-5 screenshots, documentation for existing features, sample data script for 5+ parent tasks with 3+ subtasks

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Simplicity First ✅ PASS

**Check**: Does this feature add unnecessary complexity or abstractions?

**Assessment**: No complexity violations. This is a documentation-only feature with no code changes to the application. Creating sample data and capturing screenshots are straightforward manual tasks. No clever solutions or premature optimization involved.

**Rationale**: Enhancing README with visual examples is a simple, maintainable approach to improving project documentation. No abstractions or complex systems introduced.

### II. Vim-Style Interaction ✅ PASS

**Check**: Does this feature maintain vim-style keyboard-driven workflows?

**Assessment**: Not applicable - this is documentation only. The screenshots will demonstrate the existing vim-style navigation (j/k/h/l keys), reinforcing this principle visually for new users.

**Rationale**: Feature actually promotes the vim-style interaction principle by making it visible to potential users through screenshots.

### III. Dual-Mode Architecture ✅ PASS

**Check**: Does this feature maintain consistency between TUI and CLI modes?

**Assessment**: Fully aligned. The feature explicitly documents both modes with visual examples: TUI screenshots showing the kanban interface, and CLI examples showing command output. This demonstrates the dual-mode architecture to users.

**Rationale**: Feature strengthens understanding of dual-mode architecture by providing side-by-side examples of both interaction models.

### IV. Documentation Currency ✅ PASS

**Check**: Will documentation remain synchronized with feature implementations?

**Assessment**: This feature directly implements the Documentation Currency principle by enhancing the README to accurately reflect current capabilities. The spec requires reproducible sample data (FR-014) and documents screenshot regeneration process for future UI changes.

**Rationale**: Feature is the direct application of this constitutional principle. Sample data reproducibility ensures documentation can be updated when features change.

### Constitution Check Result: ✅ ALL GATES PASSED

No violations. No complexity tracking needed. Feature is fully aligned with all constitutional principles and actively demonstrates dual-mode architecture and documentation currency principles.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Documentation and Assets Structure
docs/
└── images/
    ├── tui-kanban-board.png         # Main kanban screenshot
    ├── tui-subtasks-hierarchy.png   # Hierarchical subtask display
    ├── tui-task-form.png            # Task creation/edit form
    ├── cli-list-output.png          # CLI list command output
    └── cli-show-output.png          # CLI show command output

# Sample Data (for screenshot generation)
scripts/
└── generate-sample-data.sh          # Script to populate sample tasks

# Main documentation file
README.md                             # Enhanced with screenshots and examples

# No changes to application code
internal/                             # Existing - no modifications
cmd/                                  # Existing - no modifications
```

**Structure Decision**: This is a documentation-only feature. No application code changes required. All work involves:
1. Creating `docs/images/` directory for screenshot storage
2. Creating sample data generation script in `scripts/` directory
3. Enhancing existing `README.md` with visual content and restructured sections

The existing OnTop application structure (`internal/`, `cmd/`, etc.) remains untouched. This feature only adds documentation assets and improves existing documentation.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

N/A - All constitution checks passed with no violations.
