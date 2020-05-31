package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stekostas/fwd-dog/cache"
	"github.com/stekostas/fwd-dog/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
)

type RedirectHandler struct {
	CacheAdapter cache.Adapter
}

func NewRedirectHandler(cacheAdapter cache.Adapter) *RedirectHandler {
	return &RedirectHandler{
		CacheAdapter: cacheAdapter,
	}
}

func (r *RedirectHandler) Redirect(c *gin.Context) {
	key := c.Request.RequestURI[1:]

	if matched, _ := regexp.MatchString("[^a-zA-Z0-9.]", key); matched {
		r.notFound(c)
		return
	}

	jsonLink, err := r.CacheAdapter.Get(key)

	if err != nil {
		r.notFound(c)
		return
	}

	link := &models.Link{}
	err = json.Unmarshal([]byte(jsonLink), link)

	if err != nil {
		panic(err)
	}

	if len(link.Password) > 0 {
		r.handleUnlock(c, key, link)
		return
	}

	r.handleRedirect(c, key, link)
}

func (r *RedirectHandler) handleUnlock(c *gin.Context, key string, link *models.Link) {
	if c.Request.Method != http.MethodPost {
		c.HTML(http.StatusOK, "redirect.gohtml", gin.H{})
		return
	}

	password, ok := c.GetPostForm("password")

	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(link.Password), []byte(password))

	if err != nil {
		c.HTML(http.StatusUnauthorized, "redirect.gohtml", &TemplateContext{
			Message: "Wrong password.",
		})
		return
	}

	r.handleRedirect(c, key, link)
}

func (r *RedirectHandler) notFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.gohtml", gin.H{})
}

func (r *RedirectHandler) handleRedirect(c *gin.Context, key string, link *models.Link) {
	if strings.HasPrefix(key, ".") {
		err := r.CacheAdapter.Delete(key)

		if err != nil {
			panic(err)
		}
	}

	c.Redirect(http.StatusMovedPermanently, link.TargetUrl)
}
