#!/usr/bin/env bats

@test "can show help" {
    run ./bin/smolgit --help
    [ "${lines[0]}" = "Usage of ./bin/smolgit:" ]
}

@test "can generate config" {
    run ./bin/smolgit config
    [ "${lines[0]}" = "log:" ]
}