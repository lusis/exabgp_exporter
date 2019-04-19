#!/usr/bin/env ./test/libs/bats/bin/bats
load 'common'

@test "verify reported down from invalid standalone exporter" {
  run get_exabgp_metrics 9571
  assert_line --regexp '^exabgp_up 0$'
}