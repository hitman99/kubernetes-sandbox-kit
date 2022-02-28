package cmd

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/hitman99/kubernetes-sandbox/internal/admin"
	"github.com/hitman99/kubernetes-sandbox/internal/utils"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"time"
)

var adminCmd = &cobra.Command{
	Use: "admin",
	Run: func(cmd *cobra.Command, args []string) {
		runAdminApi()
		os.Exit(0)
	},
}

func runAdminApi() {
	logger := utils.SetupLogger()
	ac := admin.MustNewAdminClient()
	auth := admin.NewAuthMiddleware()
	mainRouter := chi.NewRouter()
	mainRouter.Route("/admin/participants", func(r chi.Router) {
		r.Use(auth.Middleware)
		r.Delete("/{pid}", ac.RemoveParticipantHandler())
		r.Delete("/", ac.ResetEnvironmentHandler())
		r.Get("/", ac.GetParticipantsHandler())
	})
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      mainRouter,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}
	logger.Info("started reg http server on port " + srv.Addr)
	logger.Fatal(srv.ListenAndServe())
}
