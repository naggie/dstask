#compdef pass
#autoload


_dstask() {
    compadd $(dstask _completions "${words[@]}")
}

compdef _dstask dstask
compdef task=dstask
