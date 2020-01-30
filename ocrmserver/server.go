package ocrmserver

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sergazyyev/crmlibrary/ocrmerrors"
	"github.com/sergazyyev/crmlibrary/ocrmmodel"
	"github.com/sirupsen/logrus"
	"net/http"
)

type BaseServer struct {
	Router              *mux.Router
	JwtKey              []byte
	JwtTokenLiveMinutes int
	Logger              *logrus.Logger
	UseAuthMidd         bool
	AuthIgnorePaths     []string
}

func (server *BaseServer) Start(port int) error {
	if server.UseAuthMidd && (server.JwtKey == nil || server.JwtTokenLiveMinutes == 0) {
		return ocrmerrors.New(ocrmerrors.ARGISNIL, "For authenticateUser middleware JwtKey and JwtTokenLiveMinutes must not be nil", "")
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", port), server)
}

func (server *BaseServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.Router.ServeHTTP(w, r)
}

func (server *BaseServer) ConfigureRouterMiddleware() {
	server.Router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodOptions}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Content-Disposition"}),
		handlers.ExposedHeaders([]string{"Authorization", "Set-Cookie", "Content-Disposition"})))
	//Middleware
	server.Router.Use(server.setRequestId)
	server.Router.Use(server.loggingRequests)
	server.Router.Use(server.authenticateUser)
}

func (server *BaseServer) RespondJson(w http.ResponseWriter, status int, data interface{}) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}

func (server *BaseServer) Error(w http.ResponseWriter, status int, err error) {
	server.respondJson(w, status, &ocrmmodel.SimpleResponse{Code: ocrmmodel.SimpleErrCode, Message: err.Error()})
}

func (server *BaseServer) respondJson(w http.ResponseWriter, status int, data interface{}) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(status)
	}
}

func (server *BaseServer) RespondByte(w http.ResponseWriter, status int, data []byte, headers map[string]string) {
	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}
	if data != nil && len(data) > 0 {
		w.WriteHeader(status)
		if _, err := w.Write(data); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(status)
	}
}
