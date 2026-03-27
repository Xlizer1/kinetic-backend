// @title           Kinetic API
// @version         1.0
// @description     Backend API for Kinetic communication platform
// @termsOfService  http://localhost:8080/terms

// @contact.name   API Support
// @contact.url    http://localhost:8080/support
// @contact.email  support@kinetic.com

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"kinetic-backend/internal/config"
	"kinetic-backend/internal/handlers"
	"kinetic-backend/internal/middleware"
	"kinetic-backend/internal/models"
	"kinetic-backend/internal/realtime"
	"kinetic-backend/internal/repositories"
	"kinetic-backend/internal/services"
	"kinetic-backend/internal/utils"

	docs "kinetic-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Server{},
		&models.Channel{},
		&models.Message{},
		&models.ServerMember{},
		&models.VoiceState{},
		&models.Presence{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	utils.SetJWTSecret(cfg.JWTSecret)

	userRepo := repositories.NewUserRepository(db)
	serverRepo := repositories.NewServerRepository(db)
	channelRepo := repositories.NewChannelRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	voiceRepo := repositories.NewVoiceRepository(db)
	presenceRepo := repositories.NewPresenceRepository(db)

	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	serverService := services.NewServerService(serverRepo)
	channelService := services.NewChannelService(channelRepo)
	messageService := services.NewMessageService(messageRepo)
	wsService := services.NewWsService(userRepo, messageRepo)
	voiceService := services.NewVoiceService(voiceRepo, channelRepo)
	presenceService := services.NewPresenceService(presenceRepo, userRepo)

	hub := realtime.NewHub()
	hub.ServerID = getServerID()
	hub.UserAuth = wsService.AuthenticateToken
	hub.SaveMessage = wsService.SaveMessage

	var sfuService *services.SFUService
	if cfg.LiveKitAPIKey != "" && cfg.LiveKitAPISecret != "" && cfg.LiveKitServerURL != "" {
		sfuService, err = services.NewSFUService(cfg.LiveKitAPIKey, cfg.LiveKitAPISecret, cfg.LiveKitServerURL)
		if err != nil {
			log.Printf("Warning: Failed to initialize SFU service: %v", err)
		} else {
			log.Printf("SFU service initialized: %s", cfg.LiveKitServerURL)
		}
	}

	if cfg.RedisHost != "" {
		redisConfig := &config.RedisConfig{
			Host:     cfg.RedisHost,
			Port:     cfg.RedisPort,
			Password: cfg.RedisPassword,
		}
		redisClient, err := redisConfig.Connect()
		if err != nil {
			log.Printf("Warning: Failed to connect to Redis: %v", err)
		} else {
			log.Printf("Redis connected: %s:%s", cfg.RedisHost, cfg.RedisPort)
			pubSubHub := realtime.NewPubSubHub(redisClient, hub, hub.ServerID)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if err := pubSubHub.Start(ctx); err != nil {
				log.Printf("Warning: Failed to start Redis pubsub: %v", err)
			}
			hub.PubSub = pubSubHub
		}
	}

	hub.JoinVoice = func(channelID, userID uint) error {
		if err := voiceService.Join(channelID, userID); err != nil {
			return err
		}
		return nil
	}
	hub.LeaveVoice = func(channelID, userID uint) error {
		if err := voiceService.Leave(channelID, userID); err != nil {
			return err
		}
		if sfuService != nil {
			sfuService.LeaveVoice(channelID, userID)
		}
		return nil
	}
	hub.GetPresenceList = func() ([]realtime.PresenceUserInfo, error) {
		users, err := presenceService.GetAllOnline()
		if err != nil {
			return nil, err
		}
		result := make([]realtime.PresenceUserInfo, len(users))
		for i, u := range users {
			result[i] = realtime.PresenceUserInfo{
				UserID:   u.UserID,
				Username: u.Username,
				Status:   u.Status,
			}
		}
		return result, nil
	}
	hub.SetPresence = func(userID uint, status string) error {
		return presenceService.SetPresence(userID, status)
	}
	hub.GetVoiceUsers = func(channelID uint) ([]realtime.VoiceUserInfo, error) {
		states, err := voiceService.GetChannelUsers(channelID)
		if err != nil {
			return nil, err
		}
		users := make([]realtime.VoiceUserInfo, len(states))
		for i, s := range states {
			users[i] = realtime.VoiceUserInfo{
				UserID: s.UserID,
			}
		}
		return users, nil
	}

	go hub.Run()

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	serverHandler := handlers.NewServerHandler(serverService)
	channelHandler := handlers.NewChannelHandler(channelService)
	messageHandler := handlers.NewMessageHandler(messageService, hub)
	voiceHandler := handlers.NewVoiceHandler(voiceService, hub)
	wsHandler := handlers.NewWsHandler(hub)

	docs.SwaggerInfo.BasePath = "/api"

	// Configure Swagger host and scheme for production (defaults set for production)
	swaggerHost := os.Getenv("SWAGGER_HOST")
	swaggerScheme := os.Getenv("SWAGGER_SCHEME")

	// If environment variables not set, use production defaults
	if swaggerHost == "" {
		swaggerHost = "api.kinetic.kite-app.online"
	}
	if swaggerScheme == "" {
		swaggerScheme = "https"
	}

	docs.SwaggerInfo.Host = swaggerHost
	docs.SwaggerInfo.Schemes = []string{swaggerScheme}

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	r.Use(middleware.LoggingMiddleware())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/verify-email", authHandler.VerifyEmail)
			auth.POST("/refresh-token", authHandler.RefreshToken)
		}

		users := api.Group("/users")
		{
			users.GET("/:id", userHandler.GetUser)
			users.GET("/@me", middleware.AuthMiddleware(), userHandler.GetMe)
			users.PATCH("/@me", middleware.AuthMiddleware(), userHandler.UpdateMe)
			users.PATCH("/@me/settings", middleware.AuthMiddleware(), userHandler.UpdateSettings)
		}

		servers := api.Group("/servers")
		{
			servers.GET("", middleware.AuthMiddleware(), serverHandler.GetServers)
			servers.POST("", middleware.AuthMiddleware(), serverHandler.CreateServer)
			servers.POST("/join", middleware.AuthMiddleware(), serverHandler.JoinServer)
			servers.GET("/:id", middleware.AuthMiddleware(), serverHandler.GetServer)
			servers.PATCH("/:id", middleware.AuthMiddleware(), serverHandler.UpdateServer)
			servers.DELETE("/:id", middleware.AuthMiddleware(), serverHandler.DeleteServer)
			servers.POST("/:id/leave", middleware.AuthMiddleware(), serverHandler.LeaveServer)
			servers.GET("/:id/channels", middleware.AuthMiddleware(), channelHandler.GetChannels)
		}

		channels := api.Group("/channels")
		{
			channels.POST("", middleware.AuthMiddleware(), channelHandler.CreateChannel)
			channels.GET("/:id", middleware.AuthMiddleware(), channelHandler.GetChannel)
			channels.PATCH("/:id", middleware.AuthMiddleware(), channelHandler.UpdateChannel)
			channels.DELETE("/:id", middleware.AuthMiddleware(), channelHandler.DeleteChannel)
			channels.GET("/:id/messages", middleware.AuthMiddleware(), messageHandler.GetMessages)
			channels.POST("/:id/messages", middleware.AuthMiddleware(), messageHandler.CreateMessage)
			channels.DELETE("/:id/messages/:messageId", middleware.AuthMiddleware(), messageHandler.DeleteMessage)
			channels.POST("/:id/voice/join", middleware.AuthMiddleware(), voiceHandler.JoinVoice)
			channels.POST("/:id/voice/leave", middleware.AuthMiddleware(), voiceHandler.LeaveVoice)
			channels.GET("/:id/voice", middleware.AuthMiddleware(), voiceHandler.GetVoiceUsers)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ws", wsHandler.HandleWebSocket)

	addr := fmt.Sprintf("0.0.0.0:%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	log.Printf("Server ID: %s", hub.ServerID)
	log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html", cfg.ServerPort)
	log.Printf("WebSocket available at ws://localhost:%s/ws", cfg.ServerPort)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getServerID() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return fmt.Sprintf("%s-%d", hostname, time.Now().Unix())
}
