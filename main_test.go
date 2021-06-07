package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeparatePathAndTag_NameOnly(t *testing.T) {
	// given
	srcName := "nginx"

	// when
	path, tag := separatePathAndTag(srcName)

	// then
	assert.NotEmpty(t, path)
	assert.NotEmpty(t, tag)
	assert.Equal(t, srcName, path)
	assert.Equal(t, "latest", tag)
}

func TestSeparatePathAndTag_NameAndTag(t *testing.T) {
	// given
	srcName := "nginx"
	srcTag := "1.21.0"

	// when
	path, tag := separatePathAndTag(fmt.Sprintf("%s:%s", srcName, srcTag))

	// then
	assert.NotEmpty(t, path)
	assert.NotEmpty(t, tag)
	assert.Equal(t, srcName, path)
	assert.Equal(t, srcTag, tag)
}

func TestSeparatePathAndTag_DomainNameAndTag(t *testing.T) {
	// given
	srcDomain := "grafana"
	srcName := "grafana"
	srcTag := "7.5.0"

	// when
	path, tag := separatePathAndTag(fmt.Sprintf("%s/%s:%s", srcDomain, srcName, srcTag))

	// then
	assert.NotEmpty(t, path)
	assert.NotEmpty(t, tag)
	assert.Equal(t, fmt.Sprintf("%s/%s", srcDomain, srcName), path)
	assert.Equal(t, srcTag, tag)
}

func TestSeparatePathAndTag_RegistryDomainNameAndTag(t *testing.T) {
	// given
	srcRegistry := "docker.io"
	srcDomain := "grafana"
	srcName := "grafana"
	srcTag := "7.5.0"

	// when
	path, tag := separatePathAndTag(fmt.Sprintf("%s/%s/%s:%s", srcRegistry, srcDomain, srcName, srcTag))

	// then
	assert.NotEmpty(t, path)
	assert.NotEmpty(t, tag)
	assert.Equal(t, fmt.Sprintf("%s/%s/%s", srcRegistry, srcDomain, srcName), path)
	assert.Equal(t, srcTag, tag)
}

func TestBuildBackupRegistryImgStr_NameOnly(t *testing.T) {
	// given
	srcName := "nginx"

	// when
	newImg := buildBackupRegistryImgStr(srcName)

	// then
	assert.NotEmpty(t, newImg)
	assert.Equal(t, fmt.Sprintf("%s/%s/%s:%s", backupRegistryUrl, backupDomain, srcName, "latest"), newImg)
}

func TestBuildBackupRegistryImgStr_NameAndTag(t *testing.T) {
	// given
	srcName := "nginx"
	srcTag := "1.21.0"

	// when
	newImg := buildBackupRegistryImgStr(fmt.Sprintf("%s:%s", srcName, srcTag))

	// then
	assert.NotEmpty(t, newImg)
	assert.Equal(t, fmt.Sprintf("%s/%s/%s:%s", backupRegistryUrl, backupDomain, srcName, srcTag), newImg)
}

func TestBuildBackupRegistryImgStr_DomainNameAndTag(t *testing.T) {
	// given
	srcDomain := "grafana"
	srcName := "grafana"
	srcTag := "7.5.0"

	// when
	newImg := buildBackupRegistryImgStr(fmt.Sprintf("%s/%s:%s", srcDomain, srcName, srcTag))

	// then
	assert.NotEmpty(t, newImg)
	assert.Equal(t, fmt.Sprintf("%s/%s/%s:%s", backupRegistryUrl, backupDomain, srcName, srcTag), newImg)
}

func TestBuildBackupRegistryImgStr_RegistryDomainNameAndTag(t *testing.T) {
	// given
	srcRegistry := "docker.io"
	srcDomain := "grafana"
	srcName := "grafana"
	srcTag := "7.5.0"

	// when
	newImg := buildBackupRegistryImgStr(fmt.Sprintf("%s/%s/%s:%s", srcRegistry, srcDomain, srcName, srcTag))

	// then
	assert.NotEmpty(t, newImg)
	assert.Equal(t, fmt.Sprintf("%s/%s/%s:%s", backupRegistryUrl, backupDomain, srcName, srcTag), newImg)
}
