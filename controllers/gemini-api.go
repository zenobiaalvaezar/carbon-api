package controllers

import (
	"carbon-api/utils"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

type GeminiAPIController struct{}

func NewGeminiAPIController() *GeminiAPIController {
	return &GeminiAPIController{}
}

type RequestPayload struct {
	Prompt string `json:"prompt"`
}

type ResponsePayload struct {
	Response string `json:"response"`
}

func generateContentA() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(`Analisis data emisi karbon berikut dan berikan prediksi untuk emisi di masa depan serta rekomendasi untuk mencapai netralitas karbon. Prediksi harus didasarkan pada catatan sebelumnya dan mencakup langkah-langkah yang dapat dilakukan.
		Data:

		Total Emisi: 38.901
		Jumlah Pohon yang Dibutuhkan: 155
		Catatan Terakhir:
		Rekor 1: Emisi Bahan Bakar = 1.905, Emisi Listrik = 203, Total Emisi = 109
		Rekor 2: Emisi Bahan Bakar = 1.905, Emisi Listrik = 203, Total Emisi = 109
		Format Output:

		Prediksi emisi untuk 6 bulan ke depan.
		Garis waktu untuk mencapai netralitas karbon jika tren saat ini berlanjut.
		Rekomendasi untuk mengurangi emisi.`,
	))
	if err != nil {
		log.Fatal(err)
	}

	printResponse(resp)
}

func (ctrl *GeminiAPIController) GeminiAPI(c echo.Context) error {
	var reqPayload RequestPayload
	if err := c.Bind(&reqPayload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	generateContentA()

	return c.JSON(http.StatusOK, "")
}

func initGenAIClient() (*genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return client, nil
}

type GenerateImageRequest struct {
	Prompt string `json:"prompt" validate:"required"`
	Size   string `json:"size" validate:"required,oneof=256x256 512x512 1024x1024"`
}

type OpenAIResponse struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func (ctrl *GeminiAPIController) GenerateImage(c echo.Context) error {
	// Parse request body
	var req GenerateImageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Call OpenAI API
	client := resty.New()
	apiKey := os.Getenv("OPENAI_API_KEY")

	prompt := "Generate an image of a user profile card with a forest conservation theme, featuring a clean and professional design. Key elements: a prominent badge displaying 'pelindung hutan,' two metrics ('donation_tree': 2 trees and 'carbon_tree': 1 unit of carbon), user name 'Fathur Rohman,' and user ID '0012.' The design should be modern and minimalist, incorporating elements like leaves, trees, or a stylized forest backdrop."
	response := OpenAIResponse{}
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"prompt": prompt,
			"n":      1,
			"size":   "1024x1024",
		}).
		SetResult(&response).
		Post("https://api.openai.com/v1/images/generations")

	if err != nil || resp.StatusCode() != http.StatusOK {
		errorMessage := "Failed to generate image"
		if resp.StatusCode() != http.StatusOK {
			errorMessage = resp.String()
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": errorMessage})
	}

	imageURL := response.Data[0].URL

	// Download the image
	imageResp, err := http.Get(imageURL)
	if err != nil || imageResp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to download image"})
	}
	defer imageResp.Body.Close()

	imageData, err := io.ReadAll(imageResp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read image data"})
	}

	// Email recipient details
	subject := "Generated Image"
	body := fmt.Sprintf("Here is the image you requested for prompt: %s", req.Prompt)

	// Send email with image attachment
	err = utils.SendEmailWithAttachment("fr081938@gmail.com", subject, body, imageData)
	if err != nil {
		fmt.Print(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to send email"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Image generated and sent successfully!"})
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}
