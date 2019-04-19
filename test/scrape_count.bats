#!/usr/bin/env ./test/libs/bats/bin/bats
load 'common'

@test "verify scrape total - embedded" {
  run get_exabgp_metrics
  assert_line --regexp '^exabgp_exporter_total_scrapes [0-9]+$'
}

@test "verify scrape total - standalone" {
  run get_exabgp_metrics 9570
  assert_line --regexp '^exabgp_exporter_total_scrapes [0-9]+$'
}

@test "verify scrape total - standalone invalid" {
  run get_exabgp_metrics 9571
  assert_line --regexp '^exabgp_exporter_total_scrapes [0-9]+$'
}