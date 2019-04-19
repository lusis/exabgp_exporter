#!/usr/bin/env ./test/libs/bats/bin/bats
load 'common'

@test "verify up - embedded" {
  run get_exabgp_metrics
  assert_line --regexp '^exabgp_up 1$'
}

@test "verify up - standalone" {
  run get_exabgp_metrics 9570
  assert_line --regexp '^exabgp_up 1$'
}