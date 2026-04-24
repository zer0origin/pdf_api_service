package dataapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pdf_service_api/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DataService struct {
	BaseUrl string
}

func (t DataService) SendMetaRequest(base64 string) (models.Meta, error) {
	if t.BaseUrl == "" {
		panic("No BaseUrl Provided")
	}

	url := fmt.Sprintf("%s/meta", t.BaseUrl)
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`{"base64": "%s"}`, base64))

	client := &http.Client{}
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return models.Meta{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return models.Meta{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Meta{}, err
	}

	data := &models.Meta{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return models.Meta{}, err
	}

	return *data, nil
}

type ExtractionRequest struct {
	DocumentUid           uuid.UUID                      `json:"documentUid,omitempty"`
	Base64EncodedDocument string                         `json:"base64EncodedDocument,omitempty"`
	Selections            map[uuid.UUID]models.Selection `json:"selections,omitempty"`
}

type ExtractionResponse map[uuid.UUID]struct {
	Text                string
	TextCoordinates     models.Coordinates
	SelectionCoordinate models.Coordinates
}

type TextData struct {
	Text            string
	TextCoordinates models.Coordinates
}

func (t DataService) SendBasicExtractionRequest(data ExtractionRequest) error {
	if t.BaseUrl == "" {
		panic("No BaseUrl Provided")
	}

	url := fmt.Sprintf("%s/extract", t.BaseUrl)
	method := "POST"

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	client := &http.Client{}
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(bytes)))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	fmt.Println(string(body))

	extractRes := &ExtractionResponse{}
	err = json.Unmarshal(body, extractRes)
	if err != nil {
		return err
	}

	fmt.Println(len(*extractRes))

	return err
}
