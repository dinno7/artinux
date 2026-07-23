package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveObjectKey(t *testing.T) {
	now := time.Date(2026, 7, 23, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		fileName string
		osName   string
		arch     string
		time     time.Time
		want     string
	}{
		{
			name:     "linux amd64",
			fileName: "abc123_app.deb",
			osName:   "linux",
			arch:     "amd64",
			time:     now,
			want:     "linux/amd64/2026/7/23/abc123_app.deb",
		},
		{
			name:     "darwin arm64",
			fileName: "xyz_bin.gz",
			osName:   "darwin",
			arch:     "arm64",
			time:     time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			want:     "darwin/arm64/2025/12/1/xyz_bin.gz",
		},
		{
			name:     "windows 386",
			fileName: "setup.exe",
			osName:   "windows",
			arch:     "386",
			time:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			want:     "windows/386/2024/1/15/setup.exe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveObjectKey(tt.fileName, tt.osName, tt.arch, tt.time)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUniqueName(t *testing.T) {
	name := uniqueName("myfile.deb")
	// UUID is 36 characters, followed by underscore and the original name
	assert.Len(t, name, 36+1+len("myfile.deb"))
	assert.Contains(t, name, "_myfile.deb")

	// Each call generates a different name
	another := uniqueName("myfile.deb")
	assert.NotEqual(t, name, another)
}

func TestNewArtifact(t *testing.T) {
	t.Run("fills defaults for empty fields", func(t *testing.T) {
		artifact, err := NewArtifact("pkg.deb", "deb", "", "", "", "", 1024)
		require.NoError(t, err)
		require.NotNil(t, artifact)

		assert.Equal(t, ".deb", artifact.Extension)
		assert.Equal(t, int64(1024), artifact.Size)
		assert.Contains(t, artifact.Name, "_pkg.deb")
		assert.NotEmpty(t, artifact.ObjectKey)
		assert.NotZero(t, artifact.UploadedAt)

		// Defaults
		assert.Equal(t, "unknown", artifact.OS)
		assert.Equal(t, "unknown", artifact.Architecture)
		assert.Equal(t, "unknown", artifact.Hostname)
		assert.Equal(t, "unknown", artifact.Username)
	})

	t.Run("preserves provided values", func(t *testing.T) {
		artifact, err := NewArtifact("server.log", "log", "web-01", "deploy", "linux", "amd64", 4096)
		require.NoError(t, err)
		require.NotNil(t, artifact)

		assert.Equal(t, "linux", artifact.OS)
		assert.Equal(t, "amd64", artifact.Architecture)
		assert.Equal(t, "web-01", artifact.Hostname)
		assert.Equal(t, "deploy", artifact.Username)
	})

	t.Run("uses file extension from original name", func(t *testing.T) {
		artifact, err := NewArtifact("archive.tar.gz", "tar.gz", "", "", "", "", 512)
		require.NoError(t, err)
		require.NotNil(t, artifact)

		assert.Equal(t, ".gz", artifact.Extension)
	})
}

func TestAddChecksum(t *testing.T) {
	artifact, err := NewArtifact("f.deb", "deb", "", "", "", "", 100)
	require.NoError(t, err)

	result := artifact.AddChecksum("abc123checksum")
	assert.Same(t, artifact, result)
	assert.Equal(t, "abc123checksum", artifact.Checksum)
}

func TestArtifactToMapAndBack(t *testing.T) {
	original, err := NewArtifact("pkg.deb", "deb", "host-01", "admin", "linux", "arm64", 2048)
	require.NoError(t, err)
	original.AddChecksum("sha256checksum")

	m := original.ToMap()

	assert.Equal(t, "linux", m["os"])
	assert.Equal(t, "arm64", m["arch"])
	assert.Equal(t, "admin", m["username"])
	assert.Equal(t, "host-01", m["hostname"])
	assert.Equal(t, original.Name, m["name"])
	assert.Equal(t, "pkg.deb", m["original_name"])
	assert.Equal(t, "sha256checksum", m["checksum"])
	assert.Equal(t, ".deb", m["ext"])
	assert.Equal(t, "2048", m["file_size"])
	assert.NotEmpty(t, m["uploaded_at"])

	restored := ArtifactFromMap(m)
	assert.Equal(t, original.Name, restored.Name)
	assert.Equal(t, original.OS, restored.OS)
	assert.Equal(t, original.Architecture, restored.Architecture)
	assert.Equal(t, original.Username, restored.Username)
	assert.Equal(t, original.Hostname, restored.Hostname)
	assert.Equal(t, original.Checksum, restored.Checksum)
	assert.Equal(t, ".deb", restored.Extension)
	assert.Equal(t, original.Size, restored.Size)
	assert.WithinDuration(t, original.UploadedAt, restored.UploadedAt, time.Second)
}

func TestArtifactFromMapWithMissingFields(t *testing.T) {
	m := map[string]string{
		"name": "test",
		"os":   "linux",
	}

	restored := ArtifactFromMap(m)
	require.NotNil(t, restored)
	assert.Equal(t, "test", restored.Name)
	assert.Equal(t, "linux", restored.OS)
	assert.Empty(t, restored.Checksum)
	assert.Equal(t, int64(0), restored.Size)
}
