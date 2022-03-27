package completions

//nolint
import _ "embed"

// Zsh completion script
//go:embed zsh.sh
var Zsh string

// Bash completion script
//go:embed bash.sh
var Bash string
