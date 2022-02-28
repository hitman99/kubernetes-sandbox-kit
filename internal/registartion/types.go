package registartion

import (
	common_http "github.com/hitman99/kubernetes-sandbox/internal/common-http"
	"net/http"
)

type Registration struct {
	User         User           `json:"user"`
	Kubernetes   KubernetesInfo `json:"kubernetes"`
	Instructions string         `json:"instructions"`
	common_http.Response
}

type User struct {
	Email string `json:"email"`
	Id    string `json:"id"`
}

type KubernetesInfo struct {
	Namespace     string `json:"namespace"`
	ServerVersion string `json:"serverVersion"`
}

func (reg *Registration) Bind(r *http.Request) error {
	return nil
}
