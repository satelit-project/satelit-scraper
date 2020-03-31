package spider

import (
	"bytes"
	"compress/flate"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/minio/minio-go/v6"

	"shitty.moe/satelit-project/satelit-scraper/config"
	"shitty.moe/satelit-project/satelit-scraper/logging"
	"shitty.moe/satelit-project/satelit-scraper/proto/data"
)

type Cache interface {
	AddHTML(data *bytes.Buffer, source data.Source, id int32) error
}

type S3Cache struct {
	// Storage configuration.
	cfg *config.Storage

	// S3 client.
	client *minio.Client

	// Logger
	log *logging.Logger
}

// Creates and returns new storage object or error if initialization failed.
func NewS3Cache(cfg *config.Storage, log *logging.Logger) (S3Cache, error) {
	secure := true
	if strings.HasPrefix(cfg.Host, "localhost") || strings.HasPrefix(cfg.Host, "127.0.0.1") {
		secure = false
	}

	client, err := minio.New(cfg.Host, cfg.Key, cfg.Secret, secure)
	if err != nil {
		return S3Cache{}, err
	}

	return S3Cache{cfg, client, log.With("s3cache", cfg.Bucket)}, nil
}

func (c S3Cache) AddHTML(data *bytes.Buffer, source data.Source, id int32) error {
	compressed, err := compress(data)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s/cached/%d.flate", strings.ToLower(source.String()), id)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.cfg.UploadTimeout)*time.Second)
	defer cancel()

	c.log.Infof("caching html: %s (%d)", source, id)
	_, err = c.client.PutObjectWithContext(ctx, c.cfg.Bucket, name, compressed, int64(compressed.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	return err
}

func compress(data *bytes.Buffer) (*bytes.Buffer, error) {
	var compressed bytes.Buffer
	zw, err := flate.NewWriter(&compressed, flate.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = zw.Write(data.Bytes())
	if err != nil {
		return nil, err
	}

	return &compressed, nil
}
