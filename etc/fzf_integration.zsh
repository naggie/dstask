function dstask () {
    local has_double_dash subcommand task_ids

    if (( $# == 1 || $# == 2 )); then
        has_double_dash="$( [[ "$1" == "--" || "$2" == "--" ]] && echo "1" || echo "0" )"
        subcommand="$( [[ "$1" == "--" ]] && echo "$2" || echo "$1" )"
        if (( $# == 1 || has_double_dash )) && { \
            [[ "${subcommand}" =~ ^(done|edit|note|open|remove|start|stop)$ ]] || \
            [[ "${subcommand}" == "modify" && -n "${WIDGET}" ]] \
        }; then
            task_ids=("${(@f)$( \
                command dstask show-open $( (( has_double_dash )) && echo "--" ) \
                | jq -r '.[] | [.id, .summary] | @tsv' \
                | fzf --multi --delimiter='\t' --nth=2 --accept-nth=1 --preview='DSTASK_FAKE_PTY=1 dstask {1} --' \
            )}")
            (( ${#task_ids[@]} > 0 )) || return
            if [[ "${subcommand}" == "modify" ]]; then
                RBUFFER=""
                LBUFFER="dstask modify ${task_ids[@]}"
            else
                xargs -n1 dstask "${subcommand}" <<< "${task_ids[@]}"
            fi
            return
        fi
    fi

    command dstask "$@"
}

function dstask-accept-line () {
    if [[ "${BUFFER}" =~ "^\s*dstask\s+modify\s*$" ]]; then
        dstask modify
        zle redisplay
    else
        zle orig-accept-line
    fi
}
zle -A accept-line orig-accept-line
zle -N accept-line dstask-accept-line
