# Register PowerShell argument completer for dstask using the built-in completion engine.
# It invokes `dstask _completions` with the current command line to get suggestions.

Register-ArgumentCompleter -Native -CommandName dstask,task -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)

    try {
        # Build args: dstask _completions <user command line>
        # We collect only argument tokens (ignore command name itself)
        $tokens = [System.Management.Automation.PSParser]::Tokenize($commandAst.Extent.Text, [ref]$null)
        $argTokens = $tokens | Where-Object { $_.Type -eq 'CommandArgument' } | ForEach-Object { $_.Content }
        $args = @('_completions') + $argTokens

        $completions = & dstask @args 2>$null
        if (-not $completions) { return }

        foreach ($c in $completions) {
            if ($c -like "$wordToComplete*") {
                [System.Management.Automation.CompletionResult]::new($c, $c, 'ParameterValue', $c)
            }
        }
    } catch {
        # no-op on errors to avoid noisy completion failures
    }
}


