package text

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var testRibDataFile = filepath.Join("testdata", "rib-out.txt")

func testGetTotalLinesInFile(t *testing.T, f string) int {
	file, err := os.Open(f)
	// nolint:errcheck,staticcheck
	defer file.Close()

	require.NoError(t, err)

	s := bufio.NewScanner(file)
	totalLines := 0
	for s.Scan() {
		totalLines++
	}
	return totalLines
}

func TestParseRibTestData(t *testing.T) {
	file, err := ioutil.ReadFile(testRibDataFile)
	require.NoError(t, err)

	totalLines := testGetTotalLinesInFile(t, testRibDataFile)

	ribs, err := RibFromBytes(file)
	require.NoError(t, err)
	require.Equal(t, totalLines, len(ribs))

}

func TestParseRibString(t *testing.T) {
	var testString = `neighbor 127.0.0.1 local-ip 127.0.0.1 local-as 64496 peer-as 64496 router-id 1.1.1.1 family-allowed in-open ipv4 unicast 192.168.88.248/29 next-hop self med 100`
	m, err := RibEntryFromString(testString)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, "127.0.0.1", m.PeerIP)
	require.Equal(t, "64496", m.PeerAS)
	require.Equal(t, "127.0.0.1", m.LocalIP)
	require.Equal(t, "64496", m.LocalAS)
	require.Equal(t, "ipv4", m.AFI)
	require.Equal(t, "unicast", m.SAFI)
	require.Equal(t, "192.168.88.248/29 next-hop self med 100", m.Details)
	require.Equal(t, "ipv4 unicast", m.Family())
}

func TestParseIPv4UnicastFull(t *testing.T) {
	var testString = `neighbor 127.0.0.1 local-ip 127.0.0.1 local-as 64496 peer-as 64496 router-id 1.1.1.1 family-allowed in-open ipv4 unicast 192.168.88.248/29 next-hop self med 100`
	m, err := RibEntryFromString(testString)
	require.NoError(t, err)
	require.NotNil(t, m)
	ipv4, err := m.IPv4Unicast()
	require.NoError(t, err)
	require.NotNil(t, ipv4)
	require.Equal(t, "192.168.88.248/29", ipv4.NLRI)
	require.Equal(t, "self", ipv4.NextHop)
	require.Equal(t, "med 100", ipv4.Attributes)
}

func TestParseIPv4UnicastNoAttributes(t *testing.T) {
	var testString = `neighbor 127.0.0.1 local-ip 127.0.0.1 local-as 64496 peer-as 64496 router-id 1.1.1.1 family-allowed in-open ipv4 unicast 192.168.88.248/29 next-hop self`
	m, err := RibEntryFromString(testString)
	require.NoError(t, err)
	require.NotNil(t, m)
	ipv4, err := m.IPv4Unicast()
	require.NoError(t, err)
	require.NotNil(t, ipv4)
	require.Equal(t, "192.168.88.248/29", ipv4.NLRI)
	require.Equal(t, "self", ipv4.NextHop)
	require.Empty(t, ipv4.Attributes)
}

func TestParseIPv6RibString(t *testing.T) {
	var testString = `neighbor 2001::2 local-ip 2001::1 local-as 64496 peer-as 64496 router-id 1.1.1.1 family-allowed in-open ipv6 unicast 2001:db8:1000::/64 next-hop self med 100`
	m, err := RibEntryFromString(testString)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, "2001::2", m.PeerIP)
	require.Equal(t, "64496", m.PeerAS)
	require.Equal(t, "2001::1", m.LocalIP)
	require.Equal(t, "64496", m.LocalAS)
	require.Equal(t, "ipv6", m.AFI)
	require.Equal(t, "unicast", m.SAFI)
	require.Equal(t, "2001:db8:1000::/64 next-hop self med 100", m.Details)
	require.Equal(t, "ipv6 unicast", m.Family())
}

func TestParseIPv6UnicastFull(t *testing.T) {
	var testString = `neighbor 2001::2 local-ip 2001::1 local-as 64496 peer-as 64496 router-id 1.1.1.1 family-allowed in-open ipv6 unicast 2001:db8:1000::/64 next-hop self med 100`
	m, err := RibEntryFromString(testString)
	require.NoError(t, err)
	require.NotNil(t, m)
	ipv6, err := m.IPv6Unicast()
	require.NoError(t, err)
	require.NotNil(t, ipv6)
	require.Equal(t, "2001:db8:1000::/64", ipv6.NLRI)
	require.Equal(t, "self", ipv6.NextHop)
	require.Equal(t, "med 100", ipv6.Attributes)
}

func TestParseIPv6UnicastNoAttributes(t *testing.T) {
	var testString = `neighbor 2001::2 local-ip 2001::1 local-as 64496 peer-as 64496 router-id 1.1.1.1 family-allowed in-open ipv6 unicast 2001:db8:1000::/64 next-hop self`
	m, err := RibEntryFromString(testString)
	require.NoError(t, err)
	require.NotNil(t, m)
	ipv6, err := m.IPv6Unicast()
	require.NoError(t, err)
	require.NotNil(t, ipv6)
	require.Equal(t, "2001:db8:1000::/64", ipv6.NLRI)
	require.Equal(t, "self", ipv6.NextHop)
	require.Empty(t, ipv6.Attributes)
}
