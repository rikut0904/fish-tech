package router

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	// headerReferer はRefererヘッダー名です。
	headerReferer = "Referer"
)

// RequireAdminOrigin は管理画面オリジンのみを許可するミドルウェアです。
func RequireAdminOrigin(allowedOrigins []string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[normalizeOrigin(origin)] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := normalizeOrigin(c.Request().Header.Get(echo.HeaderOrigin))
			if origin == "" {
				origin = extractOriginFromReferer(c.Request().Header.Get(headerReferer))
			}

			if origin == "" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "管理画面からのみ利用できます"})
			}

			if _, ok := allowed[origin]; !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "管理画面からのみ利用できます"})
			}

			return next(c)
		}
	}
}

func parseAllowedOrigins(raw string, defaultOrigin string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{defaultOrigin}
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		normalized := normalizeOrigin(part)
		if normalized == "" {
			continue
		}
		result = append(result, normalized)
	}

	if len(result) == 0 {
		return []string{defaultOrigin}
	}

	return result
}

func extractOriginFromReferer(referer string) string {
	referer = strings.TrimSpace(referer)
	if referer == "" {
		return ""
	}

	parsed, err := url.Parse(referer)
	if err != nil {
		return ""
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}

	return normalizeOrigin(parsed.Scheme + "://" + parsed.Host)
}

func normalizeOrigin(origin string) string {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return ""
	}

	origin = strings.TrimRight(origin, "/")
	return strings.ToLower(origin)
}
