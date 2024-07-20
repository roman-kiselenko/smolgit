#!/usr/bin/env bats

teardown() {
  cat output.log
  cat cfg.yaml
  rm -rf cfg.yaml
  rm -rf output.log
}

@test "can show help" {
    run ./bin/smolgit --help
    [ "${lines[0]}" = "Usage of ./bin/smolgit:" ]
}

@test "can generate config" {
    run ./bin/smolgit config
    [ "${lines[0]}" = "log:" ]
    ./bin/smolgit config > cfg.yaml
    test -f cfg.yaml
}

@test "up and running" {
    ./bin/smolgit config | yq '.log.color = false | ... comments=""' > cfg.yaml
    ./bin/smolgit --config=./cfg.yaml > output.log &
    server_pid=$!
    sleep 3
    kill $server_pid
    cat output.log | grep '"initialize web server" addr=:3080'
    cat output.log | grep '"initialize ssh server" addr=:3081'
    cat output.log | grep '"start server" brand=smolgit address=:3080'
    cat output.log | grep '"starting SSH server" addr=:3081'
    cat output.log | grep 'msg="os signal received" signal=terminated'
}