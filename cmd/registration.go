package cmd

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/hitman99/kubernetes-sandbox/internal/registartion"
	"github.com/hitman99/kubernetes-sandbox/internal/utils"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"time"
)

var registerCmd = &cobra.Command{
	Use: "registration",
	Run: func(cmd *cobra.Command, args []string) {
		runRegister()
		os.Exit(0)
	},
}

func runRegister() {
	logger := utils.SetupLogger()
	reg := registartion.New()
	mainRouter := chi.NewRouter()
	mainRouter.Route("/register", func(r chi.Router) {
		r.Post("/", reg.CreateReg())
	})
	mainRouter.Route("/kubeconfig", func(r chi.Router) {
		r.Get("/{userId}", reg.KubeconfigHandler())
	})
	fs := http.StripPrefix("/", http.FileServer(http.Dir("frontend/dist")))
	mainRouter.Route("/", func(r chi.Router) {
		r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		}))
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      mainRouter,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}
	logger.Printf("started reg http server on port %s", srv.Addr)
	logger.Fatal(srv.ListenAndServe())
}
