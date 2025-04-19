package server

import (
	"github.com/ahnaftahmid39/http-from-tcp/internal/request"
	"github.com/ahnaftahmid39/http-from-tcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)
