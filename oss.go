package main

import (
	"fmt"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Client struct {
	bucket *oss.Bucket
}

func NewClient(endpoint, accessKeyID, accessKeySecret, bucket string) (*Client, error) {
	c, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	b, err := c.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	return &Client{bucket: b}, nil
}

func (c *Client) Upload(localPath string) error {
	key := filepath.Base(localPath)
	return c.bucket.PutObjectFromFile(key, localPath)
}

func (c *Client) Download(ossKey string) error {
	destPath := filepath.Base(ossKey)
	err := c.bucket.GetObjectToFile(ossKey, destPath)
	if err != nil {
		if ossErr, ok := err.(oss.ServiceError); ok && ossErr.StatusCode == 404 {
			return fmt.Errorf("object not found: %s", ossKey)
		}
		return err
	}
	return nil
}
