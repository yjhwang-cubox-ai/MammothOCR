package main

import (
	"log"
	"net/http"

	"ocr-inference/internal/model"
	"ocr-inference/internal/service"
	"ocr-inference/pkg/utils"

	"github.com/gin-gonic/gin"
)

type OCRRequest struct {
	Image []byte `json:"image"`
}

type OCRResponse struct {
	Texts []string `json:"texts"`
	Error string   `json:"error,omitempty"`
}

func main() {
	// 모델 초기화
	detector, err := model.NewTextDetector("path/to/detector/model")
	if err != nil {
		log.Fatalf("Failed to initialize detector: %v", err)
	}

	recognizer, err := model.NewTextRecognizer("path/to/recognizer/model")
	if err != nil {
		log.Fatalf("Failed to initialize recognizer: %v", err)
	}

	// OCR 서비스 초기화
	ocrService := service.NewOCRService(detector, recognizer)

	// Gin 라우터 설정
	r := gin.Default()

	r.POST("/ocr", func(c *gin.Context) {
		var req OCRRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, OCRResponse{
				Error: "Invalid request format",
			})
			return
		}

		// 이미지 디코딩
		img, err := utils.BytesToImage(req.Image)
		if err != nil {
			c.JSON(http.StatusBadRequest, OCRResponse{
				Error: "Invalid image format",
			})
			return
		}

		// 텍스트 검출
		boxes, err := ocrService.DetectText(c.Request.Context(), img)
		if err != nil {
			c.JSON(http.StatusInternalServerError, OCRResponse{
				Error: err.Error(),
			})
			return
		}

		// 텍스트 인식
		texts, err := ocrService.RecognizeText(c.Request.Context(), img, boxes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, OCRResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, OCRResponse{
			Texts: texts,
		})
	})

	// 서버 시작
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
