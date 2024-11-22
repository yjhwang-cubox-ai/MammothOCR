package service

import (
	"context"
	"fmt"
	"image"
	"ocr-inference/pkg/utils"
)

type OCRService interface {
	DetectText(ctx context.Context, img image.Image) ([]BoundingBox, error)
	RecognizeText(ctx context.Context, img image.Image, boxes []BoundingBox) ([]string, error)
}

type BoundingBox struct {
	X1, Y1, X2, Y2 float32
}

type ocrService struct {
	detector   *model.TextDetector
	recognizer *model.TextRecognizer
}

func NewOCRService(detector *model.TextDetector, recognizer *model.TextRecognizer) OCRService {
	return &ocrService{
		detector:   detector,
		recognizer: recognizer,
	}
}

func (s *ocrService) DetectText(ctx context.Context, img image.Image) ([]BoundingBox, error) {
	boxes, err := s.detector.Detect(img)
	if err != nil {
		return nil, fmt.Errorf("text detection failed: %v", err)
	}
	return boxes, nil
}

func (s *ocrService) RecognizeText(ctx context.Context, img image.Image, boxes []BoundingBox) ([]string, error) {
	var results []string
	for _, box := range boxes {
		// 이미지 크롭
		croppedImg := utils.CropImage(img, box)

		// 텍스트 인식
		text, err := s.recognizer.Recognize(croppedImg)
		if err != nil {
			return nil, fmt.Errorf("text recognition failed: %v", err)
		}
		results = append(results, text)
	}
	return results, nil
}
