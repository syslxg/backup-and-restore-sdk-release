module azure-blobstore-backup-restore

require (
	github.com/Azure/azure-pipeline-go v0.2.1 // indirect
	github.com/Azure/azure-storage-blob-go v0.0.0-20180507052152-66ba96e49ebb
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pkg/errors v0.8.1 // indirect
	system-tests v0.0.0
)

replace system-tests => ../system-tests

replace azure-blobstore-backup-restore => ./

replace s3-blobstore-backup-restore => ../s3-blobstore-backup-restore