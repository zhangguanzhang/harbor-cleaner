package cleaner

import (
	"fmt"
	//"github.com/distribution/distribution/manifest/schema1"
	//"github.com/distribution/distribution/manifest/schema2"
	"github.com/sirupsen/logrus"

	"github.com/goharbor/harbor/src/pkg/registry"

	"github.com/zhangguanzhang/harbor-cleaner/pkg/harbor"
	"github.com/zhangguanzhang/harbor-cleaner/pkg/policy"
)

//var acceptMediaTypes = []string{schema1.MediaTypeManifest, schema2.MediaTypeManifest}

type RepoCleaner struct {
	candidate  *policy.Candidate
	client     *harbor.Client
	repoClient registry.Client
	protected  []protectedTagsMenifest
}

type protectedTagsMenifest struct {
	tags      []string
	mediaType string
	payload   []byte
}

func NewRepoCleaner(candidate *policy.Candidate, client *harbor.Client) *RepoCleaner {
	return &RepoCleaner{
		candidate: candidate,
		client:    client,
	}
}

func (c *RepoCleaner) Protect() error {
	// v1 里如果两个tag是同一个 digest，删掉一个tag后，这俩tag都没，所以这里保留下 tag，下面的 Restore 里恢复
	// v2 里没这样的烦恼，但是删完 tag 后需要清理空的 digest
	//if len(c.candidate.Protected) > 0 {
	//	repoClient := harbor.NewRepoClient(c.client, fmt.Sprintf("%s/%s", c.candidate.Project, c.candidate.Repo))
	//	c.repoClient = repoClient
	//
	//	for digestID, tags := range c.candidate.Protected {
	//		manifest, digest, err := repoClient.PullManifest(c.client.digestID, acceptMediaTypes)
	//		if err != nil {
	//			logrus.Errorf("Pulling manifest %s/%s:%s error: %v", c.candidate.Project, c.candidate.Repo, digestID, err)
	//			return err
	//		}
	//		if strings.Contains(mediaType, "application/json") {
	//			mediaType = schema1.MediaTypeManifest
	//		}
	//
	//		c.protected = append(c.protected, protectedTagsMenifest{
	//			tags:      tags,
	//			mediaType: mediaType,
	//			payload:   payload,
	//		})
	//	}
	//}

	return nil
}

func (c *RepoCleaner) Clean() (int, error) {
	count := 0
	for _, tag := range c.candidate.Tags {
		if err := c.client.DeleteTag(c.candidate.Project, c.candidate.Repo, tag.Digest, tag.Name); err != nil {
			logrus.Warningf("Clean image '%s' error: %v", fmt.Sprintf("%s/%s:%s", c.candidate.Project, c.candidate.Repo, tag), err)
		} else {
			count++
		}
	}

	return count, nil
}

func (c *RepoCleaner) Restore() error {
	//for _, r := range c.protected {
	//	logrus.Infof("Start to push back tags %v to %s/%s", r.tags, c.candidate.Project, c.candidate.Repo)
	//	for _, t := range r.tags {
	//		_, err := c.repoClient.PushManifest(t, r.mediaType, r.payload)
	//		if err != nil {
	//			logrus.Errorf("Push manifest %s/%s:%s error: %v", c.candidate.Project, c.candidate.Repo, t)
	//			return err
	//		}
	//	}
	//}

	return nil
}
