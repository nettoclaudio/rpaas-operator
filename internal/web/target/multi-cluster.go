package target

import (
	"context"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/tsuru/rpaas-operator/internal/config"
	"github.com/tsuru/rpaas-operator/internal/pkg/rpaas"
	"github.com/tsuru/rpaas-operator/pkg/observability"
	extensionsruntime "github.com/tsuru/rpaas-operator/pkg/runtime"
	"k8s.io/client-go/rest"
	sigsk8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ Factory = &multiClusterFactory{}

type missingParamsError struct {
	Msg           string   `json:"msg"`
	MissingParams []string `json:"missing_params"`
}

func (e *missingParamsError) Error() string {
	return e.Msg
}

func (e *missingParamsError) IsValidation() bool {
	return true
}

var ErrNoClusterProvided = &missingParamsError{
	Msg:           "No cluster address provided",
	MissingParams: []string{"cluster"},
}

type multiClusterFactory struct {
	tokens        sync.Map
	clusters      []config.ClusterConfig
	managersMutex sync.RWMutex
	managers      map[managersCacheKey]rpaas.RpaasManager
}

type managersCacheKey struct {
	clusterName    string
	poolName       string
	clusterAddress string
}

func NewMultiClustersFactory(clusters []config.ClusterConfig) Factory {
	return &multiClusterFactory{
		clusters:      clusters,
		tokens:        sync.Map{},
		managersMutex: sync.RWMutex{},
		managers:      map[managersCacheKey]rpaas.RpaasManager{},
	}
}

func (m *multiClusterFactory) Manager(ctx context.Context, headers http.Header) (rpaas.RpaasManager, error) {
	clusterName := headers.Get("X-Tsuru-Cluster-Name")
	address := headers.Get("X-Tsuru-Cluster-Addresses")

	if address == "" {
		return nil, ErrNoClusterProvided
	}

	poolName := headers.Get("X-Tsuru-Pool-Name")
	cacheKey := managersCacheKey{clusterName, poolName, address}

	m.managersMutex.RLock()
	manager := m.managers[cacheKey]
	m.managersMutex.RUnlock()

	if manager != nil {
		return manager, nil
	}

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("cluster.name", clusterName)
		span.SetTag("cluster.address", address)
		span.SetTag("pool.name", poolName)
	}

	bearerToken, err := m.getToken(clusterName)
	if err != nil {
		return nil, err
	}
	kubernetesRestConfig := &rest.Config{
		Host:          address,
		BearerToken:   bearerToken,
		WrapTransport: observability.OpentracingTransport,
	}
	k8sClient, err := sigsk8sclient.New(kubernetesRestConfig, sigsk8sclient.Options{Scheme: extensionsruntime.NewScheme()})
	if err != nil {
		return nil, err
	}

	manager, err = rpaas.NewK8S(kubernetesRestConfig, k8sClient, clusterName, poolName)
	if err != nil {
		return nil, err
	}

	m.managersMutex.Lock()
	defer m.managersMutex.Unlock()

	m.managers[cacheKey] = manager
	return manager, nil
}

func (m *multiClusterFactory) getToken(clusterName string) (string, error) {
	var defaultCluster *config.ClusterConfig = nil
	for _, cluster := range m.clusters {
		if cluster.Default || cluster.Name == clusterName {
			defaultCluster = &cluster
			break
		}
	}

	if defaultCluster == nil {
		return "", nil
	}

	if defaultCluster.Token != "" {
		return defaultCluster.Token, nil
	}

	if defaultCluster.TokenFile != "" {
		return m.readTokenFile(defaultCluster.TokenFile)
	}

	return "", nil
}

func (m *multiClusterFactory) readTokenFile(tokenFile string) (string, error) {
	tokenInterface, ok := m.tokens.Load(tokenFile)

	if ok {
		return tokenInterface.(string), nil
	}

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}

	m.tokens.Store(tokenFile, string(token))

	return string(token), nil
}
