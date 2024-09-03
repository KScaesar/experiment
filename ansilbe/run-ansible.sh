#!/bin/bash

set -e

# init
# 以一個新的 vm 開始, 進行 deploy 的事前準備工作, 只會執行一次
#
# deploy
# 部署服務到 remote vm, 會執行多次

show_help() {
    echo "Usage: $0 [init|deploy] --env <env> --no <build number>"
    echo "  init   : Run the init.yml playbook"
    echo "  deploy : Run the deploy.yml playbook"
    echo "  --env  : Specify the environment (e.g., staging, production)"
    echo "  --no   : Specify the Jenkins build number"
    echo
    echo "Example:"
    echo "  $0 deploy --env production --no 123"
    echo "  $0 init --env staging"
    exit 1
}

require_value() {
    local key=$1
    if [ -z "$2" ]; then
        echo "Error: '$key' requires a value" >&2
        show_help
    fi
}

info() {
    echo -e "\033[32m[Info] $1\033[0m"
}

COMMAND=$1
if [ -z "$COMMAND" ] || { [ "$COMMAND" != "init" ] && [ "$COMMAND" != "deploy" ]; }; then
    show_help
fi
shift 1

ENV=""
BUILD_NO=""
while [[ "$#" -gt 0 ]]; do
    require_value "$1" "$2"

    case $1 in
        --env)
            ENV="$2"
            info "ENV=$ENV"
            shift 2
            ;;
        --no)
            BUILD_NO="$2"
            info "BUILD_NO=$BUILD_NO"
            shift 2
            ;;
        *)
            echo "Invalid option: $1" >&2
            show_help
            ;;
    esac
done


INVENTORY_FILE="hosts.${ENV}"
if [ ! -f "$INVENTORY_FILE" ]; then
    echo "Error: Inventory file $INVENTORY_FILE does not exist"
    show_help
fi

case $COMMAND in
    init)
        ansible-playbook -i "$INVENTORY_FILE" init.yml -vv
        ;;
    deploy)
        ansible-playbook -i "$INVENTORY_FILE" --extra-vars "build_number=$BUILD_NO env=$ENV" deploy.yml -vv
        ;;
esac

