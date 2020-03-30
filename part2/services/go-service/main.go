package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

/** Added for part2 */
  "strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/mmedum/prometheus-workshop/go-service/handlers/health"

/** Added for part2 */
  "github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var responseCodes [3]int

/** Added for part2 */
var (
  pongCount = prometheus.NewCounterVec(
    prometheus.CounterOpts{
      Name: "ping_total_number_of_requests",
      Help: "Number of ping requests.",
    },
    []string{"status"},
  )
)

/** Added for part2 */
func initPrometheusMetric() {
  prometheus.MustRegister(pongCount)
}


func pong(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "pong"

	rand.Seed(time.Now().Unix())

	responseCode := responseCodes[rand.Intn(len(responseCodes))]

/** Added for part2 */
  pongCount.WithLabelValues(strconv.Itoa(responseCode)).Inc()


	render.Status(r, responseCode)

	render.JSON(w, r, response)
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.RequestID,
	)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)

	router.Mount("/metrics", promhttp.Handler())
	router.Mount("/health", health.Routes())

	router.Route("/v1", func(r chi.Router) {
		r.Get("/ping", pong)
	})

	return router
}

func initResponseCodes() {
	responseCodes[0] = 200
	responseCodes[1] = 500
	responseCodes[2] = 503
}

func main() {
	initResponseCodes()

/** Added for part2 */
  initPrometheusMetric()

	router := Routes()
	port := 80

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
	log.Println("Running")
}
