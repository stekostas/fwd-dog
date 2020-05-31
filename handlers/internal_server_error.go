package handlers

import "net/http"

type InternalServerErrorHandler struct {
	Context *Context
}

func NewInternalServerErrorHandler(context *Context) http.Handler {
	return &InternalServerErrorHandler{Context: context}
}

func (h *InternalServerErrorHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusInternalServerError)
	h.Context.Renderer.RenderTemplate("500.html", nil, writer)
}
