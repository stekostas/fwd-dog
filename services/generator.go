package services

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/stekostas/fwd-dog/cache"
	"github.com/stekostas/fwd-dog/models"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strconv"
	"time"
)

type LinkGenerator struct {
	CacheAdapter cache.Adapter
}

func NewLinkGenerator(cacheAdapter cache.Adapter) *LinkGenerator {
	return &LinkGenerator{
		CacheAdapter: cacheAdapter,
	}
}

func (g *LinkGenerator) Generate(data *models.CreateLink) (string, error) {
	if len(data.TargetUrl) < 1 {
		return "", fmt.Errorf("cannot generate link with no target URL")
	}

	key := ""
	length := 1
	limit := 6
	keyHash := g.getKeyHash(data.TargetUrl)
	link := &models.Link{
		TargetUrl: data.TargetUrl,
	}

	if data.PasswordProtected {
		if len(data.Password) == 0 {
			return "", fmt.Errorf("cannot generate password protected link without password")
		}

		link.Password = g.getPasswordHash(data.Password)
	}

	for {
		key = keyHash[:length]

		if data.SingleUse {
			key = "." + key
		}

		jsonData, jsonErr := json.Marshal(link)

		if jsonErr != nil {
			return "", fmt.Errorf("encountered error when encoding link object to JSON: %v", jsonErr.Error())
		}

		ok, err := g.CacheAdapter.SetOrFail(key, jsonData, data.Ttl)

		if ok {
			break
		}

		if length >= limit || err != nil {
			return "", fmt.Errorf("could not generate link after %d tries: %v", limit, err)
		}

		length++
	}

	return key, nil
}

func (g *LinkGenerator) getKeyHash(targetUrl string) string {
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

func (g *LinkGenerator) getPasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	return string(hash)
}
