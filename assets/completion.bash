_mfa_get_accounts() {
    mfa list 2>/dev/null | sed -n 's/^- *//p'
}

_mfa_complete() {
    local cur prev words cword
    _init_completion || return

    local subcommands="add gen rename qr list del import export"

    if [ $cword -eq 1 ]; then
        COMPREPLY=( $(compgen -W "${subcommands}" -- "$cur") )
        return 0
    fi

    if [ $cword -eq 2 ]; then
        case "$prev" in
            gen|rename|qr|del|rm|mv)
                local IFS=$'\n'
                COMPREPLY=( $(compgen -W "$(_mfa_get_accounts)" -- "$cur") )
                return 0
                ;;
            import|export)
                _filedir
                return 0
                ;;
        esac
    fi
}

complete -F _mfa_complete mfa
