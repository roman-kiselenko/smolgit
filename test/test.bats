#!/usr/bin/env bats

teardown() {
  cat output.log
  cat cfg.yaml
  rm -rf cfg.yaml
  rm -rf output.log
}

@test "smolgit can show help" {
    run ./bin/smolgit --help
    [ "${lines[0]}" = "Usage of ./bin/smolgit:" ]
}

@test "smolgit can generate config" {
    run ./bin/smolgit config
    [ "${lines[0]}" = "log:" ]
    ./bin/smolgit config > cfg.yaml
    test -f cfg.yaml
}

@test "smolgit up and running" {
    ./bin/smolgit config | yq '.log.color = false | ... comments=""' > cfg.yaml
    ./bin/smolgit --config=./cfg.yaml > output.log &
    server_pid=$!
    sleep 1
    kill $server_pid
    cat output.log | grep '"initialize web server" addr=:3080'
    cat output.log | grep '"initialize ssh server" addr=:3081'
    cat output.log | grep '"start server" brand=smolgit address=:3080'
    cat output.log | grep '"starting SSH server" addr=:3081'
    cat output.log | grep 'msg="os signal received" signal=terminated'
}

@test "web server respond 404 no repositories found" {
    ./bin/smolgit config | yq '.log.color = false | ... comments=""' > cfg.yaml
    ./bin/smolgit --config=./cfg.yaml > output.log &
    server_pid=$!
    sleep 1
    response_code=$(curl -so /dev/null -w '%{response_code}' localhost:3080/)
    echo "$response_code"
    [ "$response_code" -ne "200" ]
    kill $server_pid
    cat output.log | grep 'msg="hit route" route=/'
}