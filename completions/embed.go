package completions

//nolint
import _ "embed"

// Zsh completion script
//
//go:embed zsh.sh
var Zsh string

// Bash completion script
//
//go:embed bash.sh
var Bash string

// Fish completion script
//
//go:embed completions.fish
var Fish string

// PowerShell completion script
//
//go:embed powershell.ps1
var PowerShell string
