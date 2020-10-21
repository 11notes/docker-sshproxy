package rest

import (
	"fmt"
	
	upstreamprovider "github.com/tg123/sshpiper/sshpiperd/upstream"
)

func (p *plugin) ListPipe() ([]upstreamprovider.Pipe, error) {
	return nil, nil
}

func (p *plugin) CreatePipe(opt upstreamprovider.CreatePipeOption) error {
	
}

func (p *plugin) RemovePipe(name string) error {
	return nil
}
