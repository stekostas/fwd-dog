package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type LinkRedirectHandler struct {
	Context *Context
}

func NewLinkRedirectHandler(context *Context) http.Handler {
	return &LinkRedirectHandler{Context: context}
}

func (h *LinkRedirectHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key := vars["key"]
	target, err := h.Context.CacheAdapter.Get(key)

	if err != nil {
		notFoundHandler := NewNotFoundHandler(h.Context)
		notFoundHandler.ServeHTTP(writer, request)
		return
	}

	http.Redirect(writer, request, target, http.StatusFound)
}
