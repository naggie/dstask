#compdef dstask
#autoload


_dstask() {
    compadd -- $(dstask _completions "${words[@]}")
}

compdef _dstask dstask
compdef _dstask p0
compdef _dstask p1
compdef _dstask p2
compdef _dstask p3
compdef _dstask ds
