package libimage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizePlatform(t *testing.T) {
	type platform struct {
		os, arch, variant string
	}
	for _, test := range []struct {
		input, expected platform
	}{
		{
			platform{"", "", ""},
			platform{"", "", ""},
		},
		{
			platform{"foo", "", "garbage"},
			platform{"foo", "", "garbage"},
		},
		{
			platform{"&", "invalid", "os"},
			platform{"&", "invalid", "os"},
		},
		{
			platform{"linux", "", ""},
			platform{"linux", "", ""},
		},
		{
			platform{"LINUX", "", ""},
			platform{"linux", "", ""},
		},
		{
			platform{"", "aarch64", ""},
			platform{"", "arm64", ""},
		},
		{
			platform{"macos", "x86_64", ""},
			platform{"darwin", "amd64", ""},
		},
		{
			platform{"linux", "amd64", ""},
			platform{"linux", "amd64", ""},
		},
		{
			platform{"linux", "arm64", "v8"},
			platform{"linux", "arm64", "v8"},
		},
		{
			platform{"linux", "aarch64", ""},
			platform{"linux", "arm64", ""},
		},
	} {
		os, arch, variant := NormalizePlatform(test.input.os, test.input.arch, test.input.variant)
		assert.Equal(t, test.expected.os, os, test.input)
		assert.Equal(t, test.expected.arch, arch, test.input)
		assert.Equal(t, test.expected.variant, variant, test.input)
	}
}

func TestNormalizeName(t *testing.T) {
	const digestSuffix = "@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	for _, c := range []struct{ input, expected string }{
		{"#", ""}, // Clearly invalid
		{"example.com/busybox", "example.com/busybox:latest"},                                  // Qualified name-only
		{"example.com/busybox:notlatest", "example.com/busybox:notlatest"},                     // Qualified name:tag
		{"example.com/busybox" + digestSuffix, "example.com/busybox" + digestSuffix},           // Qualified name@digest
		{"example.com/busybox:notlatest" + digestSuffix, "example.com/busybox" + digestSuffix}, // Qualified name:tag@digest
		{"busybox:latest", "localhost/busybox:latest"},                                         // Unqualified name-only
		{"busybox:latest" + digestSuffix, "localhost/busybox" + digestSuffix},                  // Unqualified name:tag@digest
		{"localhost/busybox", "localhost/busybox:latest"},                                      // Qualified with localhost
		{"ns/busybox:latest", "localhost/ns/busybox:latest"},                                   // Unqualified with a dot-less namespace
		{"docker.io/busybox:latest", "docker.io/library/busybox:latest"},                       // docker.io without /library/
	} {
		res, err := NormalizeName(c.input)
		if c.expected == "" {
			assert.Error(t, err, c.input)
		} else {
			require.NoError(t, err, c.input)
			assert.Equal(t, c.expected, res.String())
		}
	}
}

func TestNormalizeTaggedDigestedString(t *testing.T) {
	const digestSuffix = "@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	for _, test := range []struct{ input, expected string }{
		{"$$garbage", ""},
		{"fedora", "fedora"},
		{"fedora:tag", "fedora:tag"},
		{digestSuffix, ""},
		{"docker://fedora:latest", ""},
		{"docker://fedora:latest" + digestSuffix, ""},
		{"fedora" + digestSuffix, "fedora" + digestSuffix},
		{"fedora:latest" + digestSuffix, "fedora" + digestSuffix},
		{"repo/fedora:123456" + digestSuffix, "repo/fedora" + digestSuffix},
		{"quay.io/repo/fedora:tag" + digestSuffix, "quay.io/repo/fedora" + digestSuffix},
		{"localhost/fedora:anothertag" + digestSuffix, "localhost/fedora" + digestSuffix},
		{"localhost:5000/fedora:v1.2.3.4.5" + digestSuffix, "localhost:5000/fedora" + digestSuffix},
	} {
		res, err := normalizeTaggedDigestedString(test.input)
		if test.expected == "" {
			assert.Error(t, err, "%v", test)
		} else {
			assert.NoError(t, err, "%v", test)
			assert.Equal(t, test.expected, res, "%v", test)
		}
	}
}
