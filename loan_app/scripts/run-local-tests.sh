#!/bin/bash

function print_success {
    printf "\n"
    printf '\e[1;32m%-6s\e[m\n' "*********************"
    printf '\e[1;32m%-6s\e[m\n' "* All tests passed. *"
    printf '\e[1;32m%-6s\e[m\n' "*********************"
}

function print_failure {
    printf "\n"
    printf '\e[1;31m%-6s\e[m\n' "*****************************"
    printf '\e[1;31m%-6s\e[m\n' "* One or more tests failed. *"
    printf '\e[1;31m%-6s\e[m\n' "*****************************"
}

eval $(cat test.env) go test -cover -race -p=1 -count=1 -v $(go list ./... | grep -v /vendor/)
if [ $? == 0 ]; then
    print_success
else
    print_failure
    exit 1
fi

golangci-lint -D errcheck run
if [[ $? != 0 ]]; then
  echo "Lint failed. Exiting."
  exit 1
fi

go vet ./...
if [[ $? != 0 ]]; then
  echo "Vet failed. Exiting."
  exit 1
fi

echo "Good to go to commit!"

# Why printf
# https://stackoverflow.com/a/8467449

# How we are setting environment variables:
# https://stackoverflow.com/questions/19331497/set-environment-variables-from-file