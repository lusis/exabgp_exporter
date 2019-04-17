load 'libs/bats-support/load'
load 'libs/bats-assert/load'
get_peer_metrics() {
	docker exec exabgp_exporter curl -s http://localhost:9569/metrics | grep peer
}
stop_gobgpd() {
	docker exec exabgp_exporter s6-svc -d /var/run/s6/services/gobgp
}
start_gobgpd() {
	docker exec exabgp_exporter s6-svc -u /var/run/s6/services/gobgp
}
withdraw_routes() {
	docker exec exabgp_exporter /root/withdraw.sh
}
announce_routes() {
	docker exec exabgp_exporter /root/announce.sh
}
get_route_count() {
	get_peer_metrics | grep -c peer_route_state\{
}
