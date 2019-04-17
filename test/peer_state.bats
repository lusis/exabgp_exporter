#!/usr/bin/env ./test/libs/bats/bin/bats
load 'common'

@test "verify peer_state is captured" {
  run get_peer_metrics
  assert_line --regexp '^peer_state\{.*\} [0-9]$'
}

@test "verify peer_state is down" {
  run stop_gobgpd
  sleep 2
  run get_peer_metrics
  assert_line --regexp '^peer_state\{.*\} 0$'
}

@test "verify peer_resets are captured" {
  run get_peer_metrics
  assert_line --regexp '^peer_resets\{.*\} [0-9]+$'
}

@test "verify peer_state is up" {
  run start_gobgpd
  sleep 60
  run get_peer_metrics
  assert_line --regexp '^peer_state\{.*\} 1$'
}
