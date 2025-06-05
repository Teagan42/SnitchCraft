#!/bin/bash

go mod tidy
exec "$@"
