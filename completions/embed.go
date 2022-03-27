package completions

//nolint
import _ "embed"

//go:embed zsh.sh
var Zsh string

//go:embed bash.sh
var Bash string
