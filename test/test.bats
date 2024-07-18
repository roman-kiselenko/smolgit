#!/usr/bin/env bats

@test "can generate config" {
    run ./bin/smolgit config
    [ "${lines[0]}" = "log:" ]
}