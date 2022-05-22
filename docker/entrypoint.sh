#!/bin/sh
set -e

if [ "$GITLAB_CI" == "true" ]; then
    # Gitlab CI environment detected, skipping entrypoint
    exit 0
fi

if [ "$1" == "semantic-release" ]; then
    # remove first argument
    shift 1
fi

exec semantic-release "$@"
