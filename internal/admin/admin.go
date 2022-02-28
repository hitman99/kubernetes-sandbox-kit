package admin

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/hitman99/kubernetes-sandbox/internal/kubernetes"
	"github.com/hitman99/kubernetes-sandbox/internal/registartion"
	"github.com/hitman99/kubernetes-sandbox/internal/storage"
	"github.com/hitman99/kubernetes-sandbox/internal/utils"
	"net/http"
)

type Client interface {
	// RemoveParticipantHandler removes single participant and its dependencies
	RemoveParticipantHandler() func(w http.ResponseWriter, r *http.Request)
	// ResetEnvironmentHandler resets the environment completely by removing all participants and their dependencies
	ResetEnvironmentHandler() func(w http.ResponseWriter, r *http.Request)
	GetParticipantsHandler() func(w http.ResponseWriter, r *http.Request)
}

type admin struct {
	kc     kubernetes.Client
	rc     storage.RedisClient
	logger *logrus.Logger
}

func MustNewAdminClient() Client {
	return &admin{
		kc:     kubernetes.MustNewClient(),
		rc:     storage.MustNewRedisClient(),
		logger: utils.SetupLogger(),
	}
}

func (a *admin) RemoveParticipantHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pid := chi.URLParam(r, "pid")
		a.logger.WithField("participantId", pid).Info("removing participant")
		_, err := a.rc.Get(pid).Result()
		if err != nil {
			if err != redis.Nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.Error(w, "unknown participant", http.StatusNotFound)
			}
			return
		}
		// delete kubernetes namespace
		if err := a.kc.DeleteNamespace(kubernetes.K8S_NS_PREFIX + pid); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		a.logger.WithField("namespace", kubernetes.K8S_NS_PREFIX+pid).Info("k8s namespace deleted")
		// cleanup redis
		pipe := a.rc.TxPipeline()
		pipe.Del(pid)
		pipe.LRem(storage.REDIS_LIST_KEY, 0, pid)
		_, err = pipe.Exec()
		if err := storage.WrapRedisErr(err); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (a *admin) ResetEnvironmentHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := a.rc.LRange(storage.REDIS_LIST_KEY, 0, -1).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		a.logger.WithField("count", len(res)).Info("removing all participants")
		for _, pid := range res {
			// delete kubernetes namespace
			if err := a.kc.DeleteNamespace(kubernetes.K8S_NS_PREFIX + pid); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// cleanup redis
			pipe := a.rc.TxPipeline()
			pipe.Del(pid)
			pipe.LRem(storage.REDIS_LIST_KEY, 0, pid)
			_, err = pipe.Exec()
			if err := storage.WrapRedisErr(err); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			a.logger.WithField("participantId", pid).Info("participant removed")
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (a *admin) GetParticipantsHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := a.rc.LRange(storage.REDIS_LIST_KEY, 0, -1).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pipeRes := []*redis.StringCmd{}
		participants := make([]*registartion.Registration, 0, len(res))
		pipe := a.rc.TxPipeline()
		for _, pid := range res {
			pipeRes = append(pipeRes, pipe.Get(pid))
		}
		_, err = pipe.Exec()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, pr := range pipeRes {
			var reg registartion.Registration
			err := json.Unmarshal([]byte(pr.Val()), &reg)
			if err != nil {
				a.logger.WithError(err).WithField("participant", pr.Val()).Error("failed to unmarshal")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			reg.Success = true
			participants = append(participants, &reg)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(participants)
	}
}
