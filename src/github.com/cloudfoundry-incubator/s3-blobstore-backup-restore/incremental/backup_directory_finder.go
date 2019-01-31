package incremental

import (
	"regexp"
	"sort"
)

//go:generate counterfeiter -o fakes/fake_backup_directory_finder.go . BackupDirectoryFinder
type BackupDirectoryFinder interface {
	ListBlobs() ([]BackedUpBlob, error)
}

type Finder struct {
	ID     string
	Bucket Bucket
}

func (b Finder) ListBlobs() ([]BackedUpBlob, error) {
	dirs, err := b.Bucket.ListDirectories()
	if err != nil {
		return nil, err
	}

	regex := regexp.MustCompile(`^\d{4}(_\d{2}){5}$`)

	var filteredDirs []string
	for _, dir := range dirs {
		if regex.MatchString(dir) {
			filteredDirs = append(filteredDirs, dir)
		}
	}

	if len(filteredDirs) == 0 {
		return nil, nil
	}

	lastComplete, err := b.findLastCompleteBackup(filteredDirs)
	if err != nil {
		return nil, err
	}

	if lastComplete == "" {
		return nil, nil
	}

	backupDirPath := joinBlobPath(lastComplete, b.ID)
	blobs, err := b.Bucket.ListBlobs(backupDirPath)
	if err != nil {
		return nil, err
	}

	var backedUpBlobs []BackedUpBlob
	for _, blob := range blobs {
		backedUpBlobs = append(backedUpBlobs, BackedUpBlob{
			Path:                blob.Path(),
			BackupDirectoryPath: backupDirPath,
		})
	}

	return backedUpBlobs, nil
}

func (b Finder) findLastCompleteBackup(backupDirectories []string) (string, error) {
	sort.Strings(backupDirectories)
	for _, dir := range backupDirectories {
		isComplete, err := b.Bucket.IsBackupComplete(dir)
		if err != nil {
			return "", err
		}

		if isComplete {
			return dir, nil
		}
	}

	return "", nil
}