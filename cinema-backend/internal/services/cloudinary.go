package services

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService(cloudinaryURL string) (*CloudinaryService, error) {
	if cloudinaryURL == "" {
		return nil, nil
	}
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}
	return &CloudinaryService{cld: cld}, nil
}

func (s *CloudinaryService) IsEnabled() bool {
	return s != nil && s.cld != nil
}

func (s *CloudinaryService) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	publicID := strings.TrimSuffix(header.Filename, "."+getExt(header.Filename))

	result, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:  publicID,
		Folder:    "cinema",
		Overwrite: BoolPtr(true),
	})
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}

func (s *CloudinaryService) Delete(ctx context.Context, url string) error {
	publicID := extractPublicID(url)
	if publicID == "" {
		return nil
	}
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

func getExt(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return ""
	}
	return filename[idx+1:]
}

func extractPublicID(url string) string {
	parts := strings.Split(url, "/")
	for i, p := range parts {
		if p == "cinema" && i+1 < len(parts) {
			id := strings.Join(parts[i:], "/")
			id = strings.TrimSuffix(id, "."+getExt(id))
			return id
		}
	}
	return ""
}

func BoolPtr(b bool) *bool {
	return &b
}
