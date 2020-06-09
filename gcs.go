package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"path/filepath"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type gcsClient struct {
	gcs    *storage.Client
	config config
	parser *chunksParser
}

func (c *gcsClient) createClient(config config) {
	client, err := storage.NewClient(config.Context)

	if err != nil {
		log.Fatalf("Failed to create GCS client; %v", err)
	}

	c.gcs = client
	c.config = config
}

func (c *gcsClient) close() {
	err := c.gcs.Close()

	if err != nil {
		log.Fatal(err)
	}
}

func (c *gcsClient) downloadAndProcessChunks() {
	c.parser = &chunksParser{
		LokiAddress: c.config.LokiAddress,
	}

	bucket := c.gcs.Bucket(c.config.BucketName)

	query := &storage.Query{Prefix: ""}
	err := query.SetAttrSelection([]string{"Name", "Updated"})

	if err != nil {
		log.Fatal(err)
	}

	objs := bucket.Objects(c.config.Context, query)

	for {
		attrs, err := objs.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		path := c.saveChunk(bucket, attrs.Name)

		c.parser.ParseAndSaveChunk(path)
	}
}

func (c *gcsClient) saveChunk(bucket *storage.BucketHandle, name string) string {
	obj, err := bucket.Object(name).NewReader(c.config.Context)

	if err != nil {
		log.Fatal(err)
	}

	file, err := ioutil.ReadAll(obj)
	defer obj.Close()

	if err != nil {
		log.Fatal(err)
	}

	hasher := md5.New()
	hasher.Write([]byte(name))

	path := filepath.FromSlash(c.config.ChunksPath + hex.EncodeToString(hasher.Sum(nil)))

	err = ioutil.WriteFile(path, file, 0644)

	if err != nil {
		log.Fatalf("cannot write chunk: %v", err)
	} else {
		log.Printf("chunk %v saved successfully\n", name)
	}
	return path
}
