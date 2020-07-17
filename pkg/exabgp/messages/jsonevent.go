package messages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Timestamp represents the Exabgp timestamp
type Timestamp struct {
	time.Time
}

func timestampFromFloat64(ts float64) Timestamp {
	secs := int64(ts)
	nsecs := int64((ts - float64(secs)) * 1e9)
	return Timestamp{time.Unix(secs, nsecs)}
}

// UnmarshalJSON converts the float64 exabgp timestamp format to golang time.Time
func (e *Timestamp) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		e.Time = time.Time{}
		return nil
	}
	ts, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return fmt.Errorf("unable to parse timestamp: %s", string(data))
	}
	*e = timestampFromFloat64(ts)
	return nil
}

// JSONEvent represents a message as JSON
type JSONEvent struct {
	BaseEvent
	Neighbor `json:"neighbor"`
}

// BaseEvent is represents the common data in all messages
type BaseEvent struct {
	Version string    `json:"exabgp"`
	Time    Timestamp `json:"time"`
	Host    string    `json:"host"`
	PID     int       `json:"pid"`
	PPID    int       `json:"ppid"`
	Counter int64     `json:"counter"`
	Type    string    `json:"type"`
}

// Neighbor represents the common data in a neighbor message
type Neighbor struct {
	Address struct {
		Local string `json:"local"`
		Peer  string `json:"peer"`
	} `json:"address"`
	ASN struct {
		Local int `json:"local"`
		Peer  int `json:"peer"`
	} `json:"asn"`
	Direction    string      `json:"direction"`
	State        string      `json:"state"`
	Reason       string      `json:"reason"`
	Notification interface{} `json:"notification"`
	Message      struct {
		Update UpdateMessageFull `json:"update"`
	} `json:"message"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// NotificationMessage represents a notification message
type NotificationMessage struct {
	Code    int    `json:"code"`
	SubCode int    `json:"subcode"`
	Data    string `json:"data"`
}

// UpdateMessage is a bgp update message
type UpdateMessage struct {
	Attribute `json:"attribute"`
	Announce  json.RawMessage `json:"announce"` // "announce": { "ipv4 unicast": { "192.168.1.184": [ "192.168.88.2/32" ] } }
	Withdraw  json.RawMessage `json:"withdraw"` // "withdraw": { "ipv4 unicast": [ "192.168.88.2/32" ] }
}

// UpdateMessageFull represents an update message
type UpdateMessageFull struct {
	Attribute `json:"attribute"`
	Announce  struct {
		// messages are different between compact and non-compact mode
		// compact: { "ipv4 unicast": { "192.168.1.184": [ "192.168.88.2/32" ] } }
		// non-compact: "ipv4 unicast": [ { "nlri": "192.168.88.0/24" } ] } } } } }
		IPv4Unicast map[string][]interface{}      `json:"ipv4 unicast"`
		IPv4Flow    IPv4FlowAnnounceMessage       `json:"ipv4 flow"`
		IPv6Unicast map[string][]interface{}      `json:"ipv6 unicast"`
		IPv6Flow    IPv6FlowAnnounceMessage       `json:"ipv6 flow"`
		L2VPNVpls   map[string][]L2VPNVplsMessage `json:"l2vpn vpls"`
	} `json:"announce"`
	Withdraw struct {
		// compact: { "ipv4 unicast": { [ "192.168.88.2/32" ] } }
		// non-compact: "ipv4 unicast": [ { "nlri": "192.168.88.0/24" } ] } } } } }
		IPv4Unicast []interface{}             `json:"ipv4 unicast"`
		IPv4Flow    []IPv4FlowWithdrawMessage `json:"ipv4 flow"`
		IPv6Unicast []interface{}             `json:"ipv6 unicast"`
		IPv6Flow    []IPv6FlowWithdrawMessage `json:"ipv6 flow"`
		L2VPNVpls   []L2VPNVplsMessage        `json:"l2vpn vpls"`
	} `json:"withdraw"`
	EOR EORMessage `json:"eor"`
}

// EORMessage represents an End-of-RIB message
type EORMessage struct {
	// // { "eor": { "afi" : "ipv4", "safi" : "unicast" } }
	AFI  string `json:"afi"`
	SAFI string `json:"safi"`
}

// IPv4UnicastAnnounceMessage represents an ipv4 unicast family announcement
// messages are different between compact and non-compact mode
// compact: { "ipv4 unicast": { "192.168.1.184": [ "192.168.88.2/32" ] } }
// non-compact: "ipv4 unicast": [ { "nlri": "192.168.88.0/24" } ] } } } } }
type IPv4UnicastAnnounceMessage interface{}

// IPv4UnicastWithdrawMessage represents an ipv4 unicast family withdraw
// messages are different between compact and non-compact mode
// compact: { "ipv4 unicast": { [ "192.168.88.2/32" ] } }
// non-compact: "ipv4 unicast": [ { "nlri": "192.168.88.0/24" } ] } } } } }
type IPv4UnicastWithdrawMessage interface{}

// IPv4FlowAnnounceMessage represents an ipv4 flow family announcement
type IPv4FlowAnnounceMessage map[string][]IPv4FlowMessage

// IPv4FlowMessage represents an ipv4 flow family announcement
type IPv4FlowMessage struct {
	DestinationIPv4 []string `json:"destination-ipv4"`
	SourceIPv4      []string `json:"source-ipv4"`
	String          string   `json:"string"`
}

// IPv4FlowWithdrawMessage represents an ipv4 flow family withdraw
type IPv4FlowWithdrawMessage struct {
	DestinationIPv4 []string `json:"destination-ipv4"`
	SourceIPv4      []string `json:"source-ipv4"`
	String          string   `json:"string"`
}

// IPv6UnicastAnnounceMessage represents an ipv6 flow family announcement
type IPv6UnicastAnnounceMessage interface{}

// IPv6UnicastWithdrawMessage represents an ipv6 flow family withdraw
type IPv6UnicastWithdrawMessage interface{}

// IPv6FlowAnnounceMessage represents an ipv6 flow family announcement
type IPv6FlowAnnounceMessage map[string][]IPv6FlowMessage

// IPv6FlowMessage represents an ipv6 flow family announcement
type IPv6FlowMessage struct {
	DestinationIPv6 []string `json:"destination-ipv6"`
	SourceIPv6      []string `json:"source-ipv6"`
	String          string   `json:"string"`
}

// IPv6FlowWithdrawMessage represents an ipv6 flow family withdraw
type IPv6FlowWithdrawMessage struct {
	DestinationIPv6 []string `json:"destination-ipv6"`
	SourceIPv6      []string `json:"source-ipv6"`
	String          string   `json:"string"`
}

// L2VPNVplsMessage represents an l2vpn vpls family message (announce/withdraw same)
type L2VPNVplsMessage struct {
	/*
		"l2vpn vpls": {
			"192.168.201.1": [{
				"rd": "192.168.201.1:123",
				"endpoint": 5,
				"base": 10702,
				"offset": 1,
				"size": 8
			}]
		}
	*/
	RD       string `json:"rd"`
	Endpoint int    `json:"endpoint"`
	Offset   int    `json:"offset"`
	Size     int    `json:"size"`
}

// Attribute represent BGP attributes for a message
type Attribute struct {
	Med               int64 `json:"med"`
	ExtendedCommunity []struct {
		Value  json.Number `json:"value"`
		String string      `json:"string"`
	} `json:"extended-community"`
	Community         [][]int  `json:"community"`
	ASPath            []int    `json:"as-path"`
	ConfederationPath []int    `json:"confederation-path"`
	OriginatorID      string   `json:"originator-id"`
	LocalPreference   int      `json:"local-preference"`
	Origin            string   `json:"origin"`
	ClusterList       []string `json:"cluster-list"`
}
