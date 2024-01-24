// Client for interacting with Fastly NGWAF API
package ngwaf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const apiURL = "https://dashboard.signalsciences.net/api"

// Client is the API client
type Client struct {
	email string
	token string
	corp  string
}

// NewTokenClient creates a Client using token authentication
func NewTokenClient(email, token, corp string) Client {
	return Client{
		email: email,
		token: token,
		corp:  corp,
	}
}

func (fc *Client) doRequest(method, url, reqBody string) ([]byte, error) {
	client := &http.Client{}

	var b io.Reader
	if reqBody != "" {
		b = strings.NewReader(reqBody)
	}

	req, err := http.NewRequest(method, apiURL+url, b)
	if err != nil {
		return []byte{}, err
	}

	if fc.email != "" {
		// token auth
		req.Header.Set("X-API-User", fc.email)
		req.Header.Set("X-API-Token", fc.token)
	} else {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", fc.token))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go-sigsci")

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	switch method {
	case "GET":
		if resp.StatusCode != http.StatusOK {
			return body, errMsg(body)
		}
	case "POST":
		switch resp.StatusCode {
		case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		default:
			return body, errMsg(body)
		}
	case "DELETE":
		if resp.StatusCode != http.StatusNoContent {
			return body, errMsg(body)
		}
	case "PATCH":
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return body, errMsg(body)
		}
	}

	return body, nil
}

func errMsg(b []byte) error {
	var errResp struct {
		Message string
	}

	err := json.Unmarshal(b, &errResp)
	if err != nil {
		return err
	}

	return errors.New(errResp.Message)
}

// Simulation request and response for simulation test
type Simulation struct {
	SampleRequest  string `json:"sample_request"`
	SampleResponse string `json:"sample_response"`
}

// SimulationResponse the response of the simulation test
type SimulationOutput struct {
	Data struct {
		WafResponse int      `json:"waf_response"`
		Signals     []Signal `json:"signals"`
	} `json:"data"`
}

// Signal the signal from the simulation test
type Signal struct {
	Type      string `json:"type"`
	Location  string `json:"location"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	Detector  string `json:"detector"`
	Redaction int    `json:"redaction"`
}

// getSimulationOutput gets the simulation output
func getSimulationOutput(response []byte) (SimulationOutput, error) {
	var simulationOutput SimulationOutput
	err := json.Unmarshal(response, &simulationOutput)
	if err != nil {
		return SimulationOutput{}, err
	}
	return simulationOutput, nil
}

// SimulationTest runs a simulation test
func (sc *Client) SimulationTest(sim Simulation, site string) (SimulationOutput, error) {
	b, err := json.Marshal(sim)
	if err != nil {
		return SimulationOutput{}, err
	}
	resp, err := sc.doRequest("POST", fmt.Sprintf("/v0/corps/%s/sites/%s/simulator", sc.corp, site), string(b))
	if err != nil {
		return SimulationOutput{}, err
	}
	return getSimulationOutput(resp)
}
