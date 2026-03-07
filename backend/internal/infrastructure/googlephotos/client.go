package googlephotos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	uploadEndpoint      = "https://photoslibrary.googleapis.com/v1/uploads"
	batchCreateEndpoint = "https://photoslibrary.googleapis.com/v1/mediaItems:batchCreate"
	tokenEndpoint       = "https://oauth2.googleapis.com/token"
)

// Client はGoogle Photosアップロードクライアントです。
type Client struct {
	httpClient    *http.Client
	clientID      string
	clientSecret  string
	refreshToken  string
	targetAlbumID string
}

// NewClientFromEnv は環境変数からGoogle Photosクライアントを作成します。
func NewClientFromEnv() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		clientID:   strings.TrimSpace(os.Getenv("GOOGLE_PHOTOS_CLIENT_ID")),
		clientSecret: strings.TrimSpace(
			os.Getenv("GOOGLE_PHOTOS_CLIENT_SECRET"),
		),
		refreshToken:  strings.TrimSpace(os.Getenv("GOOGLE_PHOTOS_REFRESH_TOKEN")),
		targetAlbumID: strings.TrimSpace(os.Getenv("GOOGLE_PHOTOS_ALBUM_ID")),
	}
}

// Enabled はアップロード設定の有効状態を返します。
func (c *Client) Enabled() bool {
	return c.clientID != "" && c.clientSecret != "" && c.refreshToken != ""
}

// Upload は画像をGoogle Photosへアップロードし、閲覧URLを返します。
func (c *Client) Upload(ctx context.Context, filename string, data []byte) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("Google Photos設定が不足しています")
	}
	if len(data) == 0 {
		return "", fmt.Errorf("画像データが空です")
	}

	accessToken, err := c.fetchAccessToken(ctx)
	if err != nil {
		return "", err
	}

	uploadToken, err := c.uploadRawBytes(ctx, accessToken, filename, data)
	if err != nil {
		return "", err
	}

	return c.createMediaItem(ctx, accessToken, uploadToken, filename)
}

func (c *Client) fetchAccessToken(ctx context.Context) (string, error) {
	values := url.Values{}
	values.Set("client_id", c.clientID)
	values.Set("client_secret", c.clientSecret)
	values.Set("refresh_token", c.refreshToken)
	values.Set("grant_type", "refresh_token")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("アクセストークン取得に失敗しました: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("アクセストークン取得エラー: %s", string(body))
	}

	var payload struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	if payload.AccessToken == "" {
		return "", fmt.Errorf("アクセストークンが空です")
	}

	return payload.AccessToken, nil
}

func (c *Client) uploadRawBytes(ctx context.Context, accessToken string, filename string, data []byte) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadEndpoint, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", filename)
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Google Photosアップロードに失敗しました: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("Google Photosアップロードエラー: %s", string(body))
	}

	uploadToken := strings.TrimSpace(string(body))
	if uploadToken == "" {
		return "", fmt.Errorf("uploadTokenが取得できませんでした")
	}

	return uploadToken, nil
}

func (c *Client) createMediaItem(ctx context.Context, accessToken string, uploadToken string, filename string) (string, error) {
	requestBody := map[string]any{
		"newMediaItems": []map[string]any{
			{
				"description": "fish-tech admin upload",
				"simpleMediaItem": map[string]any{
					"fileName":    filename,
					"uploadToken": uploadToken,
				},
			},
		},
	}
	if c.targetAlbumID != "" {
		requestBody["albumId"] = c.targetAlbumID
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, batchCreateEndpoint, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("mediaItems作成に失敗しました: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("mediaItems作成エラー: %s", string(body))
	}

	var response struct {
		NewMediaItemResults []struct {
			Status struct {
				Code    int32  `json:"code"`
				Message string `json:"message"`
			} `json:"status"`
			MediaItem struct {
				ProductURL string `json:"productUrl"`
				BaseURL    string `json:"baseUrl"`
			} `json:"mediaItem"`
		} `json:"newMediaItemResults"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	if len(response.NewMediaItemResults) == 0 {
		return "", fmt.Errorf("mediaItems結果が空です")
	}

	item := response.NewMediaItemResults[0]
	if item.Status.Code != 0 {
		return "", fmt.Errorf("Google Photos保存エラー: %s", item.Status.Message)
	}
	if item.MediaItem.ProductURL != "" {
		return item.MediaItem.ProductURL, nil
	}
	if item.MediaItem.BaseURL != "" {
		return item.MediaItem.BaseURL, nil
	}

	return "", fmt.Errorf("Google Photos URLが取得できませんでした")
}
