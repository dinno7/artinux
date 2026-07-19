package storage

import (
	"cmp"
	"fmt"
	"strconv"

	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/minio/minio-go/v7"
)

func ObjectStorageToArtifact(obj *minio.ObjectInfo) *entities.Artifact {
	getMetaVal := func(key string) string {
		metaKey := getHeaderRealKey(key)
		return obj.Metadata.Get(metaKey)
	}

	meta := map[string]string{
		"object_key":    obj.Key,
		"arch":          getMetaVal("arch"),
		"os":            getMetaVal("os"),
		"ext":           getMetaVal("ext"),
		"name":          getMetaVal("name"),
		"checksum":      cmp.Or(getMetaVal("checksum"), obj.ChecksumSHA256),
		"file_size":     cmp.Or(getMetaVal("file_size"), strconv.FormatInt(obj.Size, 10)),
		"hostname":      getMetaVal("hostname"),
		"username":      getMetaVal("username"),
		"original_name": getMetaVal("original_name"),
		"uploaded_at":   getMetaVal("uploaded_at"),
	}

	return entities.ArtifactFromMap(meta)
}

func getHeaderRealKey(key string) string {
	return fmt.Sprintf("x-amz-meta-%s", key)
}
