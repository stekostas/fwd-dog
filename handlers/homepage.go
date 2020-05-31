package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stekostas/fwd-dog/models"
	"github.com/stekostas/fwd-dog/services"
	"net/http"
	"net/url"
	"time"
)

type HomepageHandler struct {
	TtlOptions    map[time.Duration]string
	LinkGenerator *services.LinkGenerator
}

type TemplateContext struct {
	Success    bool
	Message    string
	Link       string
	TtlOptions map[time.Duration]string
}

func NewHomepageHandler(ttlOptions map[time.Duration]string, linkGenerator *services.LinkGenerator) *HomepageHandler {
	return &HomepageHandler{
		TtlOptions:    ttlOptions,
		LinkGenerator: linkGenerator,
	}
}

func (h *HomepageHandler) Get(c *gin.Context) {
	c.HTML(http.StatusOK, "index.gohtml", &TemplateContext{
		TtlOptions: h.TtlOptions,
	})
}

func (h *HomepageHandler) Post(c *gin.Context) {
	form := &models.CreateLinkForm{}
	template := "index.gohtml"

	err := c.ShouldBind(form)
	data, valErr := h.ValidateCreateLinkForm(form)

	if err != nil || valErr != nil {
		c.HTML(http.StatusBadRequest, template, &TemplateContext{
			Success:    false,
			Message:    "Invalid form data.",
			TtlOptions: h.TtlOptions,
		})
		return
	}

	key, genErr := h.LinkGenerator.Generate(data)

	if genErr != nil {
		panic(genErr)
	}

	c.Header("X-Fwd-Key", key)

	c.HTML(http.StatusCreated, template, &TemplateContext{
		Success: true,
		Link:    "https://" + c.Request.Host + "/" + key,
	})
}

func (h *HomepageHandler) ValidateCreateLinkForm(f *models.CreateLinkForm) (*models.CreateLink, error) {
	const checkboxTrueValue = "on"
	data := &models.CreateLink{
		TargetUrl: f.TargetUrl,
		Ttl:       time.Second * time.Duration(f.Ttl),
	}

	_, err := url.ParseRequestURI(data.TargetUrl)

	if err != nil {
		return nil, fmt.Errorf("the URL provided '%s' is not valid", data.TargetUrl)
	}

	if _, ok := h.TtlOptions[data.Ttl]; !ok {
		return nil, fmt.Errorf("ttl provided must be one of %v, %v given", h.TtlOptions, data.Ttl)
	}

	data.SingleUse = f.SingleUse == checkboxTrueValue
	data.PasswordProtected = f.PasswordProtected == checkboxTrueValue

	if data.PasswordProtected && len(f.Password) < 1 {
		return nil, fmt.Errorf("password is empty but the link was requested to be password protected")
	}

	data.Password = f.Password

	return data, nil
}
