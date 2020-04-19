package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Client is skyway-gateway client
type Client struct {
	apiKey     string
	gwHost     string
	gwPort     int
	httpClient *http.Client
}

// NewClient generates client.
// APIKey is obtained by skyway.
// gwHost and gwPort is your skyway-gateway host and port
func NewClient(apiKey, gwHost string, gwPort int) *Client {
	return &Client{
		apiKey: apiKey,
		gwHost: gwHost,
		gwPort: gwPort,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Constraints is used for media connection constraints.
type Constraints struct {
	Video               bool `json:"video"`
	VideoReceiveEnabled bool `json:"videoReceiveEnabled"`
	Audio               bool `json:"audio"`
	AudioReceiveEnabled bool `json:"audioReceiveEnabled"`
	VideoParams         struct {
		BandWidth   int    `json:"band_width"`
		Codec       string `json:"codec"`
		MediaID     string `json:"media_id"`
		RTCPID      string `json:"rtcp_id"`
		PayloadType int    `json:"payload_type"`
	} `json:"video_params"`
	AudioParams struct {
		BandWidth   int    `json:"band_width"`
		Codec       string `json:"codec"`
		MediaID     string `json:"media_id"`
		RTCPID      string `json:"rtcp_id"`
		PayloadType int    `json:"payload_type"`
	} `json:"audio_params"`
}

// MediaRedirect has IP and Port information
type MediaRedirect struct {
	Video struct {
		IPv4 string `json:"ip_v4"`
		Port int    `json:"port"`
	} `json:"video"`
	Audio struct {
		IPv4 string `json:"ip_v4"`
		Port int    `json:"port"`
	} `json:"audio"`
	VideoRTCP struct {
		IPv4 string `json:"ip_v4"`
		Port int    `json:"port"`
	} `json:"video_rtcp"`
	AudioRTCP struct {
		IPv4 string `json:"ip_v4"`
		Port int    `json:"port"`
	} `json:"audio_rtcp"`
}

// CreatPeer returns token from gateway and error
func (c *Client) CreatPeer(peerID, domain string) (string, error) {
	type params struct {
		Key    string `json:"key"`
		Domain string `json:"domain"`
		Turn   bool   `json:"turn"`
		PeerID string `json:"peer_id"`
	}

	type responseData struct {
		CommandType string `json:"command_type"`
		Params      struct {
			PeerID string `json:"peer_id"`
			Token  string `json:"token"`
		} `json:"params"`
	}

	p := params{
		Key:    c.apiKey,
		Domain: domain,
		Turn:   true,
		PeerID: peerID,
	}
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(p); err != nil {
		return "", err
	}

	resp, err := c.request(http.MethodPost, "peers", reqBody)

	respData := &responseData{}
	if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return respData.Params.Token, nil
}

// CreateMedia returns media id, IPv4 addr, udp port, and error
func (c *Client) CreateMedia(isVideo bool) (string, string, int, error) {
	type params struct {
		IsVideo bool `json:"is_video"`
	}
	type responseData struct {
		MediaID string `json:"media_id"`
		Port    int    `json:"port"`
		IPv4    string `json:"ip_v4"`
		IPv6    string `json:"ip_v6"`
	}

	p := params{isVideo}
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(p); err != nil {
		return "", "", -1, err
	}

	resp, err := c.request(http.MethodPost, "media", reqBody)

	respData := &responseData{}
	if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
		return "", "", -1, err
	}
	if err != nil {
		return "", "", -1, err
	}
	return respData.MediaID, respData.IPv4, respData.Port, nil
}

// CreateRTCP returns rtcp id, IPv4 addr, udp port for rtcp, and error
func (c *Client) CreateRTCP() (string, string, int, error) {
	type params struct{}
	type responseData struct {
		RTCPID string `json:"rtcp_id"`
		Port   int    `json:"port"`
		IPv4   string `json:"ip_v4"`
		IPv6   string `json:"ip_v6"`
	}

	p := params{}
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(p); err != nil {
		return "", "", -1, err
	}

	resp, err := c.request(http.MethodPost, "media/rtcp", reqBody)

	respData := &responseData{}
	if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
		return "", "", -1, err
	}
	if err != nil {
		return "", "", -1, err
	}
	return respData.RTCPID, respData.IPv4, respData.Port, nil
}

// Call returns media connection id and error
func (c *Client) Call(token, peerID, targetID string, constraints Constraints, redirect MediaRedirect) (string, error) {
	type params struct {
		PeerID         string        `json:"peer_id"`
		Token          string        `json:"token"`
		TargetID       string        `json:"target_id"`
		Constraints    Constraints   `json:"constraints"`
		RedirectParams MediaRedirect `json:"redirect_params"`
	}
	type responseData struct {
		CommandType string `json:"command_type"`
		Params      struct {
			MediaConnectionID string `json:"media_connection_id"`
		} `json:"params"`
	}

	p := params{
		PeerID:         peerID,
		Token:          token,
		TargetID:       targetID,
		Constraints:    constraints,
		RedirectParams: redirect,
	}
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(p); err != nil {
		return "", err
	}
	resp, err := c.request(http.MethodPost, "media/connections", reqBody)
	if err != nil {
		return "", err
	}
	respData := &responseData{}
	if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
		return "", err
	}
	fmt.Println(respData)
	return respData.Params.MediaConnectionID, nil
}

// WaitCall waits call from other peer
func (c *Client) WaitCall(peerID, token string) (string, error) {
	type responseData struct {
		Event  string `json:"event"`
		Params struct {
			PeerID string `json:"peer_id"`
			Token  string `json:"token"`
		} `json:"params"`
		CallParams struct {
			MediaConnectionID string `json:"media_connection_id"`
		} `json:"call_params"`
		DataParams struct {
			DataConnectionID string `json:"data_connection_id"`
		} `json:"data_params"`
		ErrorMessage string `json:"error_message"`
	}
	expectedEvent := "CALL"
	event := ""
	respData := &responseData{}
	for event != expectedEvent {
		log.Printf("peers/%s/events?token=%s\n", peerID, token)
		resp, err := c.request(http.MethodGet, fmt.Sprintf("peers/%s/events?token=%s", peerID, token), nil)
		if err != nil {
			return "", err
		}

		var r io.Reader = resp.Body
		r = io.TeeReader(r, os.Stderr)
		if resp.StatusCode == 408 {
			continue
		}
		if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
			return "", err
		}
		event = respData.Event
	}
	return respData.CallParams.MediaConnectionID, nil
}

// Answer apply call from other peer
func (c *Client) Answer(mediaConnectionID string, constraints Constraints, redirect MediaRedirect) error {
	type params struct {
		Constraints    Constraints   `json:"constraints"`
		RedirectParams MediaRedirect `json:"redirect_params"`
	}
	type responseData struct {
		CommandType string `json:"command_type"`
		Params      struct {
			VideoID string `json:"video_id"`
			AudioID string `json:"audio_id"`
		} `json:"params"`
	}

	p := params{
		Constraints:    constraints,
		RedirectParams: redirect,
	}
	reqBody := new(bytes.Buffer)
	if err := json.NewEncoder(reqBody).Encode(p); err != nil {
		return err
	}
	resp, err := c.request(http.MethodPost, fmt.Sprintf("media/connections/%s/answer", mediaConnectionID), reqBody)
	if err != nil {
		return err
	}
	respData := &responseData{}
	if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
		return err
	}

	return nil
}

// WatchMediaConnection returns error
func (c *Client) WatchMediaConnection(connectionID string) error {
	type responseData struct {
		Event         string `json:"event"`
		StreamOptions struct {
			IsVideo      bool `json:"is_video"`
			StreamParams struct {
				MediaID string `json:"media_id"`
				Port    int    `json:"port"`
				IPv4    string `json:"ip_v4"`
				IPv6    string `json:"ip_v6"`
			} `json:"stream_params"`
		} `json:"stream_options"`
		CloseOptions struct{} `json:"close_options"`
		ErrorMessage string   `json:"error_message"`
	}
	expectedEvent := "OPEN"
	event := ""
	for event != expectedEvent {
		log.Printf("media/connections/%s/events", connectionID)
		resp, err := c.request(http.MethodGet, fmt.Sprintf("media/connections/%s/events", connectionID), nil)
		if err != nil {
			return err
		}
		respData := &responseData{}
		if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
			return err
		}
		event = respData.Event
	}
	return nil
}

// PeerClose disconnects peer
func (c *Client) PeerClose(peerID, token string) (*http.Response, error) {
	return c.request(http.MethodDelete, fmt.Sprintf("peers/%s?token=%s", peerID, token), nil)
}

// MediaConnectionClose disconnects media connection
func (c *Client) MediaConnectionClose(connectionID string) (*http.Response, error) {
	return c.request(http.MethodDelete, fmt.Sprintf("media/connections/%s", connectionID), nil)
}

func (c *Client) request(method, path string, body io.Reader) (*http.Response, error) {
	time.Sleep(1 * time.Second)
	req, err := http.NewRequest(method, fmt.Sprintf("%s:%d/%s", c.gwHost, c.gwPort, path), body)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 300 {
		fmt.Printf("Request %s:%d/%s\n", c.gwHost, c.gwPort, path)
		return nil, fmt.Errorf("get response statud code: %d", resp.StatusCode)
	}
	return resp, nil
}
