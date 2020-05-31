package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type HomepageHandler struct {
	Context *Context
}

type TemplateContext struct {
	Success bool
	Message string
	Link    string
}

func NewHomepageHandler(context *Context) http.Handler {
	return &HomepageHandler{Context: context}
}

func (h *HomepageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	templateName := "index.html"

	if request.Method == http.MethodPost {
		h.handleFormSubmission(writer, request, templateName)
		return
	}

	h.Context.Renderer.RenderTemplate(templateName, nil, writer)
}

func (h *HomepageHandler) handleFormSubmission(writer http.ResponseWriter, request *http.Request, templateName string) {
	targetUrl := request.FormValue("url")
	ttl := request.FormValue("ttl")
	validUrl := h.isValidUrl(writer, targetUrl, templateName)
	validTtl := h.isValidTtl(writer, ttl, templateName)

	if !validUrl || !validTtl {
		return
	}

	h.generateKey(writer, request, targetUrl, ttl, templateName)
}

func (h *HomepageHandler) isValidUrl(writer http.ResponseWriter, targetUrl string, templateName string) bool {
	_, err := url.ParseRequestURI(targetUrl)

	if err == nil {
		return true
	}

	writer.WriteHeader(400)
	h.Context.Renderer.RenderTemplate(templateName, &TemplateContext{Success: false, Message: "Please enter a valid URL."}, writer)
	return false
}

func (h *HomepageHandler) isValidTtl(writer http.ResponseWriter, ttl string, templateName string) bool {
	ttlInt, _ := strconv.Atoi(ttl)

	for _, duration := range h.Context.TtlOptions {
		if time.Second*time.Duration(ttlInt) == duration {
			return true
		}
	}

	writer.WriteHeader(400)
	h.Context.Renderer.RenderTemplate(templateName, &TemplateContext{Success: false, Message: "Please select a valid expiration time."}, writer)
	return false
}

func (h *HomepageHandler) generateKey(writer http.ResponseWriter, request *http.Request, targetUrl string, ttl string, templateName string) {
	key := ""
	length := 1
	limit := 6
	ttlInt, _ := strconv.Atoi(ttl)
	keyHash := h.getKeyHash(targetUrl)

	for {
		key = keyHash[:length]

		ok, err := h.Context.CacheAdapter.SetOrFail(key, targetUrl, time.Second*time.Duration(ttlInt))

		if ok {
			break
		}

		if length >= limit || err != nil {
			internalServerErrorHandler := NewInternalServerErrorHandler(h.Context)
			internalServerErrorHandler.ServeHTTP(writer, request)
			return
		}

		length++
	}

	final := fmt.Sprintf("%s/%s", request.Host, key)

	writer.Header().Add("X-Fwd-Key", key)
	h.Context.Renderer.RenderTemplate(templateName, &TemplateContext{Success: true, Link: final}, writer)
}

func (h *HomepageHandler) getKeyHash(targetUrl string) string {
	unixNano := time.Now().UnixNano()
	timestamp := strconv.FormatInt(unixNano, 10)
	key := targetUrl + timestamp

	hasher := sha256.New()
	hasher.Write([]byte(key))
	hash := hasher.Sum(nil)
	encoded := base64.URLEncoding.EncodeToString(hash)
	pattern := regexp.MustCompile(`[^a-zA-Z0-9]`)

	return pattern.ReplaceAllString(encoded, "")
}
