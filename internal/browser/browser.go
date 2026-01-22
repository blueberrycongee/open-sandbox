package browser

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

type Service struct {
	CDPAddress string
}

func NewService(cdpAddress string) *Service {
	return &Service{CDPAddress: cdpAddress}
}

func (service *Service) Info() string {
	return service.CDPAddress
}

func (service *Service) Screenshot(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White)
	return png.Encode(file, img)
}
