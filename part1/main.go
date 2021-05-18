package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	bigCache "github.com/allegro/bigcache/v2"
	"github.com/daheige/bigcache-demo/utils"
	"github.com/daheige/tigago/gpprof"
	"github.com/daheige/tigago/gutils"
	"github.com/daheige/tigago/monitor"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port  = 1336
	cache *bigCache.BigCache
)

// H map简写
type H map[string]interface{}

// ApiResult api result
type ApiResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func InitCache() {
	// 底层默认1024个Shards 分片
	cache, _ = bigCache.NewBigCache(bigCache.DefaultConfig(10 * time.Minute))
}

func init() {
	InitCache()

	// 添加prometheus性能监控指标
	prometheus.MustRegister(monitor.WebRequestTotal)
	prometheus.MustRegister(monitor.WebRequestDuration)

	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

	// 性能监控的端口port+1000,只能在内网访问
	httpMux := gpprof.New()

	// 添加prometheus metrics处理器
	httpMux.Handle("/metrics", promhttp.Handler())
	gpprof.Run(httpMux, port+1000)
}

func main() {
	log.Println("big cache demo...")

	log.Printf("Server running on port:%d/", port)

	// register mux router
	router := RouteHandler()

	// walk router
	walkRouter(router)

	// create http services
	server := &http.Server{
		// Handler: http.TimeoutHandler(router, time.Second*6, `{code:503,"message":"services timeout"}`),
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// run http services in goroutine
	go func() {
		defer utils.Recover()

		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Println("services listen error:", err)
				return
			}

			log.Println("services will exit...")
		}
	}()

	// graceful exit
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recv signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	<-ch

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	go server.Shutdown(ctx)
	<-ctx.Done()

	log.Println("services shutdown success")
}

// RouteHandler router handler
func RouteHandler() *mux.Router {
	r := mux.NewRouter()

	r.StrictSlash(true)

	// install access log and recover handler
	// r.Use(AccessLog, RecoverHandler)

	r.Use(RecoverHandler)

	// not found handler
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	// Index Route
	r.HandleFunc("/", index)
	r.HandleFunc("/set-data", setData)
	r.HandleFunc("/get-data", getData)

	return r
}

func walkRouter(r *mux.Router) {
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}

		var queriesTemplates []string
		queriesTemplates, err = route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}

		var queriesRegexps []string
		queriesRegexps, err = route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}

		var methods []string
		methods, err = route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}

		return nil
	})

	if err != nil {
		fmt.Println("router walk error:", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("big cache test"))
}

func setData(w http.ResponseWriter, r *http.Request) {
	setCache("my-test", 123)
}

func getData(w http.ResponseWriter, r *http.Request) {
	key := "my-test"
	b, err := getCache(key)
	if err == bigCache.ErrEntryNotFound {
		log.Println("cache not fond")

		b = []byte("1234")
		setCache(key, b)
	}

	if len(b) == 0 {
		log.Println("cache is empty")
	}

	res := ApiResult{
		Code:    0,
		Message: "ok",
		Data: H{
			"value": string(b),
		},
	}

	jsonBytes, _ := json.Marshal(res)
	w.Write(jsonBytes)
}

func getCache(key string) ([]byte, error) {
	return cache.Get(key)
}

func setCache(key string, value interface{}) error {
	if v, ok := value.([]byte); ok {
		cache.Set(key, v)
		return nil
	}

	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("set key:%s cache error:%s", key, err.Error())
	}

	cache.Set(key, b)

	return nil
}

// HandlerRecover catch services recover
func HandlerRecover() {
	if err := recover(); err != nil {
		log.Println("exec panic", map[string]interface{}{
			"error":       err,
			"error_trace": string(debug.Stack()),
		})
	}
}

// NotFoundHandler not found api router
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("this page not found"))
}

// RecoverHandler recover handler
func RecoverHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("exec panic error", map[string]interface{}{
					"trace_error": string(debug.Stack()),
				})

				// services error
				http.Error(w, "services error!", http.StatusInternalServerError)
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func AccessLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		log.Println("exec begin", nil)
		log.Println("request before")
		log.Println("request uri: ", r.RequestURI)

		// x-request-id
		reqId := r.Header.Get("x-request-id")
		if reqId == "" {
			reqId = gutils.Uuid()
		}

		// log.Println("log_id: ", reqId)
		r = utils.ContextSet(r, "log_id", reqId)
		r = utils.ContextSet(r, "client_ip", r.RemoteAddr)
		r = utils.ContextSet(r, "request_method", r.Method)
		r = utils.ContextSet(r, "request_uri", r.RequestURI)
		r = utils.ContextSet(r, "user_agent", r.Header.Get("User-Agent"))

		h.ServeHTTP(w, r)

		log.Println("exec end", map[string]interface{}{
			"exec_time": time.Now().Sub(t).Seconds(),
		})

	})
}
