#!/usr/bin/env ./test/libs/bats/bin/bats
load 'common'

@test "verify no parse failures" {
  run docker logs exabgp_exporter
  refute_line --regexp '^.*unable to parse line:.*$'
}
