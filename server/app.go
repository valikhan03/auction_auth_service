package server

import(
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"


	"auction_auth_service/auth"
	"auction_auth_service/auth/repository/authdatabase"


	authhttp "auction_auth_service/auth/delivery/http"
	authUsecase "auction_auth_service/auth/usecase"


	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

type App struct {
	httpServer *http.Server
	authUC     auth.UseCase
}

func NewApp() *App {

	postgresDB := initPostgreDB()
	

	authRepos := authdatabase.NewUserRepository(postgresDB)

	return &App{
		authUC:    authUsecase.NewAuthUseCase(authRepos, "Pstre12e_9fQz", []byte("pwr12qxk90"), 10),
	}
}

func initPostgreDB() *sqlx.DB {
	db, err := sqlx.Connect("pgx", ReadPostgresConfigs())
	if err != nil {
		log.Fatal(err)
	}

	return db
}



func (a *App) Run() error {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	//needed to set up http handlers as endpoints
	authhttp.RegisterAuthHTTPEndpoints(router, a.authUC)

	authMiddleware := authhttp.NewAuthMiddleware(a.authUC)

	router.POST("/check-access", authMiddleware.Handle)

	a.httpServer = &http.Server{
		Addr:           ":8090",
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Server failed : %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
