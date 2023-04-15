package git

import (
	"reflect"
	"testing"
	"time"
)

func TestParseCommit(t *testing.T) {
	testCommit := `1c915e7dd147d4b060c2c241bb966d6f6c6ecde9__SEPARATOR__Sat, 8 Apr 2023 17:47:43 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Wed, 12 Apr 2023 23:21:43 +0000__SEPARATOR__Jean-Baptiste Kempf__SEPARATOR__jb@videolan.org
modules/gui/macosx/library/VLCLibraryWindow.h                            |  6 +++---
modules/gui/macosx/library/VLCLibraryWindowPersistentPreferences.h       | 22 +++++++++-------------
modules/gui/macosx/library/VLCLibraryWindowPersistentPreferences.m       | 30 +++++++++++++++---------------
modules/gui/macosx/library/audio-library/VLCLibraryAudioViewController.m |  4 ++--
modules/gui/macosx/library/media-source/VLCMediaSourceBaseDataSource.m   |  4 ++--
modules/gui/macosx/library/video-library/VLCLibraryVideoViewController.m |  2 +-
6 files changed, 32 insertions(+), 36 deletions(-)`

	expectedCommitAuthorLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("Received an error while creating expected commit author timezone: %s", err)
	}

	expectedCommitCommitterLocation, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatalf("Received an error while creating expected commit committer timezone: %s", err)
	}

	expectedCommitData := new(CommitData)
	expectedCommitData.id = "1c915e7dd147d4b060c2c241bb966d6f6c6ecde9"
	expectedCommitData.authorTime = time.Date(2023, 4, 8, 17, 47, 43, 0, expectedCommitAuthorLocation).Unix()
	expectedCommitData.authorName = "Claudio Cambra"
	expectedCommitData.authorEmail = "developer@claudiocambra.com"
	expectedCommitData.committerTime = time.Date(2023, 4, 12, 23, 21, 43, 0, expectedCommitCommitterLocation).Unix()
	expectedCommitData.committerName = "Jean-Baptiste Kempf"
	expectedCommitData.committerEmail = "jb@videolan.org"
	expectedCommitData.numInsertions = 32
	expectedCommitData.numDeletions = 36

	commitData, err := ParseCommit(testCommit)
	if err != nil {
		t.Fatalf("Received an error while parsing commit: %s", err)
	}

	if !reflect.DeepEqual(commitData, expectedCommitData) {
		t.Fatalf(`Parsed commit does not equal expected commit.
			Expected: %+v
			Received: %+v`, expectedCommitData, commitData)
	}
}
