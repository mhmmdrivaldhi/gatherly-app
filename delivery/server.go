package delivery

import (
	"fmt"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gatherly-app/config"
	"gatherly-app/delivery/controllers"
	"gatherly-app/delivery/middleware"
	_ "gatherly-app/docs"
	"gatherly-app/models"
	"gatherly-app/repositories"
	"gatherly-app/service"
	"gatherly-app/usecase"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	db              *gorm.DB
	userUC          usecase.UserUsecase
	ticketUC        usecase.TicketUseCase
	eventAttendeeUC usecase.EventAttendeeUseCase
	eventUC         usecase.EventsUsecase
	transactionUC   usecase.TransactionUsecase
	authUC          usecase.AuthenticationUseCase
	jwtService      service.JwtService
	midtransService service.MidtransService
	engine          *gin.Engine
	host            string
}

func (s *Server) initRoute() {
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rgAuth := s.engine.Group("/api/auth")
	controllers.NewAuthController(s.authUC, rgAuth)

	rgV1 := s.engine.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)

	// Public routes
	controllers.NewUserController(s.userUC, rgV1)
	controllers.NewTransactionController(s.transactionUC, rgV1).RegisterPublicRoutes()

	// Authenticated routes
	authGroup := rgV1.Group("")
	authGroup.Use(authMiddleware.RequireToken())
	{
		controllers.NewTicketController(s.ticketUC, authGroup).Route()
		controllers.NewEventAttendeeController(s.eventAttendeeUC, authGroup).Route()
		controllers.NewEventsController(s.eventUC, authGroup).Route()
		controllers.NewTransactionController(s.transactionUC, authGroup).Route()
	}
}

func (s *Server) initMigration() {
	err := s.db.AutoMigrate(
		&models.Ticket{},
		&models.Transactions{},
		&models.User{},
		&models.Event{},
	)

	if err != nil {
		log.Fatal("Failed to migrate: ", err)
	}

	s.db.Migrator().CreateConstraint(&models.Transactions{}, "fk_transactions_users")
	s.db.Migrator().CreateConstraint(&models.Transactions{}, "foreignKey")

	log.Println("Migrated Successfully")
}

func (s *Server) Run() {
	s.engine.Use(func(c *gin.Context) {
		// Middleware untuk mendapatkan IP asli (jika behind proxy)
		if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
			c.Set("client_ip", forwarded)
		} else {
			c.Set("client_ip", c.ClientIP())
		}
		c.Next()
	})

	s.initRoute()     // Inisialisasi routing
	s.initMigration() // Jalankan migrasi

	if err := s.engine.Run(s.host); err != nil {
		log.Fatalf("server not running on host %s, because error %v", s.host, err)
	}
}

func NewServer() *Server {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error yang terjadi :", err.Error())
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	jwtService := service.NewJwtService(cfg.TokenConfig)

	client := resty.New()
	midtransService := service.NewMidtransService(client, cfg.MidtransServerKey)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	eventRepo := repositories.NewEventsRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)
	eventAttendeeRepo := repositories.MakeNewEventAttendeeRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo)
	eventUsecase := usecase.NewEventUsecase(eventRepo, eventAttendeeRepo)
	ticketUseCase := usecase.NewTicketUseCase(ticketRepo)
	transactionUseCase := usecase.NewTransactionUsecase(transactionRepo, midtransService)
	eventAttendeeUseCase := usecase.NewEventAttendeeUseCase(eventAttendeeRepo, eventRepo, ticketRepo, transactionUseCase)
	authUseCase := usecase.NewAuthenticationUseCase(userRepo, jwtService)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		db:              db, // Simpan db ke struct Server
		ticketUC:        ticketUseCase,
		userUC:          userUsecase,
		eventUC:         eventUsecase,
		transactionUC:   transactionUseCase,
		eventAttendeeUC: eventAttendeeUseCase,
		engine:          engine,
		host:            host,
		authUC:          authUseCase,
		jwtService:      jwtService,
		midtransService: midtransService,
	}
}
