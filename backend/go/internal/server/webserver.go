package server

import (
	"cerberus-examples/internal/routes"
	"cerberus-examples/internal/services/jwtutils"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
)

type WebServer interface {
	Start()
}

type webServer struct {
	context       context.Context
	port          string
	jwtSecret     string
	publicRoutes  []routes.Routable
	privateRoutes []routes.Routable
}

func NewWebServer(context context.Context, port string, jwtSecret string, publicRoutes []routes.Routable, privateRoutes []routes.Routable) WebServer {
	return &webServer{
		context:       context,
		port:          port,
		jwtSecret:     jwtSecret,
		publicRoutes:  publicRoutes,
		privateRoutes: privateRoutes,
	}
}

func (s *webServer) Start() {
	router := gin.Default()
	applyCors(router)

	public := router.Group("/")
	api := router.Group("/api")
	api.Use(s.JWTAuthRequired)

	for _, route := range s.publicRoutes {
		route.RegisterRoutes(public)
	}

	for _, route := range s.privateRoutes {
		route.RegisterRoutes(api)
	}

	fmt.Println("Listening on port " + s.port)

	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	select {
	case <-s.context.Done():
		log.Println(s.port + " context done.")
	}
	log.Println(s.port + " Server exiting")
}

func (s *webServer) JWTAuthRequired(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userId, accountId, err := s.extractSubjectAndToken(token)
	if err != nil || userId == "" || accountId == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Set userId for route handlers
	c.Set("userId", userId)
	c.Set("accountId", accountId)

	c.Next()
}

func (s *webServer) extractSubjectAndToken(bearer string) (string, string, error) {
	if bearer == "" {
		return "", "", nil
	}
	subject, err := jwtutils.ExtractToken(bearer, s.jwtSecret, func(token *jwt.Token) interface{} {
		claims := token.Claims.(jwt.MapClaims)
		return claims["sub"]
	})
	if err != nil {
		return "", "", err
	}

	if subject == nil {
		return "", "", fmt.Errorf("empty subject in claims")
	}

	accountId, err := jwtutils.ExtractToken(bearer, s.jwtSecret, func(token *jwt.Token) interface{} {
		claims := token.Claims.(jwt.MapClaims)
		extraClaims, ok := claims[subject.(string)].(map[string]interface{})
		if !ok {
			return nil
		}
		return extraClaims["accountId"]
	})
	if err != nil {
		return "", "", err
	}

	if accountId == nil {
		return "", "", fmt.Errorf("empty accountId in extra claims")
	}

	return subject.(string), accountId.(string), nil
}

func applyCors(r *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	//hot reload CORS
	corsConfig.AllowOrigins = []string{"http://localhost:3001"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization", "CerberusAccessToken", "CerberusRefreshToken"}
	r.Use(cors.New(corsConfig))
}
