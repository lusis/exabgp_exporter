#!/usr/bin/env ./test/libs/bats/bin/bats
load 'common'

@test "verify peer_resets are captured" {
  run get_peer_metrics
  assert_line --regexp '^peer_resets\{.*\} [0-9]+$'
}
