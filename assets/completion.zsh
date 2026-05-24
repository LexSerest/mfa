#compdef mfa

_mfa_get_accounts() {
    mfa list 2>/dev/null | sed -n 's/^- *//p'
}

_mfa() {
    local context state line
    typeset -A opt_args

    local -a subcommands
    subcommands=(
        'add:Add a new account'
        'gen:Generate a code'
        'rename:Rename an account'
        'qr:Show QR code'
        'list:List all accounts'
        'del:Delete an account'
        'import:Import data'
        'export:Export data'
    )

    _arguments -C \
        '1: :->command' \
        '*:: :->args'

    case $state in
        command)
            _describe 'command' subcommands
            ;;
        args)
            case $words[1] in
                gen|rename|qr|del|rm|mv)
                    if (( CURRENT == 2 )); then
                        local -a accounts
                        accounts=(${(f)"$(_mfa_get_accounts)"})
                        _describe 'account' accounts
                    fi
                    ;;
                import|export)
                    if (( CURRENT == 2 )); then
                        _files
                    fi
                    ;;
            esac
            ;;
    esac
}

_mfa "$@"
