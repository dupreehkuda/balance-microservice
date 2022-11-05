package interfaces

import "net/http"

type Handlers interface {
}

type Stored interface {
}

type Actions interface {
}

type Middleware interface {
	CheckCompression(next http.Handler) http.Handler
	WriteCompressed(next http.Handler) http.Handler
}
