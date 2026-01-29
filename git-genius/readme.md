# Git Genius

Git Genius is a beginner-friendly, interactive Git CLI tool that helps you work with Git repositories **without memorizing commands**.

It is designed to be:
- safe for beginners
- useful for daily development
- modular and future-ready

---

## Why Git Genius?

Git is powerful, but the command-line workflow can be confusing, especially for new developers.

Git Genius provides:
- a menu-driven interface
- guided setup and validation
- smart handling of common Git problems
- clear messages instead of cryptic errors

You focus on your code. Git Genius manages the Git workflow.

---

## Core Features

### Project & Repository Management
- Select any project directory
- Initialize Git if repository does not exist
- Work with multiple projects easily
- Safe confirmation before destructive actions

### Daily Git Operations
- Git status
- Push changes with commit message
- Pull latest changes
- Fetch all remotes
- Switch branch
- Switch remote

### Smart Workflow Features
- Smart Pull (auto-stash → pull → restore changes)
- Stash manager
  - stash save
  - stash list
  - stash pop
- Undo last commit safely (changes preserved)

### Guided Setup
- Step-by-step setup wizard
- Choose project directory
- Configure default branch and remote
- GitHub username or organization support
- GitHub token help with validation
- Automatic remote configuration

### Doctor (Health Check)
- Git installation check
- Project directory validation
- Git repository detection
- Git user.name and user.email check
- Internet connectivity check
- GitHub token validation
- Error log detection with guidance

---

## Design Philosophy

- Menu-based, no command memorization required
- Safe defaults, explicit confirmations
- Modular architecture
- Read-only diagnostics (Doctor never changes data)
- Clean separation of concerns

---

## Future-Upgradable Features (Planned)

- Commit history viewer
- Amend last commit message
- Diff viewer (file-wise changes)
- GitHub repository creation via API
- Recent projects list
- Project switcher
- Plugin-style command extensions
- CI/CD helper commands

These features are **not promises**, but the architecture is intentionally built to support them.

---

## Disclaimer

Git Genius is a helper tool, not a replacement for Git knowledge.

- It wraps Git commands; it does not modify Git itself
- You are still responsible for your repositories
- Always review actions before confirming

Use responsibly, especially on important repositories.

---

## License

Open-source.  
You are free to learn, modify, and extend.

---

## Final Note

Git Genius aims to make Git **approachable**, not magical.

If you understand Git better after using this tool,  
it has done its job.