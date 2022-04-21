//go:build integration
// +build integration

package githubutil

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	stdlibversion "github.com/flyteorg/flytestdlib/version"

	"github.com/flyteorg/flytectl/pkg/platformutil"

	"github.com/stretchr/testify/assert"
)

var sandboxImageName = "cr.flyte.org/flyteorg/flyte-sandbox"
var sandboxManifest = "flyte_sandbox_manifest.yaml"

func TestGetLatestVersion(t *testing.T) {
	t.Run("Get latest release with wrong url", func(t *testing.T) {
		_, err := GetLatestVersion("fl")
		assert.NotNil(t, err)
	})
	t.Run("Get latest release", func(t *testing.T) {
		_, err := GetLatestVersion("flytectl")
		assert.Nil(t, err)
	})
}

func TestGetLatestRelease(t *testing.T) {
	release, err := GetLatestVersion("flyte")
	assert.Nil(t, err)
	assert.Equal(t, true, strings.HasPrefix(release.GetTagName(), "v"))
}

func TestCheckVersionExist(t *testing.T) {
	t.Run("Invalid Tag", func(t *testing.T) {
		_, err := CheckVersionExist("v100.0.0", "flyte")
		assert.NotNil(t, err)
	})
	t.Run("Valid Tag", func(t *testing.T) {
		release, err := CheckVersionExist("v0.15.0", "flyte")
		assert.Nil(t, err)
		assert.Equal(t, true, strings.HasPrefix(release.GetTagName(), "v"))
	})
}

func TestGetFullyQualifiedImageName(t *testing.T) {
	t.Run("Get tFully Qualified Image Name ", func(t *testing.T) {
		image, tag, err := GetFullyQualifiedImageName("dind", "", sandboxImageName, false)
		assert.Nil(t, err)
		assert.Equal(t, true, strings.HasPrefix(tag, "v"))
		assert.Equal(t, true, strings.HasPrefix(image, sandboxImageName))
	})
	t.Run("Get tFully Qualified Image Name with pre release", func(t *testing.T) {
		image, tag, err := GetFullyQualifiedImageName("dind", "", sandboxImageName, true)
		assert.Nil(t, err)
		assert.Equal(t, true, strings.HasPrefix(tag, "v"))
		assert.Equal(t, true, strings.HasPrefix(image, sandboxImageName))
	})
	t.Run("Get tFully Qualified Image Name with specific version", func(t *testing.T) {
		image, tag, err := GetFullyQualifiedImageName("dind", "v0.19.0", sandboxImageName, true)
		assert.Nil(t, err)
		assert.Equal(t, "v0.19.0", tag)
		assert.Equal(t, true, strings.HasPrefix(image, sandboxImageName))
	})
}

func TestGetSHAFromVersion(t *testing.T) {
	t.Run("Invalid Tag", func(t *testing.T) {
		_, err := GetSHAFromVersion("v100.0.0", "flyte")
		assert.NotNil(t, err)
	})
	t.Run("Valid Tag", func(t *testing.T) {
		release, err := GetSHAFromVersion("v0.15.0", "flyte")
		assert.Nil(t, err)
		assert.Greater(t, len(release), 0)
	})
}

func TestGetAssetsFromRelease(t *testing.T) {
	t.Run("Successful get assets", func(t *testing.T) {
		assets, err := GetAssetsFromRelease("v0.15.0", sandboxManifest, flyte)
		assert.Nil(t, err)
		assert.NotNil(t, assets)
		assert.Equal(t, sandboxManifest, *assets.Name)
	})

	t.Run("Failed get assets with wrong name", func(t *testing.T) {
		assets, err := GetAssetsFromRelease("v0.15.0", "test", flyte)
		assert.NotNil(t, err)
		assert.Nil(t, assets)
	})
	t.Run("Successful get assets with wrong version", func(t *testing.T) {
		assets, err := GetAssetsFromRelease("v100.15.0", "test", flyte)
		assert.NotNil(t, err)
		assert.Nil(t, assets)
	})
}

func TestGetAssetsName(t *testing.T) {
	t.Run("Get Assets name", func(t *testing.T) {
		expected := fmt.Sprintf("flytectl_%s_386.tar.gz", cases.Title(language.English).String(runtime.GOOS))
		arch = platformutil.Arch386
		assert.Equal(t, expected, getFlytectlAssetName())
	})
}

func TestCheckBrewInstall(t *testing.T) {
	symlink, err := CheckBrewInstall(platformutil.Darwin)
	assert.Nil(t, err)
	assert.Equal(t, len(symlink), 0)
	symlink, err = CheckBrewInstall(platformutil.Linux)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(symlink))
}

func TestGetUpgradeMessage(t *testing.T) {
	var darwin = platformutil.Darwin
	var linux = platformutil.Linux
	var windows = platformutil.Linux

	var version = "v0.2.20"
	stdlibversion.Version = "v0.2.10"
	message, err := GetUpgradeMessage(version, darwin)
	assert.Nil(t, err)
	assert.Equal(t, 157, len(message))

	version = "v0.2.09"
	message, err = GetUpgradeMessage(version, darwin)
	assert.Nil(t, err)
	assert.Equal(t, 63, len(message))

	version = "v"
	message, err = GetUpgradeMessage(version, darwin)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(message))

	version = "v0.2.20"
	message, err = GetUpgradeMessage(version, windows)
	assert.Nil(t, err)
	assert.Equal(t, 157, len(message))

	version = "v0.2.20"
	message, err = GetUpgradeMessage(version, linux)
	assert.Nil(t, err)
	assert.Equal(t, 157, len(message))
}