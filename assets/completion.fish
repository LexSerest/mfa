complete -c mfa -f

function __mfa_get_accounts
    mfa list 2>/dev/null | string match -r '^- \s*(.*)' | string replace -r '^- \s*' ''
end

function __mfa_is_second_arg
    set -l tokens (commandline -co)
    if test (count $tokens) -eq 2
        return 0
    end
    return 1
end


complete -c mfa -n "not __fish_seen_subcommand_from add gen rename qr list del import export" \
    -a "add gen rename qr list del import export"

complete -c mfa -n "__fish_seen_subcommand_from gen rename qr del rm mv; and __mfa_is_second_arg" \
    -a "(__mfa_get_accounts)"

complete -c mfa -n "__fish_seen_subcommand_from import export; and __mfa_is_second_arg" -F
