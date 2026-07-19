package storage

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"

	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/minio/minio-go/v7"
)

func ObjectStorageToArtifact(obj *minio.ObjectInfo) *entities.Artifact {
	meta := map[string]string{
		"object_key":    obj.Key,
		"arch":          getMetaVal(obj, "arch"),
		"os":            getMetaVal(obj, "os"),
		"ext":           getMetaVal(obj, "ext"),
		"name":          getMetaVal(obj, "name"),
		"checksum":      cmp.Or(getMetaVal(obj, "checksum"), obj.ChecksumSHA256),
		"file_size":     cmp.Or(getMetaVal(obj, "file_size"), strconv.FormatInt(obj.Size, 10)),
		"hostname":      getMetaVal(obj, "hostname"),
		"username":      getMetaVal(obj, "username"),
		"original_name": getMetaVal(obj, "original_name"),
		"uploaded_at":   getMetaVal(obj, "uploaded_at"),
	}

	return entities.ArtifactFromMap(meta)
}

func getMetaVal(obj *minio.ObjectInfo, key string) string {
	metaKey := getMetaHeaderKey(key)
	return cmp.Or(obj.Metadata.Get(metaKey), obj.UserMetadata[getMetaKeyPascalCase(metaKey)])
}

func getMetaHeaderKey(key string) string {
	return fmt.Sprintf("x-amz-meta-%s", key)
}

func getMetaKeyPascalCase(metaKey string) string {
	parts := strings.Split(strings.ToLower(metaKey), "-")
	for i := range parts {
		firstChar := string(parts[i][0])
		parts[i] = strings.ToUpper(firstChar) + parts[i][1:]
	}
	return strings.Join(parts, "-")
}
