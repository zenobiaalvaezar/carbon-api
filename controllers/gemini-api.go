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

	c.Response().Header().Set("Content-Type", "text/plain")
	c.Response().WriteHeader(http.StatusOK)

	_, _ = c.Response().Write([]byte("Generating image... Please wait.\n"))
	c.Response().Flush()

	client := resty.New()
	apiKey := os.Getenv("OPENAI_API_KEY")

	prompt := "a badge with the name 'Fathur Rohman' and a label 'Pelindung Hutan' (tree protector)"
	response := OpenAIResponse{}
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":   "dall-e-3",
			"prompt":  prompt,
			"n":       1,
			"size":    "1024x1024",
			"quality": "standard",
		}).
		SetResult(&response).
		Post("https://api.openai.com/v1/images/generations")

	if err != nil || resp.StatusCode() != http.StatusOK {
		_, _ = c.Response().Write([]byte("Failed to generate image.\n"))
		c.Response().Flush()
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate image"})
	}

	imageURL := response.Data[0].URL

	_, _ = c.Response().Write([]byte("Downloading image...\n"))
	c.Response().Flush()

	imageResp, err := http.Get(imageURL)
	if err != nil || imageResp.StatusCode != http.StatusOK {
		_, _ = c.Response().Write([]byte("Failed to download image.\n"))
		c.Response().Flush()
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to download image"})
	}
	defer imageResp.Body.Close()

	imageData, err := io.ReadAll(imageResp.Body)
	if err != nil {
		_, _ = c.Response().Write([]byte("Failed to read image data.\n"))
		c.Response().Flush()
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read image data"})
	}

	subject := "Generated Image"
	body := fmt.Sprintf("Here is the image you requested for prompt: %s", req.Prompt)

	_, _ = c.Response().Write([]byte("Sending email...\n"))
	c.Response().Flush()

	err = utils.SendEmailWithAttachment("fr081938@gmail.com", subject, body, imageData)
	if err != nil {
		_, _ = c.Response().Write([]byte("Failed to send email.\n"))
		c.Response().Flush()
		fmt.Print(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to send email"})
	}

	_, _ = c.Response().Write([]byte("Image generated and sent successfully!\n"))
	c.Response().Flush()

	return nil
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
