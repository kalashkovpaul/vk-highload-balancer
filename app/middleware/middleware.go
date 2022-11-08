package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var HitsCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "hits_counter",
	Help: "Number of hits to the server",
})

func IncCounter(c *gin.Context) {
	HitsCounter.Inc()
}
