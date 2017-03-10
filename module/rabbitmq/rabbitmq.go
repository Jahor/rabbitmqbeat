package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseUrl string
	http    *http.Client
	auth    *auth
}

type auth struct {
	username string
	password string
}

type Api interface {
	Nodes() ([]common.MapStr, error)
	Overview() (common.MapStr, error)
	Queues() ([]common.MapStr, error)
	Connections() ([]common.MapStr, error)
}

func NewClient(host, username, password string, tls *outputs.TLSConfig) (Api, error) {
	var transport http.RoundTripper
	var url string
	host = strings.TrimPrefix(strings.TrimPrefix(host, "http://"), "https://")
	if tls != nil && tls.IsEnabled() {
		cfg, err := outputs.LoadTLSConfig(tls)

		if err != nil {
			logp.Err("Can not load SSL config", err)
			return nil, err
		}

		ssl := cfg.BuildModuleConfig(strings.Split(host, ":")[0])

		transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       ssl,
		}
		url = "https://" + host
	} else {
		url = "http://" + host
		transport = http.DefaultTransport
	}

	c := Client{
		baseUrl: url,
		http: &http.Client{
			Transport: transport,
		},
		auth: &auth{
			username: username,
			password: password,
		},
	}
	return Api(&c), nil
}

func (c *Client) Nodes() ([]common.MapStr, error) {
	n, err := c.getMany("nodes")
	return n, err
}

func (c *Client) Queues() ([]common.MapStr, error) {
	q, err := c.getMany("queues")
	return q, err
}

func (c *Client) Connections() ([]common.MapStr, error) {
	q, err := c.getMany("connections")
	return q, err
}

func (c *Client) Overview() (common.MapStr, error) {
	o, err := c.getOne("overview")
	return o, err
}

func (c *Client) getOne(call string) (common.MapStr, error) {
	var q common.MapStr
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%v", c.baseUrl, call), nil)
	if err != nil {
		return q, err
	}
	req.SetBasicAuth(c.auth.username, c.auth.password)
	resp, err := c.http.Do(req)
	if err != nil {
		return q, err
	}
	if resp.StatusCode != 200 {
		return q, fmt.Errorf(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return q, err
	}
	err = json.Unmarshal(body, &q)
	return q, err
}

func (c *Client) getMany(call string) ([]common.MapStr, error) {
	var q []common.MapStr
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%v", c.baseUrl, call), nil)
	if err != nil {
		return q, err
	}
	req.SetBasicAuth(c.auth.username, c.auth.password)
	resp, err := c.http.Do(req)
	if err != nil {
		return q, err
	}
	if resp.StatusCode != 200 {
		return q, fmt.Errorf(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return q, err
	}
	err = json.Unmarshal(body, &q)
	return q, err
}
