package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/moby/moby/client"
	"github.com/octago/polygon/api"
)

type PolygonServer struct {
	dockerCli    *client.Client
	templatesDir string
}

func New(templatesDir string) (*PolygonServer, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &PolygonServer{
		dockerCli:    cli,
		templatesDir: templatesDir,
	}, nil
}

func (p *PolygonServer) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateReply, error) {
	template := req.TemplateId
	templatePath := filepath.Join(p.templatesDir, template)
	if _, err := os.Lstat(p.templatesDir); err != nil {
		return nil, err
	}
	configFile := filepath.Join(templatePath, "config.json")
	hostConfigFile := filepath.Join(templatePath, "hostconfig.json")

	rawConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var cfg container.Config
	if err := json.Unmarshal(rawConfig, &cfg); err != nil {
		return nil, err
	}

	rawHostConfig, err := ioutil.ReadFile(hostConfigFile)
	if err != nil {
		return nil, err
	}
	var hostCfg container.HostConfig
	if err := json.Unmarshal(rawHostConfig, &hostCfg); err != nil {
		return nil, err
	}

	cont, err := p.dockerCli.ContainerCreate(ctx, &cfg, &hostCfg, nil, "")
	if err != nil {
		return nil, err
	}
	if err := p.dockerCli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	return &api.CreateReply{Stand: &api.Stand{Id: cont.ID}}, nil
}

func (p *PolygonServer) Get(context.Context, *api.GetRequest) (*api.GetReply, error) {
	panic("not implemented")
}

func (p *PolygonServer) Attach(stream api.PolygonServer_AttachServer) error {
	ctx := stream.Context()
	execCfg := types.ExecConfig{
		Tty:          true,
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          []string{"sh"},
	}
	firstMsg, err := stream.Recv()
	if err != nil {
		return err
	}
	contID := firstMsg.StandId
	execResp, err := p.dockerCli.ContainerExecCreate(ctx, contID, execCfg)
	if err != nil {
		return err
	}
	hijack, err := p.dockerCli.ContainerExecAttach(ctx, execResp.ID, execCfg)
	if err != nil {
		return err
	}
	defer hijack.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				fmt.Println("ERROR", err)
				break
			}
			if _, err := hijack.Conn.Write(msg.Chunk); err != nil {
				fmt.Println("ERROR WRITING STDIN", err)
			}
		}
		wg.Done()
	}()
	go func() {
		for {
			b := make([]byte, 4096)
			read, err := hijack.Reader.Read(b)
			if err != nil {
				fmt.Println("ERROR READING STDOUT", err)
				break
			}
			if err := stream.Send(&api.StreamChunk{Chunk: b[:read]}); err != nil {
				fmt.Println("ERROR", err)
				break
			}
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func (p *PolygonServer) CancelStand(context.Context, *api.CancelRequest) (*api.CancelReply, error) {
	panic("not implemented")
}
