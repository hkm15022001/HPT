package httpServer

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hkm15022001/Supply-Chain-Event-Management/api/middleware"
	"github.com/hkm15022001/Supply-Chain-Event-Management/api/router"
	"github.com/hkm15022001/Supply-Chain-Event-Management/internal/handler"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

// RunServer will start 2 server for app and web
func RunServer(store *handler.Store) {
	// gin.SetMode(gin.ReleaseMode)
	// export GIN_MODE=debug

	webServer := &http.Server{
		Addr:         os.Getenv("WEB_PORT"),
		Handler:      webRouter(store),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := webServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

func webRouter(store *handler.Store) http.Handler {
	e := gin.Default()

	e.Static("/api/images", os.Getenv("IMAGE_FILE_PATH"))
	e.Static("/api/qrcode", os.Getenv("QR_CODE_FILE_PATH"))
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	e.MaxMultipartMemory = 8 << 20 // 8 MiB

	api := e.Group("/api")
	// Active web auth
	if os.Getenv("RUN_WEB_AUTH") == "yes" {
		api.Use(middleware.ValidateWebSession())
	}
	api.Use(func(c *gin.Context) {
		c.Set("store", store)
		c.Next()
	})
	router.WebUserRoutes(api, store)
	return e
}
