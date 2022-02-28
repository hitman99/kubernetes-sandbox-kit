package kubernetes

import (
    stdErr "errors"
    "fmt"
    "github.com/hitman99/kubernetes-sandbox/internal/config"
    "github.com/hitman99/kubernetes-sandbox/internal/utils"
    log "github.com/sirupsen/logrus"
    v1 "k8s.io/api/core/v1"
    v1rbac "k8s.io/api/rbac/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "sync"
)

const K8S_NS_PREFIX = "sandbox-"

type PodInfo struct {
    Ip       string
    Hostname string
}

type Client interface {
    CreateNamespace(namespace string) error
    DeleteNamespace(namespace string) error
    GetKubeconfig(namespace string) (string, error)
    GetVersion() string
}

type client struct {
    clientset *kubernetes.Clientset
    logger    *log.Logger
    config.KubernetesConfig
    _m sync.Mutex
}

func MustNewClient() Client {
    logger := utils.SetupLogger()
    cfg, updates := config.GetKubernetesConfig()
    var (
        err        error
        restConfig *rest.Config
    )
    if cfg.DevMode {
        restConfig, err = clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
    } else {
        restConfig, err = rest.InClusterConfig()
    }
    if err != nil {
        logger.WithError(err).Fatal("failed to connect to kubernetes")
    }
    // creates the clientset
    clientset, err := kubernetes.NewForConfig(restConfig)
    if err != nil {
        logger.Fatalf("cannot init kubernetes config: %s", err.Error())
    }
    c := &client{
        clientset:        clientset,
        logger:           logger,
        _m:               sync.Mutex{},
        KubernetesConfig: cfg,
    }
    c.syncConfig(updates)
    return c
}

func (k *client) CreateNamespace(namespace string) error {
    _, err := k.clientset.CoreV1().Namespaces().Get(namespace, v1meta.GetOptions{})
    if err == nil || errors.IsNotFound(err) {
        _, err := k.clientset.CoreV1().Namespaces().Create(&v1.Namespace{
            ObjectMeta: v1meta.ObjectMeta{
                Name: namespace,
            },
        })

        if err != nil {
            k.logger.WithError(err).WithField("namespace", namespace).Error("cannot create namespace")
            return err
        }
        _, err = k.clientset.RbacV1().RoleBindings(namespace).Create(&v1rbac.RoleBinding{
            ObjectMeta: v1meta.ObjectMeta{
                Name: "admin",
            },
            Subjects: []v1rbac.Subject{{
                Kind:      "ServiceAccount",
                APIGroup:  "",
                Name:      "default",
                Namespace: namespace,
            }},
            RoleRef: v1rbac.RoleRef{
                APIGroup: "rbac.authorization.k8s.io",
                Kind:     "ClusterRole",
                Name:     "admin",
            },
        })
        if err != nil {
            k.logger.WithError(err).WithField("namespace", namespace).Error("cannot create rolebinding")
            return err
        }
    } else {
        k.logger.WithField("namespace", namespace).Info("namespace already exists")
    }
    return nil
}

func (k *client) GetKubeconfig(namespace string) (string, error) {
    sa, err := k.clientset.CoreV1().ServiceAccounts(namespace).Get("default", v1meta.GetOptions{})
    if err != nil {
        return "", fmt.Errorf("kubernetes error", err)
    }
    if len(sa.Secrets) == 1 {
        secret, err := k.clientset.CoreV1().Secrets(namespace).Get(sa.Secrets[0].Name, v1meta.GetOptions{})
        if err != nil {
            return "", fmt.Errorf("kubernetes error", err)
        }
        k._m.Lock()
        defer k._m.Unlock()
        return fmt.Sprintf(kubeconfig, k.ApiCA, k.ApiURI, namespace, secret.Data["token"]), nil
    } else {
        return "", stdErr.New("no secrets found")
    }
}

func (k *client) DeleteNamespace(namespace string) error {
    err := k.clientset.CoreV1().Namespaces().Delete(namespace, &v1meta.DeleteOptions{})
    if err != nil && !errors.IsNotFound(err) {
        k.logger.Printf("cannot delete namespace %s, %s", namespace, err.Error())
        return err
    }
    return nil
}

func (k *client) syncConfig(updates <-chan config.KubernetesConfig) {
    go func() {
        for c := range updates {
            k._m.Lock()
            k.KubernetesConfig = c
            k._m.Unlock()
        }
    }()
}

func (k *client) GetVersion() string {
    ver, _ := k.clientset.ServerVersion()
    return ver.GitVersion
}
