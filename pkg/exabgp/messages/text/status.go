package text

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

/*
Peer            AS        up/down state       |     #sent     #recvd
127.0.0.1       64496        down idle                  0          0
192.168.1.1     64496     0:00:01 established          45          0
*/
var summaryHeaderLine = `Peer            AS        up/down state       |     #sent     #recvd`
var rxSummary = `(?P<peer_ip>\S+)\s+(?P<peer_as>\d+)\s+(?P<status>.*)\s+(?P<state>idle|active|connect|opensent|openconfirm|established)\s+(?P<sent>\d+)\s+(?P<recvd>\d+)$`

func parseSummaryLine(s string) (map[string]string, error) {
	md := make(map[string]string)
	re := regexp.MustCompile(rxSummary)
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

func SummaryEntryFromString(s string) (*NeighborSummary, error) {
	ns := &NeighborSummary{}
	md, err := parseSummaryLine(s)
	if err != nil {
		return nil, err
	}
	ns.IPAddress = md["peer_ip"]
	ns.AS = md["peer_as"]
	if md["status"] == "down" {
		ns.Status = "down"
	} else {
		ns.Status = "up"
	}
	// we need to account for pre-4.0.8 summary reporting for down connections
	if md["state"] != "established" {
		ns.Status = "down"
	}
	ns.State = md["state"]

	ns.Sent, _ = strconv.Atoi(md["sent"])
	ns.Received, _ = strconv.Atoi(md["recvd"])
	return ns, nil
}

func SummariesFromBytes(b []byte) ([]*NeighborSummary, error) {
	var sums []*NeighborSummary
	reader := bufio.NewReader(bytes.NewReader(b))
	for {
		l, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if string(l) != summaryHeaderLine {
			r, err := SummaryEntryFromString(string(l))
			if err != nil {
				return sums, err
			}
			sums = append(sums, r)
		}
	}
	return sums, nil
}

// Neighbor represents a neighbor summary
type NeighborSummary struct {
	IPAddress string
	AS        string
	Status    string
	State     string
	Sent      int
	Received  int
}
