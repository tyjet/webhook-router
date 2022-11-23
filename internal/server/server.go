package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tyjet/webhook-test/internal/config"
)

const httpXHubSignature256 = "HTTP_X_HUB_SIGNATURE_256"

type Server struct {
	cfg config.Config
}

func NewServer(cfg config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Register() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		fmt.Println("[DBG] root route")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200 - OK"))
	})
	http.Handle("/payload", s.verifyWebhookSignature(http.HandlerFunc(s.postWebhook)))
}

func (s *Server) Start() {
	http.ListenAndServe("localhost:8128", nil)
}

func (s *Server) postWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body {error: %s}\n", err.Error())
		http.Error(w, "could not ready body", http.StatusInternalServerError)
		return
	}
	fmt.Println("body", body)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) verifyWebhookSignature(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("could not read body {error: %s}\n", err.Error())
			http.Error(w, "could not ready body", http.StatusInternalServerError)
			return
		}

		ghSignature := r.Header.Get(httpXHubSignature256)
		if ghSignature == "" {
			msg := "missing " + httpXHubSignature256 + " header"
			fmt.Println(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		ghSignatureFields := strings.Split(ghSignature, "=")
		if len(ghSignatureFields) != 2 {
			msg := "malformed GitHub signature " + ghSignature
			fmt.Println(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		method := ghSignatureFields[0]
		if method != "sha256" {
			msg := "expected SHA256, but was " + method
			fmt.Println(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		ghDigest := ghSignatureFields[1]
		mac := hmac.New(sha256.New, []byte(s.cfg.GHWebhookSecret))
		mac.Write(body)
		digest := mac.Sum(nil)
		if !hmac.Equal(digest, []byte(ghDigest)) {
			msg := fmt.Sprintf("GitHub signature does not match expected signature {expected: %s, actual: %s}", digest, ghDigest)
			fmt.Println(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		h.ServeHTTP(w, r)
	})
}
