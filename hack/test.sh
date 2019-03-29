#!/bin/bash

go test ./... | sed ''/ok/s//$(printf "\033[32mok\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
