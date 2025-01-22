package main

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoint := "minio-api.test"
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")

	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now setup

	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	for _, bucket := range buckets {
		log.Println(bucket.Name)

		objects := minioClient.ListObjects(context.Background(), bucket.Name, minio.ListObjectsOptions{
			WithMetadata: true,
		})
		for object := range objects {
			if object.Err != nil {
				log.Println(object.Err)
				return
			}
			log.Println("key", object.Key)
			log.Println("UserTag", object.UserTags)
			log.Println("UserMetadata", object.UserMetadata)
			log.Println()

			// for k, v := range object.UserMetadata {
			// 	log.Printf("%s: %s\n", k, v)
			// }
		}
	}

	/*
		err = minioClient.FGetObject(context.Background(), "sample", "2025-01-21/docker-desktop-amd64.deb", "tmp/myobject.jar", minio.GetObjectOptions{})
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("File downloaded successfully")
	*/

	info, err := minioClient.FPutObject(context.Background(), "sample", uuid.NewString(), "tmp/nrjmx_linux_2.6.0_noarch.tar.gz", minio.PutObjectOptions{
		ContentType: "application/java-archive",
		UserTags: map[string]string{
			"application-tag": "miniotest-tag",
			"domain-tag":      "minio-tag",
		},
		UserMetadata: map[string]string{
			"application": "miniotest",
			"domain":      "minio",
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("File uploaded successfully",info.VersionID)

	o, err := minioClient.GetObjectACL(context.Background(), "sample", "0c615b8c-9218-4121-95b1-b78f431853db")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(o.UserMetadata)

}
