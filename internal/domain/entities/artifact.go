package entities

import (
	"cmp"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Artifact struct {
	ObjectKey    string    `json:"object_key"`
	Name         string    `json:"name"`
	OriginalName string    `json:"original_name"`
	Extension    string    `json:"extension"`
	Size         int64     `json:"size"`
	OS           string    `json:"os"`
	Architecture string    `json:"architecture"`
	Username     string    `json:"username"`
	Hostname     string    `json:"hostname"`
	Checksum     string    `json:"checksum"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

func NewArtifact(
	fileName, extension,
	hostname, username,
	osName, arch string,
	size int64,
) (*Artifact, error) {
	osName = cmp.Or(osName, "unknown")
	arch = cmp.Or(arch, "unknown")
	uniqueFileName := uniqueName(fileName)
	hostname = cmp.Or(hostname, "unknown")
	username = cmp.Or(username, "unknown")
	uploadedAt := time.Now().UTC()
	objectKey := resolveObjectKey(uniqueFileName, osName, arch, uploadedAt)
	ext := filepath.Ext(fileName)

	return &Artifact{
		ObjectKey:    objectKey,
		Name:         uniqueFileName,
		OS:           osName,
		Architecture: arch,
		Extension:    ext,
		Size:         size,
		Username:     username,
		Hostname:     hostname,
		OriginalName: fileName,
		UploadedAt:   uploadedAt,
	}, nil
}

func (a *Artifact) AddChecksum(checksum string) *Artifact {
	a.Checksum = checksum
	return a
}

// Pattern: {os}/{arch}/{YYYY}/{MM}/{DD}/{checksum_prefix}_{original_filename}
func resolveObjectKey(fileName, osName, arch string, uploadedAt time.Time) string {
	return fmt.Sprintf(
		"%s/%s/%d/%d/%d/%s",
		osName,
		arch,
		uploadedAt.Year(),
		uploadedAt.Month(),
		uploadedAt.Day(),
		fileName,
	)
}

func uniqueName(name string) string {
	return fmt.Sprintf("%s_%s", uuid.New(), name)
}

func (a *Artifact) ToMap() map[string]string {
	m := map[string]string{
		"os":            a.OS,
		"arch":          a.Architecture,
		"username":      a.Username,
		"hostname":      a.Hostname,
		"name":          a.Name,
		"original_name": a.OriginalName,
		"checksum":      a.Checksum,
		"ext":           a.Extension,
		"file_size":     strconv.FormatInt(a.Size, 10),
		"uploaded_at":   a.UploadedAt.Format(time.RFC3339),
	}
	return m
}

func ArtifactFromMap(m map[string]string) *Artifact {
	size, _ := strconv.ParseInt(m["file_size"], 10, 64)
	uploadedAt, _ := time.Parse(time.RFC3339, m["uploaded_at"])
	return &Artifact{
		Name:         m["name"],
		ObjectKey:    m["object_key"],
		OS:           m["os"],
		Architecture: m["arch"],
		OriginalName: m["original_name"],
		Checksum:     m["checksum"],
		Username:     m["username"],
		Hostname:     m["hostname"],
		UploadedAt:   uploadedAt,
		Extension:    m["ext"],
		Size:         size,
	}
}
