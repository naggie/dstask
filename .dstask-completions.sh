_dstask() {
    cur="${COMP_WORDS[COMP_CWORD]}"
    COMPREPLY=( $(dstask _completions "${COMP_WORDS[@]}") )
}

complete -F _dstask dstask
complete -F _dstask task
complete -F _dstask n
complete -F _dstask t
