package portfolio

import (
	"log"

	"github.com/Jojojojodr/portfolio/internal/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunServer(port string) {
	log.Println("Starting Server")
	port = ":" + port

	log.Println("Running Server")

	svr := gin.Default()
	err := svr.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalf("Could not set trusted proxies: %v", err)
	}

	routers.FrontendRouter(svr)
	routers.V1Router(svr)

	log.Println("Server is running at http://localhost" + port)
	svr.Use(cors.Default())
	svr.Run(port)
}
