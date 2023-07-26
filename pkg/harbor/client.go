package harbor

import (
	"context"
	"time"

	"github.com/mittwald/goharbor-client/v5/apiv2"
	"github.com/mittwald/goharbor-client/v5/apiv2/model"
	client_cf "github.com/mittwald/goharbor-client/v5/apiv2/pkg/config"

	"github.com/zhangguanzhang/harbor-cleaner/pkg/config"
)

type Client struct {
	ctx    context.Context
	config *config.C
	client *apiv2.RESTClient
}

var APIClient *Client

func NewClient(ctx context.Context, conf *config.C) (*Client, error) {
	c, err := apiv2.NewRESTClientForHost(conf.Host, conf.Auth.User, conf.Auth.Password, &client_cf.Options{
		PageSize: -1,
		Page:     -1,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		ctx:    ctx,
		config: conf,
		client: c,
	}, nil
}

func (c *Client) AllProjects(name string) ([]*model.Project, error) {
	return c.client.ListProjects(c.ctx, name)
}

func (c *Client) ListAllRepositories(project string) ([]*model.Repository, error) {
	return c.client.ListRepositories(c.ctx, project)
}

func (c *Client) ListTags(projectName string, repoName string) ([]*Tag, error) {
	artifacts, err := c.client.ListArtifacts(c.ctx, projectName, repoName)
	if err != nil {
		return nil, err
	}

	tags := make([]*Tag, 0)

	for _, arti := range artifacts {
		for _, tag := range arti.Tags {
			tags = append(tags, &Tag{
				Digest:  arti.Digest,
				Name:    tag.Name,
				Created: time.Time(tag.PushTime),
			})
		}
	}

	return tags, nil
}

func (c *Client) DeleteTag(projectName, repoName, reference, tagName string) error {
	if err := c.client.DeleteTag(c.ctx, projectName, repoName, reference, tagName); err != nil {
		return err
	}

	art, err := c.client.GetArtifact(c.ctx, projectName, repoName, reference)
	if err != nil {
		return err
	}
	// 只有一个 tag 存在，则直接删除该 digest
	if len(art.Tags) == 0 {
		return c.client.DeleteArtifact(c.ctx, projectName, repoName, reference)
	}

	return nil
}
