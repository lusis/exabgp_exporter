package text

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var testSummaryDataFile = filepath.Join("testdata", "neighbor-summary.txt")

func TestParseSummaryTestData(t *testing.T) {
	file, err := ioutil.ReadFile(testSummaryDataFile)
	require.NoError(t, err)

	totalLines := testGetTotalLinesInFile(t, testSummaryDataFile)
	parsedEvents, err := SummariesFromBytes(file)
	require.NoError(t, err)
	require.Equal(t, totalLines-1, len(parsedEvents))
}
