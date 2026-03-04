package router

import (
	"firebase-auth-go/internal/interface/handler"
	"firebase-auth-go/middleware"
	"github.com/gin-gonic/gin"
)

// Ginルーターを構築してCORSおよびルートを設定
func NewRouter() *gin.Engine {
	r := gin.Default()

	// CORS設定
	// ALLOWED_ORIGINSから許可オリジンを読み込む
	// 未設定の場合はローカル開発用のデフォルト値を使用。
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		for _, o := range strings.Split(allowedOriginsEnv, ",") {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	} else {
		allowedOrigins = []string{
			"http://localhost:3000",
			"http://localhost:5173",
		}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))


	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 認証ルート
	auth := r.Group("/auth")
	{
		// 公開エンドポイント（認証不要）
		auth.POST("/signup", handler.SignUp)
		auth.POST("/signin", handler.SignIn)
		auth.POST("/send-verification", handler.SendVerificationEmail)

		// 保護されたエンドポイント（Firebase IDトークン必須）
		protected := auth.Group("")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/me", handler.GetMe)
			protected.DELETE("/me", handler.DeleteAccount)
		}
	}

	// 魚図鑑ルート
	fish := r.Group("/fish")
	{
	// handler. を付けて呼び出す
        fish.GET("", handler.GetFishList)
        fish.GET("/:id", handler.GetFishDetail)
        fish.POST("/:id/like", middleware.AuthRequired(), handler.ToggleLike)
    }

	// 開発用ルート
	dev := r.Group("/dev")
	{
		dev.POST("/seed", handler.SeedData)
	}

	return r
}