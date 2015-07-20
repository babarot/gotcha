#!/bin/bash

export PLATFORM

user=b4b4r07
repo=gotcha

ink() {
    if [ "$#" -eq 0 -o "$#" -gt 2 ]; then
        echo "Usage: ink <color> <text>"
        echo "Colors:"
        echo "  black, white, red, green, yellow, blue, purple, cyan, gray"
        return 1
    fi

    local open="\033["
    local close="${open}0m"
    local black="0;30m"
    local red="1;31m"
    local green="1;32m"
    local yellow="1;33m"
    local blue="1;34m"
    local purple="1;35m"
    local cyan="1;36m"
    local gray="0;37m"
    local white="$close"

    local text="$1"
    local color="$close"

    if [ "$#" -eq 2 ]; then
        text="$2"
        case "$1" in
            black | red | green | yellow | blue | purple | cyan | gray | white)
                eval color="\$$1"
                ;;
        esac
    fi

    printf "${open}${color}${text}${close}"
}

logging() {
    if [ "$#" -eq 0 -o "$#" -gt 2 ]; then
        echo "Usage: ink <fmt> <msg>"
        echo "Formatting Options:"
        echo "  TITLE, ERROR, WARN, INFO, SUCCESS"
        return 1
    fi

    local color=
    local text="$2"

    case "$1" in
        TITLE)
            color=yellow
            ;;
        ERROR | WARN)
            color=red
            ;;
        INFO)
            color=green
            ;;
        SUCCESS)
            color=green
            ;;
        *)
            text="$1"
    esac

    timestamp() {
        ink gray "["
        ink purple "$(date +%H:%M:%S)"
        ink gray "] "
    }

    timestamp; ink "$color" "$text"; echo
}

ok() {
    logging SUCCESS "$1"
}

die() {
    logging ERROR "$1" 1>&2
}

lower() {
    if [ $# -eq 0 ]; then
        cat <&0
    elif [ $# -eq 1 ]; then
        if [ -f "$1" -a -r "$1" ]; then
            cat "$1"
        else
            echo "$1"
        fi
    else
        return 1
    fi | tr "[:upper:]" "[:lower:]"
}

ostype() {
    uname | lower
}

# os_detect export the PLATFORM variable as you see fit
os_detect() {
    export PLATFORM
    case "$(ostype)" in
        *'linux'*)  PLATFORM='linux'   ;;
        *'darwin'*) PLATFORM='darwin'  ;;
        *'bsd'*)    PLATFORM='bsd'     ;;
        *)          PLATFORM='unknown' ;;
    esac
}

main() {
    logging TITLE "== Bootstraping enhancd =="
    logging INFO "Installing dependencies..."
    sleep 1
    echo

    os_detect

    # equals to
    # but this one liner needs jq
    # curl --fail -X GET https://api.github.com/repos/b4b4r07/gomi/releases/latest | jq '.assets[0].browser_download_url' | xargs curl -L -O
    # http://stackoverflow.com/questions/24987542/is-there-a-link-to-github-for-downloading-a-file-in-the-latest-release-of-a-repo
    # http://stackoverflow.com/questions/18384873/how-to-list-the-releases-of-a-repository
    # http://stackoverflow.com/questions/5207269/releasing-a-build-artifact-on-github
    releases="$( curl -s -L https://github.com/"${user}"/"${repo}"/releases/latest |
    egrep -o '/'"${user}"'/'"${repo}"'/releases/download/[^"]*' |
    grep $PLATFORM )"

    # github releases not available
    if [ -z "$releases" ]; then
        die "URL that can be used as Github releases was not found"
        exit 1
    fi

    # download github releases for user's platform
    echo "$releases" | wget --base=http://github.com/ -i -

    # install repo
    re=$(uname -m | grep -o "..$")
    for i in $releases
    do
        bin="$(basename "$i" | grep "$re")"
        if [ -f "$bin" ]; then
            cp "$bin" "$repo"
            chmod 755 "$repo"
            logging INFO "installing to ${PATH%%:*}..."
            sudo install -m 0755 "$repo" "${PATH%%:*}"
            break
        fi
    done

    # remove the intermediate files
    for i in $releases
    do
        rm -f $(basename "$i")
    done

    # log
    if [ -x "${PATH%%:*}"/"$repo" ]; then
        ok "$repo: sucessfully installed"
        rm -f "$repo"
    else
        die "$repo: incomplete or unsuccessful installations"
        echo "please put the binary to somewhere you want"
        echo "(on UNIX-ly systems, /usr/local/bin or the like)"
        exit 1
    fi
}

if echo "$-" | grep -q "i"; then
    # -> source a.sh
    return

else
    # three patterns
    # -> cat a.sh | bash
    # -> bash -c "$(cat a.sh)"
    # -> bash a.sh
    if [ "$0" = "${BASH_SOURCE:-}" ]; then
        # -> bash a.sh
        exit
    fi

    if [ -n "${BASH_EXECUTION_STRING:-}" ] || [ -p /dev/stdin ]; then
        trap "die 'terminated'; exit 1" INT ERR
        # -> cat a.sh | bash
        # -> bash -c "$(cat a.sh)"
        main
    fi
fi
