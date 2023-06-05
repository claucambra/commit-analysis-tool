package authorgroups

import (
	"encoding/json"
	"os"
	"testing"
)

const testCommitsFile = "../../../test/data/log.txt"

const testGroupName = "VideoLAN"
const testGroupDomain = "videolan.org"

func testEmailGroups() map[string][]string {
	return map[string][]string{
		testGroupName: {testGroupDomain},
	}
}
func testGroupData(t *testing.T) *GroupData {
	marshalledTestGroupDataBytes, err := os.ReadFile("testdata/domaingroupsreport_test_orggroupdata.json")
	if err != nil {
		t.Fatalf("Could not read test marshalled group data file: %s", err)
	}

	var testGroupData GroupData
	err = json.Unmarshal(marshalledTestGroupDataBytes, &testGroupData)
	if err != nil {
		t.Fatalf("Could not unmarshall marshalled group data file: %s", err)
	}

	return &testGroupData
}

func testUnknownGroupData(t *testing.T) *GroupData {
	marshalledTestUnknownGroupDataBytes, err := os.ReadFile("testdata/domaingroupsreport_test_commgroupdata.json")
	if err != nil {
		t.Fatalf("Could not read test marshalled group data file: %s", err)
	}

	var testUnknownGroupData GroupData
	err = json.Unmarshal(marshalledTestUnknownGroupDataBytes, &testUnknownGroupData)
	if err != nil {
		t.Fatalf("Could not unmarshall marshalled group data file: %s", err)
	}

	return &testUnknownGroupData
}
