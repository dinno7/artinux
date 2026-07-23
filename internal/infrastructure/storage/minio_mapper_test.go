package storage

import (
	"net/http"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

func TestGetMetaVal(t *testing.T) {
	// TEST: Get from Metadata
	obj := &minio.ObjectInfo{
		Metadata: http.Header{
			"X-Amz-Meta-File_size": []string{"InHeader"},
		},
		UserMetadata: minio.StringMap{
			"X-Amz-Meta-File_size": "InUserMetadataWithPrefix",
			"File_size":            "InUserMetadata",
		},
	}
	val := getMetaVal(obj, "file_size")
	assert.Equal(t, "InHeader", val)

	// TEST: Fallback to user metadata
	obj = &minio.ObjectInfo{
		Metadata: nil,
		UserMetadata: minio.StringMap{
			"X-Amz-Meta-File_size": "InUserMetadataWithPrefix",
			"File_size":            "InUserMetadata",
		},
	}
	val = getMetaVal(obj, "file_size")
	assert.Equal(t, "InUserMetadataWithPrefix", val)

	// TEST: Fallback to user metadata(with prefix)
	obj = &minio.ObjectInfo{
		UserMetadata: minio.StringMap{
			"File_size": "InUserMetadata",
		},
	}
	val = getMetaVal(obj, "file_size")
	assert.Equal(t, "InUserMetadata", val)

	// TEST: Not exists
	obj = &minio.ObjectInfo{}
	val = getMetaVal(obj, "file_size")
	assert.Equal(t, "", val)
}

func TestGetMetaKeyPascalCase(t *testing.T) {
	metaKey := getMetaKeyPascalCase("file_size")
	assert.Equal(t, "File_size", metaKey)

	metaKey = getMetaKeyPascalCase("ObjectKey")
	assert.Equal(t, "Objectkey", metaKey)

	metaKey = getMetaKeyPascalCase("x-amz-meta-file_size")
	assert.Equal(t, "X-Amz-Meta-File_size", metaKey)
}

func TestGetMetaHeaderKey(t *testing.T) {
	metaKey := getMetaHeaderKey("file_size")
	assert.Equal(t, "x-amz-meta-file_size", metaKey)

	metaKey = getMetaHeaderKey("File_siZe")
	assert.Equal(t, "x-amz-meta-file_size", metaKey)

	metaKey = getMetaHeaderKey("SomeKey")
	assert.Equal(t, "x-amz-meta-somekey", metaKey)
}
