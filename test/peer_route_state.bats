#!/usr/bin/env ./test/libs/bats/bin/bats
load 'common'

@test "verify count of peer routes" {
  run get_route_count
  assert_output '32'
}

@test "verify peer routes are withdrawn" {
  run withdraw_routes
  run get_peer_metrics
  sleep 5
  # we don't care how many, if one is being withdraw then all should but counter updates take time
  assert_line --regexp '^peer_route_state\{.*\} 0$'
}

@test "verify peer routes are announced" {
  run announce_routes
  run get_peer_metrics
  sleep 5
  # we don't care how many, if one is being withdraw then all should but counter updates take time
  assert_line --regexp '^peer_route_state\{.*\} 1$'
}
