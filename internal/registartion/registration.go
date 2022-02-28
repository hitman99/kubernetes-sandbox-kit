package registartion

import (
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-redis/redis/v7"
	"github.com/gofrs/uuid"
	common_http "github.com/hitman99/kubernetes-sandbox/internal/common-http"
	"github.com/hitman99/kubernetes-sandbox/internal/config"
	"github.com/hitman99/kubernetes-sandbox/internal/kubernetes"
	"github.com/hitman99/kubernetes-sandbox/internal/storage"
	"github.com/hitman99/kubernetes-sandbox/internal/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"regexp"
)

type Reg struct {
	redisClient      *redis.Client
	kubeClient       kubernetes.Client
	logger           *logrus.Logger
	emailRegex       *regexp.Regexp
	instructionsPath string
}

func New() *Reg {
	logger := utils.SetupLogger()
	cfg, _ := config.Get()
	rcli := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})
	_, err := rcli.Ping().Result()
	if err != nil {
		logger.WithError(err).Fatal("cannot reach redis, exiting")
	}
	r := &Reg{
		redisClient:      rcli,
		logger:           logger,
		kubeClient:       kubernetes.MustNewClient(),
		emailRegex:       regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"),
		instructionsPath: cfg.InstructionsPath,
	}
	return r
}

func (s *Reg) CreateReg() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reg := &Registration{}
		err := render.Bind(r, reg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !s.validateEmail(reg.User.Email) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(common_http.Response{
				Success: false,
				Message: "invalid email",
			})
			return
		}
		available, err := s.emailAvailable(reg.User.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(common_http.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		if !available {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(common_http.Response{
				Success: false,
				Message: "user with this email already exists",
			})
			return
		}
		reg.User.Id = uuid.Must(uuid.NewV4()).String()
		reg.Kubernetes.Namespace = kubernetes.K8S_NS_PREFIX + reg.User.Id
		s.logger.WithFields(logrus.Fields{"email": reg.User.Email, "namespace": reg.Kubernetes.Namespace, "pid": reg.User.Id}).Info("new user")
		err = s.persistParticipant(reg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(common_http.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		err = s.kubeClient.CreateNamespace(reg.Kubernetes.Namespace)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(common_http.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		reg.Kubernetes.ServerVersion = s.kubeClient.GetVersion()
		w.WriteHeader(http.StatusCreated)
		reg.Success = true
		data, err := os.ReadFile(s.instructionsPath)
		if err == nil {
			reg.Instructions = base64.StdEncoding.EncodeToString(data)
		}
		err = json.NewEncoder(w).Encode(reg)
		if err != nil {
			s.logger.WithError(err).WithField("namespace", reg.Kubernetes.Namespace).Error("cannot send response")
		}
	}
}

func (s *Reg) ListRegs() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := s.redisClient.LRange(storage.REDIS_LIST_KEY, 0, -1).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			studs, err := json.Marshal(students)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(common_http.Response{
					Success: false,
					Message: err.Error(),
				})
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(studs)
		}
	}
}

func (s *Reg) KubeconfigHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pid := chi.URLParam(r, "userId")
		blob, err := s.redisClient.Get(pid).Result()
		if err != nil {
			if err != redis.Nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(common_http.Response{
					Success: false,
					Message: "storage failed, please try again",
				})
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(common_http.Response{
					Success: false,
					Message: "this participant does not exist, please register again",
				})
				return
			}
		}

		reg := &Registration{}
		err = json.Unmarshal([]byte(blob), reg)
		if err != nil {
			s.logger.WithError(err).Error("cannot unmarshal")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		kubeconfig, err := s.kubeClient.GetKubeconfig(reg.Kubernetes.Namespace)
		if err != nil {
			s.logger.WithError(err).WithField("namespace", reg.Kubernetes.Namespace).Error("cannot get kubeconfig for namespace")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(kubeconfig))
		}
	}
}

func (s *Reg) persistParticipant(r *Registration) error {
	blob, err := json.Marshal(r)
	if err != nil {
		return err
	}
	pipe := s.redisClient.TxPipeline()
	pipe.LPush(storage.REDIS_LIST_KEY, r.User.Id)
	pipe.Set(r.User.Id, blob, 0)
	_, err = pipe.Exec()
	return storage.WrapRedisErr(err)
}

func (s *Reg) validateEmail(email string) bool {
	return s.emailRegex.MatchString(email)
}

func (s *Reg) emailAvailable(email string) (bool, error) {
	res, err := s.redisClient.LRange(storage.REDIS_LIST_KEY, 0, -1).Result()
	if err != nil {
		return false, err
	}
	pipeRes := []*redis.StringCmd{}
	pipe := s.redisClient.TxPipeline()
	for _, pid := range res {
		pipeRes = append(pipeRes, pipe.Get(pid))
	}
	_, err = pipe.Exec()
	if err != nil {
		return false, err
	}
	for _, pr := range pipeRes {
		var reg Registration
		err := json.Unmarshal([]byte(pr.Val()), &reg)
		if err != nil {
			s.logger.WithError(err).WithField("participant", pr.Val()).Error("failed to unmarshal")
			return false, err
		}
		if reg.User.Email == email {
			return false, nil
		}
	}
	return true, nil
}

func (s *Reg) GetInstructions() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(s.instructionsPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(common_http.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(base64.StdEncoding.EncodeToString(data))
	}
}
