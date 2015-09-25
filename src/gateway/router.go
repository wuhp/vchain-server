package main

import (
    "fmt"
    "log"
    "time"
    "runtime"
    "net/http"

    "github.com/gorilla/mux"

    "gateway/handler"
)

type Route struct {
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

var routes = []Route{
    Route{"GET", "/api/v1/ping", handler.Ping},

    Route{"GET",    "/api/v1/instances",      handler.QueryInstance },
    Route{"POST",   "/api/v1/instances",      handler.CreateInstance},
    Route{"GET",    "/api/v1/instances/{id}", handler.GetInstance   },
    Route{"PUT",    "/api/v1/instances/{id}", handler.UpdateInstance},
    Route{"DELETE", "/api/v1/instances/{id}", handler.DeleteInstance},

    Route{"GET",  "/api/v1/mysql",       handler.QueryMapping},
    Route{"GET",  "/api/v1/mysql/{pid}", handler.GetMapping  },
    Route{"POST", "/api/v1/provision",   handler.Provision   },
    Route{"POST", "/api/v1/unprovision", handler.Unprovision },
}

type InnerResponseWriter struct {
    StatusCode int
    isSet      bool
    http.ResponseWriter
}

func (i *InnerResponseWriter) WriteHeader(status int) {
    if !i.isSet {
        i.StatusCode = status
        i.isSet = true
    }

    i.ResponseWriter.WriteHeader(status)
}

func (i *InnerResponseWriter) Write(b []byte) (int, error) {
    i.isSet = true
    return i.ResponseWriter.Write(b)
}

func wrapper(inner http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        s := time.Now()
        wr := &InnerResponseWriter{
            StatusCode:     200,
            isSet:          false,
            ResponseWriter: w,
        }

        defer func() {
            if err := recover(); err != nil {
                wr.WriteHeader(http.StatusInternalServerError)
                buf := make([]byte, 1<<16)
                stackSize := runtime.Stack(buf, true)
                log.Printf("Panic: %v\n%s\n", err, string(buf[0:stackSize]))
                fmt.Fprintf(w, fmt.Sprintln(err))
            }

            d := time.Now().Sub(s)
            log.Printf("%s %s %d %s\n", r.Method, r.RequestURI, wr.StatusCode, d.String())
        }()

        inner.ServeHTTP(wr, r)
    })
}

func NewRouter() *mux.Router {
    router := mux.NewRouter()
    for _, route := range routes {
        router.Methods(route.Method).Path(route.Pattern).HandlerFunc(wrapper(route.HandlerFunc))
    }

    return router
}
