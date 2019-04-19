#!/usr/bin/env ./test/libs/bats/bin/bats
load 'common'

@test "verify no parse failures - embedded" {
  run get_exabgp_metrics
  assert_line --regexp '^exabgp_exporter_parse_failures 0$'
}

@test "verify no parse failures - standalone" {
  run get_exabgp_metrics 9570
  assert_line --regexp '^exabgp_exporter_parse_failures 0$'
}