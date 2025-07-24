package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/salamsites/minio-pkg"
	sminiochat "github.com/salamsites/minio-pkg/client/chat"
	"github.com/salamsites/minio-pkg/client/feed"
	"github.com/salamsites/minio-pkg/client/music"
	"github.com/salamsites/minio-pkg/client/user"
	"github.com/salamsites/minio-pkg/test"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logFile := pwd + "/test/test_files/" + test.LogFileName
	logger := slog.GetLogger(pwd+"/test/test_files", test.LogFileName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for SIGINT and SIGTERM signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// ---old
	//sminioClient, err := sminio.NewClient(sminio.Options{
	//	Endpoint:        "10.192.1.45:9000",
	//	AccessKeyID:     "minioadmin",
	//	SecretAccessKey: "minioadmin",
	//})

	sminioUserClient, err := user.NewUserClient(sminio.Options{
		Endpoint:        "10.192.1.115:9000",
		AccessKeyID:     "admin",
		SecretAccessKey: "password",
	})

	sminioFeedClient, err := feed.NewFeedClient(sminio.Options{
		Endpoint:        "10.192.1.115:9000",
		AccessKeyID:     "admin",
		SecretAccessKey: "password",
	})

	sminioMusicClient, err := music.NewMusicClient(sminio.Options{
		Endpoint:        "10.192.1.115:9000",
		AccessKeyID:     "admin",
		SecretAccessKey: "password",
	})

	sminioChatClient, err := sminiochat.NewChatClient(sminio.Options{
		Endpoint:        "10.192.1.115:9000",
		AccessKeyID:     "admin",
		SecretAccessKey: "password",
	})
	// Create a server
	srv := &http.Server{
		Addr:         test.PORT,
		Handler:      getRouter(cancel, sminioUserClient, sminioFeedClient, sminioMusicClient, sminioChatClient, logger, pwd),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	//
	go func() {
		log.Println("Server started")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Wait for either a shutdown signal or an API call to shut down
	select {
	case sig := <-sigCh:
		log.Printf("Received signal: %v", sig)
	case <-ctx.Done():
		log.Println("Received shutdown request from API")
	}

	// Start the shutdown process
	log.Println("Shutting down...")

	err = os.Remove(logFile)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("File deleted successfully.")

	err = os.RemoveAll(pwd + "/test/uploaded_files")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("File deleted successfully.")

	// Create a context with a timeout to allow the server some time to finish ongoing requests
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully")

}

func getRouter(cancel context.CancelFunc, sminioUserClient sminio.UserClient, sminioFeedClient sminio.FeedClient, sminioMusicClient sminio.MusicClient, sminioChatClient sminio.ChatClient, logger *slog.Logger, pwd string) *mux.Router {

	router := mux.NewRouter()

	limitter := shttp.NewRateLimiter()
	middleware := shttp.NewMiddleware(logger, "", limitter)

	routerHandler := &handler{
		cancel:            cancel,
		logger:            logger,
		middleware:        middleware,
		pwd:               pwd,
		sminioUserClient:  sminioUserClient,
		sminioFeedClient:  sminioFeedClient,
		sminioMusicClient: sminioMusicClient,
		sminioChatClient:  sminioChatClient,
	}
	routerHandler.Register(router)

	return router
}

type handler struct {
	cancel            context.CancelFunc
	logger            *slog.Logger
	middleware        *shttp.Middleware
	pwd               string
	sminioUserClient  sminio.UserClient
	sminioFeedClient  sminio.FeedClient
	sminioMusicClient sminio.MusicClient
	sminioChatClient  sminio.ChatClient
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(test.UploadImageURL, h.middleware.Base(h.uploadHandler)).Methods(http.MethodPost)
	router.HandleFunc(test.ShutdownURL, h.middleware.Base(h.shutdown)).Methods(http.MethodGet)
}

func (h *handler) uploadHandler(w http.ResponseWriter, r *http.Request) shttp.Response {
	fmt.Println("start uploadHandler\n ")
	sizes, err := h.sminioUserClient.UploadAvatar(r.Context(), 321, r, test.KEY)
	if err.StatusCode > 0 {
		fmt.Printf("\nerr: \n%v\n", err)
		return shttp.Result.SetStatusCode(err.StatusCode).SetData(err.Message)
	}
	fmt.Printf("avatar size: %d\n", sizes)
	fmt.Println("avatar uploaded successfully")

	sizess, err := h.sminioFeedClient.UploadFeed(r.Context(), 321, r, test.KEY)
	if err.StatusCode > 0 {
		fmt.Printf("\nerr: \n%v\n", err)
		return shttp.Result.SetStatusCode(err.StatusCode).SetData(err.Message)
	}
	fmt.Printf("feed size: %d\n", sizess)
	fmt.Println("feed uploaded successfully")

	chatSize, err := h.sminioChatClient.UploadFile(r.Context(), 321, r, test.KEY)
	if err.StatusCode > 0 {
		fmt.Printf("\nerr: \n%v\n", err)
		return shttp.Result.SetStatusCode(err.StatusCode).SetData(err.Message)
	}
	fmt.Printf("chat size: %d\n", chatSize)
	fmt.Println("file uploaded successfully")

	var musicImagePath string
	msize, err := h.sminioMusicClient.UploadMusicPhoto(r.Context(), 321, musicImagePath)
	if err.StatusCode > 0 {
		fmt.Printf("\nerr: \n%v\n", err)
		return shttp.Result.SetStatusCode(err.StatusCode).SetData(err.Message)
	}
	fmt.Printf("music photo size: %d\n", msize)
	fmt.Println("music photo uploaded successfully")

	var musicPath string
	result, err := h.sminioMusicClient.UploadMusic(r.Context(), 321, musicPath)
	if err.StatusCode > 0 {
		fmt.Printf("\nerr: \n%v\n", err)
		return shttp.Result.SetStatusCode(err.StatusCode).SetData(err.Message)
	}
	fmt.Println("music uploaded successfully", result)

	return shttp.Success
}

func (h *handler) shutdown(w http.ResponseWriter, r *http.Request) shttp.Response {

	h.cancel()

	return shttp.Success.SetData("Success shutdown")
}
