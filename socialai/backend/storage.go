package backend

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"socialai/util"

	"cloud.google.com/go/storage"
)

// MediaBackend lets the course use GCS in the cloud while keeping a practical
// local file-storage option for development on a personal computer.
var MediaBackend MediaStorage

type MediaStorage interface {
	Save(reader io.Reader, objectName string) (string, error)
	Delete(mediaURL string) error
}

type localStorageBackend struct {
	directory     string
	publicBaseURL string
}

type googleCloudStorageBackend struct {
	client *storage.Client
	bucket string
}

func InitMediaBackend(config *util.ApplicationConfig) error {
	if config.StorageConfig != nil && config.StorageConfig.Driver == "local" {
		directory := config.StorageConfig.Directory
		if directory == "" {
			directory = "uploads"
		}
		if err := os.MkdirAll(directory, 0o755); err != nil {
			return err
		}

		publicBaseURL := config.StorageConfig.PublicBaseURL
		if configuredURL := os.Getenv("SOCIALAI_PUBLIC_BASE_URL"); configuredURL != "" {
			publicBaseURL = configuredURL
		}

		MediaBackend = &localStorageBackend{
			directory:     directory,
			publicBaseURL: strings.TrimRight(publicBaseURL, "/"),
		}
		return nil
	}

	if config.GCSConfig == nil || config.GCSConfig.Bucket == "" {
		return fmt.Errorf("a GCS bucket is required when storage.driver is not local")
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		return err
	}

	MediaBackend = &googleCloudStorageBackend{
		client: client,
		bucket: config.GCSConfig.Bucket,
	}
	return nil
}

func (backend *localStorageBackend) Save(reader io.Reader, objectName string) (string, error) {
	path := filepath.Join(backend.directory, objectName)
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", err
	}

	if backend.publicBaseURL == "" {
		return "/uploads/" + url.PathEscape(objectName), nil
	}
	return backend.publicBaseURL + "/" + url.PathEscape(objectName), nil
}

func (backend *localStorageBackend) Delete(mediaURL string) error {
	objectName, err := mediaObjectName(mediaURL)
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join(backend.directory, objectName))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (backend *googleCloudStorageBackend) Save(reader io.Reader, objectName string) (string, error) {
	ctx := context.Background()
	object := backend.client.Bucket(backend.bucket).Object(objectName)
	writer := object.NewWriter(ctx)

	if _, err := io.Copy(writer, reader); err != nil {
		_ = writer.Close()
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	// The course uses object ACLs. This requires uniform bucket-level access
	// to be disabled, exactly as discussed in the GCS lesson.
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", backend.bucket, url.PathEscape(objectName)), nil
}

func (backend *googleCloudStorageBackend) Delete(mediaURL string) error {
	objectName, err := mediaObjectName(mediaURL)
	if err != nil {
		return err
	}

	return backend.client.Bucket(backend.bucket).Object(objectName).Delete(context.Background())
}

func mediaObjectName(mediaURL string) (string, error) {
	parsedURL, err := url.Parse(mediaURL)
	if err != nil {
		return "", err
	}

	objectName, err := url.PathUnescape(filepath.Base(parsedURL.Path))
	if err != nil {
		return "", err
	}
	if objectName == "" || objectName == "." {
		return "", fmt.Errorf("invalid media URL")
	}
	return objectName, nil
}
