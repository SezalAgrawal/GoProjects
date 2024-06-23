#!/bin/bash

function loadenvfile {
    echo "Loading $1"
    for i in $(cat $1 | grep "^[^#;]"); do
        export $i
    done
}

loadenvfile development.env
go run cmd/loan_app/main.go