package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"path/filepath"
	"strings"
)

const (
	ENDPOINT          = "https://oss-cn-shenzhen.aliyuncs.com"
	ACCESS_KEY_ID     = "my-access-key-id"
	ACCESS_KEY_SECRET = "my-access-key-secret"
	SAVE_DIR          = "app-pic/"
	BUCKET_NAME       = "my-bucket-name"
)

var (
	u string
	d string
	h bool
	l bool
)

func init() {
	flag.StringVar(&u, "u", "", "upload `file` to aliyun oss")
	flag.StringVar(&d, "d", "", "delete `file` from aliyun oss")

	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&l, "l", false, "list files")

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stdout, `A tool for uploading file to aliyun oss
Usage: alioss [-hl] [-u file] [-d file]

Options:
`)
	flag.PrintDefaults()
}

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

func uploadObj(client *oss.Client, bucketName string, localFile string) {
	fmt.Println("Upload:", localFile)
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	basename := filepath.Base(localFile)
	objectKey := SAVE_DIR + basename
	options := []oss.Option{oss.ObjectACL(oss.ACLPublicRead)}
	err = bucket.PutObjectFromFile(objectKey, localFile, options...)
	if err != nil {
		handleError(err)
	}
	picUrl := ENDPOINT[:8] + bucketName + "." + ENDPOINT[8:] + "/" + objectKey
	fmt.Println("Upload success:", picUrl)
}

func deleteObj(client *oss.Client, bucketName string, objNames []string) {
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	length := len(objNames)
	objectKeys := make([]string, length)
	for i := 0; i < length; i++ {
		objectKeys[i] = SAVE_DIR + objNames[i]
	}
	delRes, err := bucket.DeleteObjects(objectKeys)
	if err != nil {
		handleError(err)
	}
	fmt.Println("Delete resource:", delRes)
}

func listObj(client *oss.Client, bucketName string) {
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	lor, err := bucket.ListObjects(oss.Prefix(SAVE_DIR), oss.Delimiter("/"))
	if err != nil {
		handleError(err)
	}
	for idx, obj := range lor.Objects[1:] {
		segments := strings.Split(obj.Key, "/")
		filename := segments[len(segments)-1]
		picUrl := ENDPOINT[:8] + bucketName + "." + ENDPOINT[8:] + "/" + obj.Key
		fmt.Printf("%d [%s](%s)\n", idx+1, filename, picUrl)
	}
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		os.Exit(0)
	}

	client, err := oss.New(ENDPOINT, ACCESS_KEY_ID, ACCESS_KEY_SECRET)
	if err != nil {
		handleError(err)
	}
	bucketName := BUCKET_NAME

	if l {
		listObj(client, bucketName)
	}

	if u != "" {
		files := flag.Args()[:]
		files = append(files, u)
		for _, localFile := range files {
			uploadObj(client, bucketName, localFile)
		}
	}

	if d != "" {
		files := flag.Args()[:]
		files = append(files, d)
		deleteObj(client, bucketName, files)
	}

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}
}
