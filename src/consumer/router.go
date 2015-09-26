package main

import (
    "fmt"
    "log"
    "time"
    "runtime"
    "net/http"

    "github.com/gorilla/mux"

    "consumer/handler"
)

type Route struct {
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

var routes = []Route{
    Route{"GET", "/api/v1/{pid}/ping", handler.Ping},

    Route{"POST", "/api/v1/{pid}/requests",     handler.PostRequest   },
    Route{"POST", "/api/v1/{pid}/request-logs", handler.PostRequestLog},

    Route{"GET", "/api/v1/{pid}/services",      handler.GetServices    },
    Route{"GET", "/api/v1/{pid}/service-chain", handler.GetServiceChain},
    Route{"GET", "/api/v1/{pid}/request-types", handler.GetRequestTypes},

    Route{"GET", "/api/v1/{pid}/invoke-chains",                                         handler.GetAllInvokeChains        },
    Route{"GET", "/api/v1/{pid}/invoke-chains/{service}/{category}",                    handler.GetInvokeChains           },
    Route{"GET", "/api/v1/{pid}/invoke-chains/{service}/{category}/{id}",               handler.GetInvokeChain            },
    Route{"GET", "/api/v1/{pid}/invoke-chains/{service}/{category}/{id}/root-requests", handler.GetInvokeChainRootRequests},

    Route{"GET", "/api/v1/{pid}/request-logs",        handler.GetAllRequestLogs },
    Route{"GET", "/api/v1/{pid}/request-logs/{uuid}", handler.GetSomeRequestLogs},

    Route{"GET", "/api/v1/{pid}/requests",                      handler.GetRequests           },
    Route{"GET", "/api/v1/{pid}/requests/{uuid}",               handler.GetRequest            },
    Route{"GET", "/api/v1/{pid}/requests/{uuid}/invoke-chain",  handler.GetRequestInvokeChain },
    Route{"GET", "/api/v1/{pid}/requests/{uuid}/group",         handler.GetRequestRequestGroup},
    Route{"GET", "/api/v1/{pid}/requests/{uuid}/root-request",  handler.GetRequestRootRequest },
    Route{"GET", "/api/v1/{pid}/requests/{uuid}/parent",        handler.GetRequestParent      },
    Route{"GET", "/api/v1/{pid}/requests/{uuid}/children",      handler.GetRequestChildren    },
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
