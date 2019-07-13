package logger

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/mgutz/ansi"
)

// Data contains data to be logged
type Data map[string]interface{}

var namespace = "service-namespace"

var version = ""

var logMutex sync.Mutex

var redisClient *redis.Client

// SetNamespace sets the logger namespace
func SetNamespace(ns string) {
	namespace = ns
}

// SetVersion sets the app version
func SetVersion(v string) {
	version = v
}

// SetStreamLogger takes a redis client a sets it as the stream logger
func SetStreamLogger(cli *redis.Client) {
	redisClient = cli
}

// Handler represents an http handler
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		rc := &responseCapture{w, 200}

		h.ServeHTTP(rc, req)

		elapsed := time.Since(start)
		end := time.Now()

		data := Data{
			"start":    start.Format("02-01-2006 15:04:05.999999"),
			"end":      end.Format("02-01-2006 15:04:05.999999"),
			"duration": elapsed,
			"status":   rc.statusCode,
			"method":   req.Method,
			"path":     req.URL.Path,
		}
		if len(req.URL.RawQuery) > 0 {
			data["query"] = req.URL.Query()
		}

		Request(data)

	})

}

type responseCapture struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseCapture) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseCapture) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (r *responseCapture) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := r.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, errors.New("log: response does not implement http.Hijacker")
}

// Error logs the error to stdout
func Error(err error, data ...Data) {
	jsonlog, _ := strconv.ParseBool(os.Getenv("JSON_LOG"))

	if err == nil {
		return
	}
	var d Data
	if len(data) > 0 {
		d = data[0]
	}

	if d == nil {
		d = make(Data)
	}

	d["error"] = err
	d["message"] = err.Error()
	d["type"] = "error"
	d["created_at"] = time.Now().Format("02-01-2006 15:04:05.999999")
	d["namespace"] = namespace

	logjson(d, jsonlog)

	consolelog(err.Error(), d, ansi.Red)
}

// Request logs a request
func Request(data ...Data) {
	jsonlog, _ := strconv.ParseBool(os.Getenv("JSON_LOG"))

	message := "request"

	var d Data
	if len(data) > 0 {
		d = data[0]
	}

	if d == nil {
		d = make(Data)
	}

	d["message"] = message
	d["type"] = "request"
	d["created_at"] = time.Now().Format("02-01-2006 15:04:05.999999")
	d["namespace"] = namespace

	logjson(d, jsonlog)

	consolelog(message, d, ansi.Cyan)
}

// OK logs the message to stdout (green)
func OK(message string, data ...Data) {
	jsonlog, _ := strconv.ParseBool(os.Getenv("JSON_LOG"))

	var d Data
	if len(data) > 0 {
		d = data[0]
	}

	if d == nil {
		d = make(Data)
	}

	d["message"] = message
	d["type"] = "ok"
	d["created_at"] = time.Now().Format("02-01-2006 15:04:05.999999")
	d["namespace"] = namespace

	logjson(d, jsonlog)

	consolelog(message, d, ansi.Green)
}

// Trace logs the message to stdout (blue)
func Trace(message string, data ...Data) {
	jsonlog, _ := strconv.ParseBool(os.Getenv("JSON_LOG"))

	var d Data
	if len(data) > 0 {
		d = data[0]
	}

	if d == nil {
		d = make(Data)
	}

	d["message"] = message
	d["type"] = "trace"
	d["created_at"] = time.Now().Format("02-01-2006 15:04:05.999999")
	d["namespace"] = namespace

	logjson(d, jsonlog)
	consolelog(message, d, ansi.Blue)
}

func logjson(data Data, jsonlog bool) {

	b, _ := json.Marshal(data)

	if redisClient != nil {
		go redisClient.SAdd("log-stream", string(b))
	}

	if jsonlog {
		fmt.Fprintf(os.Stdout, "%s\n", b)
	}
}

func consolelog(msg string, data Data, colour string) {
	logMutex.Lock()
	defer logMutex.Unlock()
	fmt.Fprintf(os.Stdout, "%s%s %s %s\n", colour, data["created_at"], msg, ansi.DefaultFG)
	for k, v := range data {
		fmt.Fprintf(os.Stdout, "\t -- %s: %+v\n", k, v)
	}
}
