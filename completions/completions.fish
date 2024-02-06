#!/usr/bin/env fish

complete -f -c dstask -a (echo (dstask _completions) | string collect)
#complete -f -c task -a (echo (task _completions) | string collect)
#complete -f -c t -a (echo (t _completions) | string collect)

