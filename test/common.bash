load 'libs/bats-support/load'
load 'libs/bats-assert/load'


get_exabgp_metrics() {
	local port=${1:-9569}
	docker exec exabgp_exporter curl -s http://localhost:${port}/metrics | grep exabgp
}

get_peer_metrics() {
	local port=${1:-9569}
	docker exec exabgp_exporter curl -s http://localhost:${port}/metrics | grep peer
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
	local port=${1:-9569}
	get_peer_metrics ${port}| grep -c exabgp_state_route\{
}
