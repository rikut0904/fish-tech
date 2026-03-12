package hotpepper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	placeDomain "fish-tech/internal/domain/place"
)

const searchEndpoint = "https://webservice.recruit.co.jp/hotpepper/gourmet/v1/"
const smallAreaEndpoint = "https://webservice.recruit.co.jp/hotpepper/small_area/v1/"
const middleAreaEndpoint = "https://webservice.recruit.co.jp/hotpepper/middle_area/v1/"

// Client はHotPepperグルメサーチAPIクライアントです。
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// SmallArea はHotPepperのsmall_area情報です。
type SmallArea struct {
	Code          string
	Name          string
	MiddleArea    string
	LargeAreaCode string
}

// NewClientFromEnv は環境変数からHotPepperクライアントを作成します。
func NewClientFromEnv() *Client {
	return &Client{
		apiKey:     strings.TrimSpace(os.Getenv("HOTPEPPER_API_KEY")),
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

// SearchRecommendedPlaces はおすすめ店舗を検索します。
func (c *Client) SearchRecommendedPlaces(ctx context.Context, condition placeDomain.SearchCondition) ([]placeDomain.RecommendedPlace, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("HOTPEPPER_API_KEY が設定されていません")
	}

	values := url.Values{}
	values.Set("key", c.apiKey)
	values.Set("format", "json")
	values.Set("keyword", condition.Keyword)
	values.Set("count", strconv.Itoa(condition.Count))
	values.Set("start", strconv.Itoa((condition.Page-1)*condition.Count+1))
	values.Set("order", "4")
	if condition.LargeArea != "" {
		values.Set("large_area", condition.LargeArea)
	}
	if condition.SmallAreaCode != "" {
		values.Set("small_area", condition.SmallAreaCode)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchEndpoint+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HotPepper APIへの接続に失敗しました: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("HotPepper APIエラー: %s", string(body))
	}

	var payload struct {
		Results struct {
			Error []struct {
				Message string `json:"message"`
			} `json:"error"`
			Shop []struct {
				ID      string          `json:"id"`
				Name    string          `json:"name"`
				Address string          `json:"address"`
				Lat     json.RawMessage `json:"lat"`
				Lng     json.RawMessage `json:"lng"`
				Genre   struct {
					Name string `json:"name"`
				} `json:"genre"`
				SmallArea struct {
					Code string `json:"code"`
				} `json:"small_area"`
				Card       string `json:"card"`
				CouponURLs struct {
					PC string `json:"pc"`
				} `json:"coupon_urls"`
				Photo struct {
					PC struct {
						L string `json:"l"`
					} `json:"pc"`
					Mobile struct {
						L string `json:"l"`
					} `json:"mobile"`
				} `json:"photo"`
			} `json:"shop"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if len(payload.Results.Error) > 0 {
		return nil, fmt.Errorf("HotPepper APIエラー: %s", payload.Results.Error[0].Message)
	}

	places := make([]placeDomain.RecommendedPlace, 0, len(payload.Results.Shop))
	for _, shop := range payload.Results.Shop {
		places = append(places, placeDomain.RecommendedPlace{
			ID:            shop.ID,
			Name:          shop.Name,
			Address:       shop.Address,
			Lat:           parseJSONValueToString(shop.Lat),
			Lng:           parseJSONValueToString(shop.Lng),
			Coupon:        shop.CouponURLs.PC,
			Genre:         shop.Genre.Name,
			Card:          shop.Card,
			Logo:          firstNonEmpty(shop.Photo.PC.L, shop.Photo.Mobile.L),
			SmallAreaCode: strings.TrimSpace(shop.SmallArea.Code),
			LargeAreaCode: strings.TrimSpace(condition.LargeArea),
		})
	}

	return places, nil
}

// IsValidSmallAreaCode は指定コードが大エリア配下のsmall_areaとして存在するかを返します。
func (c *Client) IsValidSmallAreaCode(ctx context.Context, largeAreaCode string, smallAreaCode string) (bool, error) {
	if c.apiKey == "" {
		return false, fmt.Errorf("HOTPEPPER_API_KEY が設定されていません")
	}

	trimmedLargeAreaCode := strings.TrimSpace(largeAreaCode)
	trimmedSmallAreaCode := strings.TrimSpace(smallAreaCode)
	if trimmedLargeAreaCode == "" || trimmedSmallAreaCode == "" {
		return false, nil
	}

	smallAreas, err := c.FetchSmallAreasByLargeArea(ctx, trimmedLargeAreaCode)
	if err != nil {
		return false, err
	}

	for _, area := range smallAreas {
		if strings.EqualFold(strings.TrimSpace(area.Code), trimmedSmallAreaCode) {
			return true, nil
		}
	}

	return false, nil
}

// FetchSmallAreasByLargeArea はlarge_area配下のsmall_area一覧を返します。
func (c *Client) FetchSmallAreasByLargeArea(ctx context.Context, largeAreaCode string) ([]SmallArea, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("HOTPEPPER_API_KEY が設定されていません")
	}

	trimmedLargeAreaCode := strings.TrimSpace(largeAreaCode)
	if trimmedLargeAreaCode == "" {
		return nil, nil
	}

	middleAreas, err := c.fetchMiddleAreaCodesByLargeArea(ctx, trimmedLargeAreaCode)
	if err != nil {
		return nil, err
	}

	results := make([]SmallArea, 0, 64)
	for _, middleAreaCode := range middleAreas {
		smallAreas, fetchErr := c.fetchSmallAreaCodesByMiddleArea(ctx, trimmedLargeAreaCode, middleAreaCode)
		if fetchErr != nil {
			return nil, fetchErr
		}
		results = append(results, smallAreas...)
	}

	return results, nil
}

func (c *Client) fetchMiddleAreaCodesByLargeArea(ctx context.Context, largeAreaCode string) ([]string, error) {
	values := url.Values{}
	values.Set("key", c.apiKey)
	values.Set("format", "json")
	values.Set("large_area", strings.TrimSpace(largeAreaCode))
	values.Set("count", "100")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, middleAreaEndpoint+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HotPepper APIへの接続に失敗しました: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("HotPepper APIエラー: %s", string(body))
	}

	var payload struct {
		Results struct {
			Error []struct {
				Message string `json:"message"`
			} `json:"error"`
			MiddleArea []struct {
				Code string `json:"code"`
			} `json:"middle_area"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if len(payload.Results.Error) > 0 {
		return nil, fmt.Errorf("HotPepper APIエラー: %s", payload.Results.Error[0].Message)
	}

	codes := make([]string, 0, len(payload.Results.MiddleArea))
	for _, area := range payload.Results.MiddleArea {
		trimmed := strings.TrimSpace(area.Code)
		if trimmed != "" {
			codes = append(codes, trimmed)
		}
	}

	return codes, nil
}

func (c *Client) fetchSmallAreaCodesByMiddleArea(ctx context.Context, largeAreaCode string, middleAreaCode string) ([]SmallArea, error) {
	values := url.Values{}
	values.Set("key", c.apiKey)
	values.Set("format", "json")
	values.Set("middle_area", strings.TrimSpace(middleAreaCode))
	values.Set("count", "100")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, smallAreaEndpoint+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HotPepper APIへの接続に失敗しました: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("HotPepper APIエラー: %s", string(body))
	}

	var payload struct {
		Results struct {
			Error []struct {
				Message string `json:"message"`
			} `json:"error"`
			SmallArea []struct {
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"small_area"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if len(payload.Results.Error) > 0 {
		return nil, fmt.Errorf("HotPepper APIエラー: %s", payload.Results.Error[0].Message)
	}

	codes := make([]SmallArea, 0, len(payload.Results.SmallArea))
	for _, area := range payload.Results.SmallArea {
		code := strings.TrimSpace(area.Code)
		if code != "" {
			codes = append(codes, SmallArea{
				Code:          code,
				Name:          strings.TrimSpace(area.Name),
				MiddleArea:    strings.TrimSpace(middleAreaCode),
				LargeAreaCode: strings.TrimSpace(largeAreaCode),
			})
		}
	}

	return codes, nil
}

func parseJSONValueToString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}

	var valueString string
	if err := json.Unmarshal(raw, &valueString); err == nil {
		return strings.TrimSpace(valueString)
	}

	var valueFloat float64
	if err := json.Unmarshal(raw, &valueFloat); err == nil {
		return strconv.FormatFloat(valueFloat, 'f', -1, 64)
	}

	return strings.Trim(string(raw), "\"")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}

	return ""
}
