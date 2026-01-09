package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"dothefortune_server/internal/config"
	"dothefortune_server/internal/repository"
)

type AIService interface {
	GenerateFortuneText(fortuneMap map[string]string, todayStem, todayBranch string, category string) (string, error)
}

type aiService struct {
	cfg *config.Config
}

func NewAIService(fortuneRepo repository.FortuneRepository, userRepo repository.UserRepository, cfg *config.Config) AIService {
	return &aiService{
		cfg: cfg,
	}
}

func (s *aiService) GenerateFortuneText(fortuneMap map[string]string, todayStem, todayBranch string, category string) (string, error) {
	if s.cfg.GeminiAPIKey == "" {
		return "", errors.New("Gemini API key not configured")
	}

	prompt := s.buildFortunePrompt(fortuneMap, todayStem, todayBranch, category)
	return s.callGemini(prompt)
}

func (s *aiService) buildFortunePrompt(fortuneMap map[string]string, todayStem, todayBranch string, category string) string {
	dayStem := fortuneMap["day_stem"]
	dayBranch := fortuneMap["day_branch"]

	basePrompt := fmt.Sprintf(
		"당신은 따뜻하고 희망찬 조언을 주는 운세 AI입니다. 사용자의 사주 정보(일간: %s%s, 오늘의 일진: %s%s)를 바탕으로 %s에 대한 운세를 1~2문장으로 작성해주세요. 말투는 '~해요', '~할 수 있어요' 등 부드러운 경어체를 사용하고, 부정적인 운일 경우 '조심하세요'보다는 '잠시 쉬어가는 게 좋아요'처럼 우회적으로 표현해주세요.",
		dayStem, dayBranch, todayStem, todayBranch, category,
	)

	return basePrompt
}

type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (s *aiService) callGemini(prompt string) (string, error) {
	reqBody := GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", s.cfg.GeminiAPIKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error: %s", string(body))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no text in response")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

