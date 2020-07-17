package text

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
)

// line format:
// neighbor <string> local-ip <string> local-as <int> peer-as <int> router-id <string> family-allowed in-open <afi> <safi> <details>
var rxParseRIBLine = `^neighbor (?P<neighbor>\S+) local-ip (?P<local_ip>\S+) local-as (?P<local_as>\d+) peer-as (?P<peer_as>\d+) router-id (?P<router_id>\S+) family-allowed in-open (?P<afi>\S+) (?P<safi>\S+) (?P<details>.*)$`
var rxParseUnicast = `^(?P<nlri>\S+) next-hop (?P<next_hop>\S+)(| (?P<attributes>.*))$`

func parseUnicastLine(s string) (map[string]string, error) {
	md := make(map[string]string)
	re := regexp.MustCompile(rxParseUnicast)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return md, fmt.Errorf("unable to parse line")
	}
	keys := re.SubexpNames()
	if len(keys) != 0 {
		for i, name := range keys {
			if i != 0 {
				md[name] = matches[i]
			}
		}
	}
	return md, nil
}

func parseRIBLine(s string) (map[string]string, error) {
	md := make(map[string]string)
	re := regexp.MustCompile(rxParseRIBLine)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return md, fmt.Errorf("unable to parse line")
	}
	keys := re.SubexpNames()
	if len(keys) != 0 {
		for i, name := range keys {
			if i != 0 {
				md[name] = matches[i]
			}
		}
	}
	return md, nil
}

// RibFromBytes takes a byte slice and returns a collection of RIBMessage
func RibFromBytes(b []byte) ([]*RIBMessage, error) {
	var ribs []*RIBMessage
	reader := bufio.NewReader(bytes.NewReader(b))
	for {
		l, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		r, err := RibEntryFromString(string(l))
		if err != nil {
			return ribs, err
		}
		ribs = append(ribs, r)
	}

	return ribs, nil
}

// RibEntryFromString takes a text string and returns a RIBMessage
func RibEntryFromString(s string) (*RIBMessage, error) {
	nm := &RIBMessage{}
	md, err := parseRIBLine(s)
	if err != nil {
		return nil, err
	}

	nm.PeerIP = md["neighbor"]
	nm.LocalIP = md["local_ip"]
	nm.PeerAS = md["peer_as"]
	nm.LocalAS = md["local_as"]
	nm.AFI = md["afi"]
	nm.SAFI = md["safi"]
	nm.RouterID = md["router_id"]
	nm.Details = md["details"]
	return nm, nil
}

// RIBMessage represents the common elements in a text-based encoding exabgp message
type RIBMessage struct {
	PeerIP   string
	PeerAS   string
	LocalIP  string
	LocalAS  string
	AFI      string
	SAFI     string
	RouterID string
	Details  string
}

// Family returns the family of the rib entry
func (m *RIBMessage) Family() string {
	return m.AFI + " " + m.SAFI
}

// IPv4Unicast returns an ipv4 unicast from a rib line
func (m *RIBMessage) IPv4Unicast() (*IPv4UnicastAnnounceTextMessage, error) {
	if m.Family() != "ipv4 unicast" {
		return nil, fmt.Errorf("wrong entry family: %s", m.Family())
	}
	nm := &IPv4UnicastAnnounceTextMessage{}
	res, err := parseUnicastLine(m.Details)
	if err != nil {
		return nil, err
	}
	nm.NLRI = res["nlri"]
	nm.NextHop = res["next_hop"]
	nm.Attributes = res["attributes"]
	return nm, nil
}

// IPv4Flow returns an ipv4 flow from a rib line
func (m *RIBMessage) IPv4Flow() (*IPv4FlowAnnounceTextMessage, error) {
	return nil, nil
}

// IPv6Unicast returns an ipv6 unicast from a rib line
func (m *RIBMessage) IPv6Unicast() (*IPv6UnicastAnnounceTextMessage, error) {
	if m.Family() != "ipv6 unicast" {
		return nil, fmt.Errorf("wrong entry family: %s", m.Family())
	}
	nm := &IPv6UnicastAnnounceTextMessage{}
	res, err := parseUnicastLine(m.Details)
	if err != nil {
		return nil, err
	}
	nm.NLRI = res["nlri"]
	nm.NextHop = res["next_hop"]
	nm.Attributes = res["attributes"]
	return nm, nil
}

// IPv6Flow returns an ipv6 flow from a rib line
func (m *RIBMessage) IPv6Flow() (*IPv6FlowAnnounceTextMessage, error) {
	return nil, nil
}

// IPv4UnicastAnnounceTextMessage represents an ipv4-unicast announce in a text-based encoded exabgp message
type IPv4UnicastAnnounceTextMessage struct {
	NLRI       string
	NextHop    string
	Attributes string
}

// IPv4MplsVPNAnnounceTextMessage represents an ipv4-mpls-vpn announce in a text-based encoded exabgp message
type IPv4MplsVPNAnnounceTextMessage struct {
	NLRI               string
	Label              int
	NextHop            string
	RouteDistinguisher string
	Community          string
	Origin             string
	ASPath             []int
	ExtendedCommunity  string
	LocalPreference    int
	OriginatorID       string
}

// IPv4FlowAnnounceTextMessage represents an ipv4-flow announce in a text-based encoded exabgp message
type IPv4FlowAnnounceTextMessage struct {
	DestinationIPv4   string
	SourceIPv4        string
	Protocol          string
	SourcePort        string
	DestinationPort   string
	ExtendedCommunity string
}

// IPv6UnicastAnnounceTextMessage represents an ipv6-unicast announce in a text-based encoded exabgp message
type IPv6UnicastAnnounceTextMessage struct {
	NLRI       string
	NextHop    string
	Attributes string
}

// IPv6MplsVPNAnnounceTextMessage represents an ipv6-mpls-vpn announce in a text-based encoded exabgp message
type IPv6MplsVPNAnnounceTextMessage struct {
	NLRI               string
	Label              int
	NextHop            string
	RouteDistinguisher string
	Community          string
	Origin             string
	ASPath             []int
	ExtendedCommunity  string
	LocalPreference    int
	OriginatorID       string
}

// IPv6FlowAnnounceTextMessage represents an ipv6-flow announce in a text-based encoded exabgp message
type IPv6FlowAnnounceTextMessage struct {
	DestinationIPv6   string
	SourceIPv6        string
	Protocol          string
	SourcePort        string
	DestinationPort   string
	ExtendedCommunity string
}
