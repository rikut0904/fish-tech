package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"bytes"
	"os"


	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var authClient *auth.Client

func init() {
	godotenv.Load()
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Firebase初期化失敗: %v", err)
	}
	authClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Authクライアント取得失敗: %v", err)
	}
}

// CORS設定用ミドルウェア
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")//本番環境では実際のフロントエンドURLに変更
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	json.NewDecoder(r.Body).Decode(&req)

	// 1. ユーザー作成
	u, err := authClient.CreateUser(context.Background(), (&auth.UserToCreate{}).Email(req.Email).Password(req.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. メールリンク生成
	link, _ := authClient.EmailVerificationLinkWithSettings(context.Background(), req.Email, &auth.ActionCodeSettings{
		//URL: "http://localhost:3000/login", // 確認後のリダイレクト先
		HandleCodeInApp: true,
	})

	// 3. メール送信 (SendGrid)
	sendEmail(req.Email, link)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "確認メールを送信しました", "uid": u.UID})
}


func sendEmail(toEmail, link string) {  //firebase authenticationに対応させる
	from := mail.NewEmail("Auth System", os.Getenv("FROM_EMAIL"))
	to := mail.NewEmail("User", toEmail)
	content := fmt.Sprintf("<p>確認リンク: <a href='%s'>こちらをクリック</a></p>", link)
	message := mail.NewSingleEmail(from, "メール確認", to, "", content)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	client.Send(message)
}



func main() {
	http.HandleFunc("/register", enableCORS(registerHandler))
	http.HandleFunc("/login", enableCORS(loginHandler))
	// 認証が必要なサンプルAPI
	http.HandleFunc("/api/private", enableCORS(authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("これは秘密のデータです"))
	})))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}



type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FirebaseLoginResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
	Registered   bool   `json:"registered"`
}



func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

apiKey := os.Getenv("FIREBASE_API_KEY")
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)

	payload, _ := json.Marshal(map[string]interface{}{
		"email":             req.Email,
		"password":          req.Password,
		"returnSecureToken": true,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		http.Error(w, "Failed to connect to Firebase", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	var fbResp FirebaseLoginResponse
	json.NewDecoder(resp.Body).Decode(&fbResp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fbResp)
}



func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		idToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if idToken == "" {
			http.Error(w, "認証トークンが必要です", http.StatusUnauthorized)
			return
		}
		token, err := authClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			http.Error(w, "無効なトークンです", http.StatusUnauthorized)
			return
		}
		log.Printf("User Verified: %s", token.UID)
		next(w, r)
	}
}