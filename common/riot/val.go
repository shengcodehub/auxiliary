package val

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gonxtserver/common/types"
	consts "gonxtserver/consts/val"
	"io"
	"net/http"
)

var riotValClient *RiotValClient

func RiotClientSetUp(apiKey string, region consts.Region) {
	// 创建一个新的 Discord 会话
	innerClient := InnerHttpClient{apiKey, region}
	riotValClient = &RiotValClient{innerClient}
}

func GetRiotValClient() *RiotValClient {
	return riotValClient
}

type RiotValClient struct {
	innerClient InnerHttpClient
}

const (
	apiURLFormat      = "%s://%s.%s%s"
	baseURL           = "api.riotgames.com"
	scheme            = "https"
	apiTokenHeaderKey = "X-Riot-Token"
)

const (
	endpointBase                  = "/val"
	endpointContentBase           = endpointBase + "/content/v1"
	endPointGetContent            = endpointContentBase + "/contents?%s"
	endpointStatusBase            = endpointBase + "/status/v1"
	endpointGetPlatformData       = endpointStatusBase + "/platform-data"
	endpointRankedBase            = endpointBase + "/ranked/v1"
	endpointGetLeaderboardByActID = endpointRankedBase + "/leaderboards/by-act/%s"
	endpointMatchBase             = endpointBase + "/match/v1"
	endpointMatchByID             = endpointMatchBase + "/matches/%s"
	endpointMatchListByPUUID      = endpointMatchBase + "/matchlists/by-puuid/%s"
	endpointRecentMatchesByQueue  = endpointMatchBase + "/recent-matches/by-queue/%s"
)

const (
	endpointRiotBase              = "/riot"
	endpointAccountBase           = endpointRiotBase + "/account/v1"
	endAccountActiveRegionByPUUID = endpointAccountBase + "/active-shards/by-game/val/by-puuid/%s"
)

func (rc *RiotValClient) GetAccountActiveRegionByPUUID(ctx context.Context, puuid string) (*types.ValAccountActiveShardInfo, error) {
	logx.Infof("GetMatchListByPUUID: %s", puuid)
	url := endAccountActiveRegionByPUUID
	var activeShard *types.ValAccountActiveShardInfo

	if err := rc.innerClient.GetInto(ctx, fmt.Sprintf(url, puuid), consts.Region(consts.RouteAsia), &activeShard); err != nil {
		logx.Errorf("RiotValClient,GetAccountActiveRegionByPUUID,puuid:%s, fail:%s", puuid, err.Error())
		return nil, err
	}
	return activeShard, nil
}

func (rc *RiotValClient) GetMatchListByPUUID(ctx context.Context, region consts.Region, puuid string) (*types.ValMatchList, error) {
	logx.Infof("GetMatchListByPUUID: %s", puuid)
	url := endpointMatchListByPUUID
	var matchList *types.ValMatchList
	if err := rc.innerClient.GetInto(ctx, fmt.Sprintf(url, puuid), region, &matchList); err != nil {
		logx.Infof("RiotValClient,GetMatchListByPUUID,puuid:%s, fail:%s", puuid, err.Error())
		return nil, err
	}
	return matchList, nil
}

func (rc *RiotValClient) GetMatchByID(ctx context.Context, region consts.Region, matchID string) (*types.ValMatch, error) {
	logx.Infof("GetMatchByID: %s", matchID)
	url := endpointMatchByID
	var match *types.ValMatch
	if err := rc.innerClient.GetInto(ctx, fmt.Sprintf(url, matchID), region, &match); err != nil {
		logx.Infof("RiotValClient,GetMatchByID,matchID:%s, fail:%s", matchID, err.Error())
		return nil, err
	}
	return match, nil
}

type InnerHttpClient struct {
	APIKey string
	Region consts.Region
}

func (c *InnerHttpClient) GetInto(ctx context.Context, endpoint string, region consts.Region, target interface{}) error {

	response, err := c.Get(ctx, endpoint, region)
	if err != nil {
		return err
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return err
	}
	return nil
}

// PostInto processes a POST request and saves the response body into the given target.
func (c *InnerHttpClient) PostInto(ctx context.Context, endpoint string, region consts.Region, body, target interface{}) error {
	response, err := c.Post(ctx, endpoint, region, body)
	if err != nil {
		return err
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return err
	}
	return nil
}

// Put processes a PUT request.
func (c *InnerHttpClient) Put(ctx context.Context, endpoint string, region consts.Region, body interface{}) error {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return err
	}
	_, err := c.DoRequest(ctx, "PUT", endpoint, region, buf)
	return err
}

// Get processes a GET request.
func (c *InnerHttpClient) Get(ctx context.Context, endpoint string, region consts.Region) (*http.Response, error) {
	return c.DoRequest(ctx, "GET", endpoint, region, nil)
}

// Post processes a POST request.
func (c *InnerHttpClient) Post(ctx context.Context, endpoint string, region consts.Region, body interface{}) (*http.Response, error) {

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}
	return c.DoRequest(ctx, "POST", endpoint, region, buf)
}

// DoRequest processes a http.Request and returns the response.
// Rate-Limiting and retrying is handled via the corresponding response headers.
func (c *InnerHttpClient) DoRequest(ctx context.Context, method string, endpoint string, region consts.Region, body io.Reader) (*http.Response, error) {

	request, err := c.NewRequest(ctx, method, endpoint, region, body)
	if err != nil {
		return nil, err
	}
	//response, err := c.Do(request)
	httpclient := http.DefaultClient
	response, err := httpclient.Do(request)
	if err != nil {
		logx.Errorf("InnerHttpClient error request, err:%v", err)
		return nil, err
	}
	if response.StatusCode == http.StatusServiceUnavailable {
		logx.Infof("service unavailable,request riot api: %s", endpoint)
		return nil, errors.New("riot service unavailable")
	}
	if response.StatusCode == http.StatusTooManyRequests {

		return nil, ErrRateLimitExceeded
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		//logx.Infof("response: %v", response.Status)
		err, ok := StatusToError[response.StatusCode]
		if !ok {
			err = Error{
				Message:    "unknown err reason",
				StatusCode: response.StatusCode,
			}
		}
		return nil, err
	}
	return response, nil
}

// NewRequest returns a new http.Request with necessary headers et.
func (c *InnerHttpClient) NewRequest(ctx context.Context, method string, endpoint string, region consts.Region, body io.Reader) (*http.Request, error) {

	request, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf(apiURLFormat, scheme, region, baseURL, endpoint), body)
	if err != nil {
		return nil, err
	}
	request.Header.Add(apiTokenHeaderKey, c.APIKey)
	request.Header.Add("Accept", "application/json")
	return request, nil
}
