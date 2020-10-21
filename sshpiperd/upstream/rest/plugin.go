package rest

import (
	"log"

	"github.com/tg123/sshpiper/sshpiperd/upstream"
)

type plugin struct {
	Config struct {
		URL	string `long:"upstream-rest-url" default:"https://127.0.0.1" description:"URL of REST API" env:"SSHPIPERD_UPSTREAM_REST_URL"  ini-name:"upstream-rest-url"`
	}

	logger *log.Logger
}

// The name of the Plugin
func (p *plugin) GetName() string {
	return "rest"
}

// A ref to a struct which holds the options for the plugins
// will be populated by cmd or other plugin runners
func (p *plugin) GetOpts() interface{} {
	return &p.Config
}

// Will be called before the Plugin is used to ensure the Plugin is ready
func (p *plugin) Init(logger *log.Logger) error {
	p.logger = logger
	return nil
}

func (p *plugin) GetHandler() upstream.Handler {
	return p.findUpstream
}

func init() {
	upstream.Register("rest", &plugin{})
}
