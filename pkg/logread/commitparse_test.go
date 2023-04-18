package logread

import (
	"reflect"
	"testing"
	"time"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
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

	expectedCommitData := new(common.CommitData)
	expectedCommitData.Id = "1c915e7dd147d4b060c2c241bb966d6f6c6ecde9"
	expectedCommitData.AuthorTime = time.Date(2023, 4, 8, 17, 47, 43, 0, expectedCommitAuthorLocation).Unix()
	expectedCommitData.AuthorName = "Claudio Cambra"
	expectedCommitData.AuthorEmail = "developer@claudiocambra.com"
	expectedCommitData.CommitterTime = time.Date(2023, 4, 12, 23, 21, 43, 0, expectedCommitCommitterLocation).Unix()
	expectedCommitData.CommitterName = "Jean-Baptiste Kempf"
	expectedCommitData.CommitterEmail = "jb@videolan.org"
	expectedCommitData.NumInsertions = 32
	expectedCommitData.NumDeletions = 36
	expectedCommitData.NumFilesChanged = 6

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

func TestParseCommitLog(t *testing.T) {
	testCommitLog := `e130355483777647ab5fca208631dc18f4ef0d0a__SEPARATOR__Sun, 16 Apr 2023 17:28:40 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Sun, 16 Apr 2023 17:28:40 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com
modules/gui/macosx/library/audio-library/VLCLibraryCollectionViewAudioGroupSupplementaryDetailView.m | 9 +++++++--
1 file changed, 7 insertions(+), 2 deletions(-)

68e12e00434561f03fbfd58bac96404a873e5a18__SEPARATOR__Sun, 16 Apr 2023 17:28:13 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Sun, 16 Apr 2023 17:28:13 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com
modules/gui/macosx/library/audio-library/VLCLibraryAudioDataSource.m | 32 ++++++++++++++++++--------------
1 file changed, 18 insertions(+), 14 deletions(-)

91e0c789d8983c7c9c4a8874e8a3ce1f2ac4d330__SEPARATOR__Sun, 16 Apr 2023 16:58:03 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Sun, 16 Apr 2023 17:00:57 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com
modules/gui/macosx/library/audio-library/VLCLibraryAudioGroupDataSource.h |  2 +-
modules/gui/macosx/library/audio-library/VLCLibraryAudioGroupDataSource.m | 14 +++++++-------
2 files changed, 8 insertions(+), 8 deletions(-)

5129448f614ff5ac6475d23f77e700056d6e4076__SEPARATOR__Sun, 16 Apr 2023 16:49:17 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Sun, 16 Apr 2023 16:49:17 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com
modules/gui/macosx/library/audio-library/VLCLibraryAlbumTableCellView.m                         | 7 ++++++-
modules/gui/macosx/library/audio-library/VLCLibraryCollectionViewAlbumSupplementaryDetailView.m | 8 +++++++-
2 files changed, 13 insertions(+), 2 deletions(-)

d1de6dc7eab11db6accdf2b9ca2d30b91f41eb84__SEPARATOR__Sun, 16 Apr 2023 16:40:32 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Sun, 16 Apr 2023 16:40:32 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com
modules/gui/macosx/library/audio-library/VLCLibraryAlbumTableCellView.m                         | 5 +++--
modules/gui/macosx/library/audio-library/VLCLibraryCollectionViewAlbumSupplementaryDetailView.m | 6 +++---
2 files changed, 6 insertions(+), 5 deletions(-)

a256c8b2efbc82db85cd403457e63b2b9b5dec37__SEPARATOR__Sun, 16 Apr 2023 16:39:59 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Sun, 16 Apr 2023 16:39:59 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com
modules/gui/macosx/library/audio-library/VLCLibraryAlbumTracksDataSource.h |  3 +++
modules/gui/macosx/library/audio-library/VLCLibraryAlbumTracksDataSource.m | 14 +++++++++++++-
2 files changed, 16 insertions(+), 1 deletion(-)

254097627ba476a81ae4c628347a8c19130086c7__SEPARATOR__Sun, 16 Apr 2023 16:30:34 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com__SEPARATOR__Sun, 16 Apr 2023 16:30:34 +0800__SEPARATOR__Claudio Cambra__SEPARATOR__developer@claudiocambra.com
modules/gui/macosx/library/audio-library/VLCLibraryAlbumTracksDataSource.m | 5 ++++-
1 file changed, 4 insertions(+), 1 deletion(-)`

	expectedCommitCount := 7
	parsedCommitLog, err := ParseCommitLog(testCommitLog)
	receivedCommitCount := len(parsedCommitLog)

	if err != nil {
		t.Fatalf("Received error parsing commit log: %s", err)
	}

	if receivedCommitCount != expectedCommitCount {
		t.Fatalf(`Received a different amount of commits than expected. 
			Expected %d, received %d`, expectedCommitCount, receivedCommitCount)
	}
}
