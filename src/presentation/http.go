package rest

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

type ServerHttp interface {
	StartServer(port string) error
	GetGinRouterGroup(relativePath string) *gin.RouterGroup
}

type serverHttpGin struct {
	ginEngine *gin.Engine
}

func (ref *serverHttpGin) StartServer(port string) error {
	var g errgroup.Group

	g.Go(func() error {
		address := fmt.Sprintf(":%s", port)
		err := endless.ListenAndServe(address, ref.ginEngine)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func NewServerHttpGin() ServerHttp {
	r := gin.New()
	r.Use(gin.Recovery())

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	return &serverHttpGin{
		ginEngine: r,
	}
}

func (ref *serverHttpGin) GetGinRouterGroup(relativePath string) *gin.RouterGroup {
	ref.ginEngine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-LimitType"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	return ref.ginEngine.Group(relativePath)
}

func RegisterInfraApi(ginRouterGroup *gin.RouterGroup, diagnosticMode bool) {
	if diagnosticMode {
		pprof.RouteRegister(ginRouterGroup, "pprof")
	}
	ginRouterGroup.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}
