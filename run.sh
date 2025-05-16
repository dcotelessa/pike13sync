#!/bin/bash
export $(grep -v '^#' .env | xargs)
go run cmd/pike13sync/main.go "$@"
