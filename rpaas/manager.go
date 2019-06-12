package rpaas

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/tsuru/rpaas-operator/pkg/apis/extensions/v1alpha1"
)

type ConfigurationBlock struct {
	Name    string `form:"block_name" json:"block_name"`
	Content string `form:"content" json:"content"`
}

// ConfigurationBlockHandler defines some functions to handle the custom
// configuration blocks from an instance.
type ConfigurationBlockHandler interface {
	// DeleteBlock removes the configuration block named by blockName. It returns
	// a nil error meaning it was successful, otherwise a non-nil one which
	// describes the reached problem.
	DeleteBlock(ctx context.Context, instanceName, blockName string) error

	// ListBlocks returns all custom configuration blocks from instance (which
	// name is instanceName). It returns a nil error meaning it was successful,
	// otherwise a non-nil one which describes the reached problem.
	ListBlocks(ctx context.Context, instanceName string) ([]ConfigurationBlock, error)

	// UpdateBlock overwrites the older configuration block content with the one.
	// Whether the configuration block entry does not exist, it will already be
	// created with the new content. It returns a nil error meaning it was
	// successful, otherwise a non-nil one which describes the reached problem.
	UpdateBlock(ctx context.Context, instanceName string, block ConfigurationBlock) error
}

type File struct {
	Name    string
	Content []byte
}

func (f File) SHA256() string {
	return fmt.Sprintf("%x", sha256.Sum256(f.Content))
}

func (f File) MarshalJSON() ([]byte, error) {
	return json.Marshal(&map[string]interface{}{
		"name":    f.Name,
		"content": f.Content,
		"sha256":  f.SHA256(),
	})
}

type ExtraFileHandler interface {
	CreateExtraFiles(ctx context.Context, instanceName string, files ...File) error
	DeleteExtraFiles(ctx context.Context, instanceName string, filenames ...string) error
	GetExtraFiles(ctx context.Context, instanceName string) ([]File, error)
	UpdateExtraFiles(ctx context.Context, instanceName string, files ...File) error
}

type Route struct {
	Path        string `form:"path" json:"path"`
	Content     string `form:"content" json:"content,omitempty"`
	Destination string `form:"destination" json:"destination,omitempty"`
	HTTPSOnly   bool   `form:"https_only" json:"https_only"`
}

type RouteHandler interface {
	DeleteRoute(ctx context.Context, instanceName, path string)
	GetRoutes(ctx context.Context, instanceName string) ([]Route, error)
	UpdateRoute(ctx context.Context, instanceName string, route Route) error
}

type CreateArgs struct {
	Name        string   `json:"name" form:"name"`
	Plan        string   `json:"plan" form:"plan"`
	Team        string   `json:"team" form:"team"`
	User        string   `json:"user" form:"user"`
	Tags        []string `json:"tags" form:"tags"`
	EventID     string   `json:"eventid" form:"eventid"`
	Description string   `json:"description" form:"description"`
	Flavor      string   `json:"flavor" form:"flavor"`
	IP          string   `json:"ip" form:"ip"`
}

type PodStatusMap map[string]PodStatus

type PodStatus struct {
	Running bool   `json:"running"`
	Status  string `json:"status"`
	Address string `json:"address"`
}

type RpaasManager interface {
	ConfigurationBlockHandler
	ExtraFileHandler
	RouteHandler

	UpdateCertificate(ctx context.Context, instance, name string, cert tls.Certificate) error
	CreateInstance(ctx context.Context, args CreateArgs) error
	DeleteInstance(ctx context.Context, name string) error
	GetInstance(ctx context.Context, name string) (*v1alpha1.RpaasInstance, error)
	GetInstanceAddress(ctx context.Context, name string) (string, error)
	GetInstanceStatus(ctx context.Context, name string) (PodStatusMap, error)
	Scale(ctx context.Context, name string, replicas int32) error
	GetPlans(ctx context.Context) ([]v1alpha1.RpaasPlan, error)
}
