#!/usr/bin/env bash

set -e

tmux new-session -d -s services
tab=0

function new_tab() {
    name="$1"
    path="$2"
    tab=$(($tab + 1))
    tmux new-window -t services:"$tab" -n "$name"
    tmux send-keys -t services:"$tab" "cd $path; make run" enter
}

for ms in $(ls -d cmd/*); do
    name=$(basename "$ms")
    path="$ms"
    new_tab "$name" "$path"
done

tmux select-window -t services:1

tmux attach -t services
