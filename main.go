package main

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
)

// DockerHub, Quay, GCP GCR, AWS ECR, Azure Container Registry
var (
	containerRegistries = []string{"docker.io", "quay.io", "gcr.io", "eu.gcr.io", "us.gcr.io", "asia.gcr.io", "ecr.aws", "mcr.microsoft.com"}

	backupRegistryUrl      = "docker.io"
	backupDomain           = "bygui86"
	backupRegistryUsername = ""
	backupRegistryPassword = ""

	nginxImgStr   = "nginx:1.21.0"
	grafanaImgStr = "grafana/grafana:7.5.0"
)

func main() {
	fmt.Println("Start go-containerregistry tests")

	// pull(nginxImgStr)

	// listTags(nginxImgStr)

	copy(nginxImgStr)

	copy(grafanaImgStr)

	fmt.Println("go-containerregistry tests completed")
}

func pull(imgStr string) {
	fmt.Printf("Pull %s \n", imgStr)

	img, pullErr := crane.Pull(imgStr)
	if pullErr != nil {
		panic(pullErr)
	}

	hash, hasErr := img.Digest()
	if hasErr != nil {
		panic(hasErr)
	}

	fmt.Printf("%s - digest %s %s \n", imgStr, hash.Algorithm, hash.Hex)
}

func listTags(imgStr string) {
	fmt.Printf("List tags of %s \n", imgStr)

	tags, tagsErr := crane.ListTags(imgStr)
	if tagsErr != nil {
		panic(tagsErr)
	}

	fmt.Printf("%s tags: %s \n", imgStr, strings.Join(tags, ", "))
}

func copy(imgStr string) {
	newImgStr := buildBackupRegistryImgStr(imgStr)

	fmt.Printf("Copy %s to %s \n", imgStr, newImgStr)

	auth := &authn.Basic{
		Username: backupRegistryUsername,
		Password: backupRegistryPassword,
	}

	err := crane.Copy(
		imgStr, newImgStr,
		crane.WithAuth(auth),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copy %s to %s COMPLETED \n", imgStr, newImgStr)
}

func buildBackupRegistryImgStr(imgStr string) string {
	path, tag := separatePathAndTag(imgStr)
	pathParts := strings.Split(path, "/")
	var name string
	switch length := len(pathParts); {
	case length == 3:
		name = pathParts[2]
	case length == 2:
		name = pathParts[1]
	case length == 1:
		name = pathParts[0]
	default:
		panic(fmt.Errorf("invalid img path %s", path)) // TODO improve
	}

	return fmt.Sprintf("%s/%s/%s:%s", backupRegistryUrl, backupDomain, name, tag)
}

func separatePathAndTag(imgStr string) (string, string) {
	tagSplit := strings.Split(imgStr, ":")

	if len(tagSplit) == 1 { // [0] path parts - tag = latest
		return tagSplit[0], "latest"

	} else { // [0] path parts - [1] tag
		return tagSplit[0], tagSplit[1]
	}
}
