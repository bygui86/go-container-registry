package main

import (
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// DockerHub, Quay, GCP GCR, AWS ECR, Azure Container Registry
var (
	// containerRegistries = []string{"docker.io", "quay.io", "gcr.io", "eu.gcr.io", "us.gcr.io", "asia.gcr.io", "ecr.aws", "mcr.microsoft.com"}

	backupRegistryUrl      = "docker.io"
	backupDomain           = "bygui86"
	backupRegistryUsername = ""
	backupRegistryPassword = ""

	nginxImgStr        = "nginx:1.21.0"
	nginxImgStrNoTag   = "nginx"
	grafanaImgStr      = "grafana/grafana:7.5.0"
	grafanaImgStrNoTag = "grafana/grafana"
	polarisImgStr      = "quay.io/fairwinds/polaris:3.2.1"
	polarisImgStrNoTag = "quay.io/fairwinds/polaris"
)

func main() {
	fmt.Println("Start go-containerregistry tests")

	// listTags(nginxImgStrNoTag, false)
	// listTags(grafanaImgStrNoTag, false)
	// listTags(polarisImgStrNoTag, false)
	//
	// pull(nginxImgStr, false)
	// pull(grafanaImgStr, false)
	// pull(polarisImgStr, false)
	//
	// copy(nginxImgStr, true)
	// copy(grafanaImgStr, true)

	// across different registries
	polarisImg := pull(polarisImgStr, false)
	push(polarisImg, buildBackupRegistryImgStr(polarisImgStr), true)

	fmt.Println("go-containerregistry tests completed")
}

func listTags(imgStr string, auth bool) {
	fmt.Printf("List tags of %s \n", imgStr)

	var tags []string
	var err error
	if auth {
		tags, err = crane.ListTags(
			imgStr,
			crane.WithAuth(&authn.Basic{
				Username: backupRegistryUsername,
				Password: backupRegistryPassword,
			}),
		)
	} else {
		tags, err = crane.ListTags(imgStr)
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s tags: %s \n", imgStr, strings.Join(tags, ", "))
}

func pull(imgStr string, auth bool) v1.Image {
	fmt.Printf("Pull %s \n", imgStr)

	var img v1.Image
	var pullErr error
	if auth {
		img, pullErr = crane.Pull(
			imgStr,
			crane.WithAuth(&authn.Basic{
				Username: backupRegistryUsername,
				Password: backupRegistryPassword,
			}),
		)
	} else {
		img, pullErr = crane.Pull(imgStr)
	}
	if pullErr != nil {
		panic(pullErr)
	}

	hash, hashErr := img.Digest()
	if hashErr != nil {
		panic(hashErr)
	}

	fmt.Printf("%s - digest %s %s \n", imgStr, hash.Algorithm, hash.Hex)
	return img
}

func push(img v1.Image, imgStr string, auth bool) {
	fmt.Printf("Push %s \n", imgStr)

	var err error
	if auth {
		err = crane.Push(
			img, imgStr,
			crane.WithAuth(&authn.Basic{
				Username: backupRegistryUsername,
				Password: backupRegistryPassword,
			}),
		)
	} else {
		err = crane.Push(img, imgStr)
	}
	if err != nil {
		panic(err)
	}
}

func copy(srcImg string, auth bool) {
	tgtImg := buildBackupRegistryImgStr(srcImg)

	fmt.Printf("Copy %s to %s \n", srcImg, tgtImg)

	var err error
	if auth {
		err = crane.Copy(
			srcImg, tgtImg,
			crane.WithAuth(&authn.Basic{
				Username: backupRegistryUsername,
				Password: backupRegistryPassword,
			}),
		)
	} else {
		err = crane.Copy(srcImg, tgtImg)
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copy %s to %s COMPLETED \n", srcImg, tgtImg)
}

func copyAcrossRegistries(srcImg string) {
	tgtImg := buildBackupRegistryImgStr(srcImg)

	fmt.Printf("Copy %s to %s \n", srcImg, tgtImg)

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
		panic(fmt.Errorf("invalid img path %s", path))
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
