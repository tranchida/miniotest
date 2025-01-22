package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoint := "192.168.49.10:9000"
	accessKeyID := "ZCSuCGOy3GIN8YTxOwXJ"
	secretAccessKey := "73KUHqDWykOkpELJwWo7PQqyC3ZAn2piAt8f2XkP"

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

	minioClient.FPutObject(context.Background(), "sample", uuid.NewString(), "tmp/postgresql-42.7.5.jar", minio.PutObjectOptions{
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

	log.Println("File uploaded successfully")

	o, err := minioClient.GetObjectACL(context.Background(), "sample", "0c615b8c-9218-4121-95b1-b78f431853db")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(o.UserMetadata)

}
