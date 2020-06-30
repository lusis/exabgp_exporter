#!/usr/bin/env ./test/libs/bats/bin/bats
load 'common'

@test "verify peer routes announce - embedded" {
  run announce_routes
  sleep 5
  run get_peer_metrics
  # we don't care how many, if one is being withdraw then all should but counter updates take time
  assert_line --regexp '^exabgp_state_route\{.*\} 1$'
}

@test "verify peer routes announce - standalone" {
  run announce_routes
  sleep 5
  run get_peer_metrics 9570
  assert_line --regexp '^exabgp_state_route\{.*\} 1$'
}

@test "verify peer routes ipv4 announce - embedded" {
  run announce_routes
  sleep 5
  run get_peer_metrics
  assert_line --regexp '^exabgp_state_route\{.*family="ipv4 unicast".+nlri="192\.168\.0\.0/24".*\} 1$'
}

@test "verify peer routes ipv6 announce - embedded" {
  run announce_routes
  sleep 5
  run get_peer_metrics
  assert_line --regexp '^exabgp_state_route\{.*family="ipv6 unicast".+nlri="2001:db8:1000::/64".*\} 1$'
}

@test "verify peer routes ipv4 announce - standalone" {
  run announce_routes
  sleep 5
  run get_peer_metrics 9570
  assert_line --regexp '^exabgp_state_route\{.*family="ipv4 unicast".+nlri="192\.168\.0\.0/24".*\} 1$'
}

@test "verify peer routes ipv6 announce - standalone" {
  run announce_routes
  sleep 5
  run get_peer_metrics 9570
  assert_line --regexp '^exabgp_state_route\{.*family="ipv6 unicast".+nlri="2001:db8:1000::/64".*\} 1$'
}

@test "verify count of peer routes - embedded" {
  run announce_routes
  sleep 5
  run get_route_count
  assert_output '34'
}

@test "verify count of peer routes - standalone" {
  run announce_routes
  sleep 5
  run get_route_count 9570
  assert_output '34'
}

@test "verify peer routes withdraw - embedded" {
  run withdraw_routes
  sleep 5
  run get_peer_metrics
  # we don't care how many, if one is being withdraw then all should but counter updates take time
  assert_line --regexp '^exabgp_state_route\{.*\} 0$'
}

@test "verify peer routes withdraw - standalone" {
  run withdraw_routes
  sleep 5
  run get_peer_metrics 9570
  # standalone exporter should not have any results
  refute_line --regexp '^exabgp_state_route\{.*\} 0$'
}
