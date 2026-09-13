package portfolio

import (
	"log"

	"github.com/Jojojojodr/portfolio/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
	Port   string
}

func (s *Server) Start() {
	log.Println("Starting Server on port " + s.Port)
	log.Println("Server is running at http://localhost:" + s.Port)
	s.Engine.Use(cors.Default())
	if err := s.Engine.Run(":" + s.Port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func NewServer(port string) *Server {
	if config.AppConfig.Gin.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	for _, trustedProxy := range config.AppConfig.Gin.TrustedProxies {
		err := engine.SetTrustedProxies([]string{trustedProxy})
		if err != nil {
			log.Fatalf("Could not set trusted proxies: %v", err)
		}
	}

	return &Server{
		Engine: engine,
		Port:   port,
	}
}
