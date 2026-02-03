package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"api_kino/cmd"
	"api_kino/config/app"
	"api_kino/config/database"
	"api_kino/routes"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

func main() {
	if app.Config().GenerateModel {
		cmd.GenerateModel()
		return
	}
	router := routes.Router()
	if app.Config().EnableFreeMem {
		go periodicFree(30 * time.Second)
	}
	s := &http.Server{
		Addr:           ":" + app.Config().Port,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	database.ConnectAll()
	//_, _ = database.Connect(database.DBConfig())
	go func() {
		if err := s.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			fmt.Println(err)
			//zap.L().Fatal("HTTP Server", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("======================================= Shutting down server... =======================================")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("======================================= Server forced to shutdown =======================================", err)
	}
	fmt.Println("======================================= Server exiting =======================================")
}

func periodicFree(d time.Duration) {
	tick := time.Tick(d)
	for range tick {
		fmt.Println("======================================= Memory Released =======================================")
		debug.FreeOSMemory()
	}
}
