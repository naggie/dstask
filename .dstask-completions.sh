_dstask() {
    cur="${COMP_WORDS[COMP_CWORD]}"
    COMPREPLY=( $(dstask _completions $cur) )
}

complete -o nospace -F _dstask dstask
complete -o nospace -F _dstask task
