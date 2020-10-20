#compdef dstask
#autoload


_dstask() {
    compadd -- $(dstask _completions "${words[@]}")
}

compdef _dstask dstask
