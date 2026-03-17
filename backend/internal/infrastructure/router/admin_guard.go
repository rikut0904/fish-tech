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
func RequireAdminOrigin(allowedOrigins []string, swaggerOrigins []string, swaggerRefererPaths []string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[normalizeOrigin(origin)] = struct{}{}
	}

	allowedSwaggerOrigins := make(map[string]struct{}, len(swaggerOrigins))
	for _, origin := range swaggerOrigins {
		allowedSwaggerOrigins[normalizeOrigin(origin)] = struct{}{}
	}

	allowedSwaggerPaths := make(map[string]struct{}, len(swaggerRefererPaths))
	for _, refererPath := range swaggerRefererPaths {
		normalized := normalizeRefererPath(refererPath)
		if normalized == "" {
			continue
		}
		allowedSwaggerPaths[normalized] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := normalizeOrigin(c.Request().Header.Get(echo.HeaderOrigin))
			referer := c.Request().Header.Get(headerReferer)
			if origin == "" {
				origin = extractOriginFromReferer(referer)
			}

			if origin == "" {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "管理画面からのみ利用できます"})
			}

			if _, ok := allowed[origin]; !ok {
				if !isAllowedSwaggerRequest(origin, referer, allowedSwaggerOrigins, allowedSwaggerPaths) {
					return c.JSON(http.StatusForbidden, map[string]string{"error": "管理画面からのみ利用できます"})
				}
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

func extractPathFromReferer(referer string) string {
	referer = strings.TrimSpace(referer)
	if referer == "" {
		return ""
	}

	parsed, err := url.Parse(referer)
	if err != nil {
		return ""
	}

	return normalizeRefererPath(parsed.Path)
}

func isAllowedSwaggerRequest(origin string, referer string, allowedOrigins map[string]struct{}, allowedPaths map[string]struct{}) bool {
	if _, ok := allowedOrigins[origin]; !ok {
		return false
	}

	refererPath := extractPathFromReferer(referer)
	if refererPath == "" {
		return false
	}

	_, ok := allowedPaths[refererPath]
	return ok
}

func normalizeRefererPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = strings.TrimRight(path, "/")
	if path == "" {
		return "/"
	}
	return path
}

func normalizeOrigin(origin string) string {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return ""
	}

	origin = strings.TrimRight(origin, "/")
	return strings.ToLower(origin)
}
