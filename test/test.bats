#!/usr/bin/env bats

@test "can show help" {
    run ./bin/smolgit --help
    [ "${lines[0]}" = "Usage of ./bin/smolgit:" ]
}

@test "can generate config" {
    export config=./config.yaml
    run ./bin/smolgit config
    [ "${lines[0]}" = "log:" ]
    run ./bin/smolgit config > "$config"
    test -f "$config"
}