package handlers

import "net/http"

type NotFoundHandler struct {
	Context *Context
}

func NewNotFoundHandler(context *Context) http.Handler {
	return &NotFoundHandler{Context: context}
}

func (h *NotFoundHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	h.Context.Renderer.RenderTemplate("404.html", nil, writer)
}
