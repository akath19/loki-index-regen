package main

import (
	"context"
	"log"
	"os"
)

//merchantPortal stores all required merchantPortal to be used by the app
type config struct {
	ProjectID     string
	BucketName    string
	Context       context.Context
	ContextCancel context.CancelFunc
	LokiAddress   string
	ChunksPath    string
}

//GetConfig parses environment variables and returns a merchantPortal struct
func getConfig() config {
	projectID := os.Getenv("PROJECT_ID")
	bucketName := os.Getenv("BUCKET_NAME")
	credentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	lokiAddress := os.Getenv("LOKI_ADDRESS")
	chunksPath := os.Getenv("CHUNKS_PATH")

	if projectID == "" {
		log.Fatal("PROJECT_ID environment variable not found or empty, please set PROJECT_ID to the Google Cloud Project ID where your bucket is located in")
	}

	if bucketName == "" {
		log.Fatal("BUCKET_NAME environment variable not found or empty, please set BUCKET_NAME to the name of the bucket where chunks will be downloaded from")
	}

	if credentials == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable not found or empty, please set GOOGLE_APPLICATION_CREDENTIALS to a path where a service account JSON file can be found")
	}

	if lokiAddress == "" {
		log.Fatal("LOKI_ADDRESS environment variable not found or empty, please set LOKI_ADDRESS to the Loki service address")
	}

	if chunksPath == "" {
		log.Fatal("CHUNKS_PATH environment variable not found or empty, please set CHUNKS_PATH to the path where chunks will be stored before processing")
	}

	ctx := context.Background()

	return config{
		ProjectID:   projectID,
		BucketName:  bucketName,
		Context:     ctx,
		LokiAddress: lokiAddress,
		ChunksPath:  chunksPath,
	}
}
