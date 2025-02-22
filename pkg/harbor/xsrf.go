package harbor

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/zhangguanzhang/harbor-cleaner/pkg/config"
)

func SetXSRFToken(client *http.Client, conf *config.C, targetReq *http.Request) error {
	if !conf.XSRF.Enabled {
		return nil
	}

	req, err := http.NewRequest(http.MethodGet, PingURL(conf.Host), nil)
	if err != nil {
		logrus.Error(err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		logrus.Errorf("ping Harbor: %s error: %s", conf.Host, b)
		return fmt.Errorf("%s", b)
	}

	for _, c := range resp.Cookies() {
		if c.Name == "_xsrf" {
			if v, ok := GetSecureCookie(conf.XSRF.Key, c.Value); ok {
				targetReq.Header.Add("X-Xsrftoken", v)
			}
		}

		targetReq.AddCookie(c)
	}

	return nil
}

func GetSecureCookie(Secret, val string) (string, bool) {
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)

	if len(parts) != 3 {
		return "", false
	}

	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	h := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(h, "%s%s", vs, timestamp)

	if fmt.Sprintf("%02x", h.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}
