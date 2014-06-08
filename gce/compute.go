package gce

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"

	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/compute/v1"
)

type gcloudCredentialsCache struct {
	Data []gceConfig
}

type gceConfig struct {
	Credential gceCredential
	Key        gceKey
	ProjectId  string `json:"projectId"`
}

type gceCredential struct {
	ClientId     string `json:"Client_Id"`
	ClientSecret string `json:"Client_Secret"`
	RefreshToken string `json:"Refresh_Token"`
}

type gceKey struct {
	Scope string
}

// Gets the OAuth2 token from the current user's gcloud crendentials.
func getOauthToken() (*gceCredential, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	confPath := path.Join(usr.HomeDir, ".config/gcloud/credentials")
	f, err := os.Open(confPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load gcloud credentials: %q", confPath)
	}
	defer f.Close()
	cache := &gcloudCredentialsCache{}
	if err := json.NewDecoder(f).Decode(cache); err != nil {
		return nil, err
	}
	if len(cache.Data) == 0 {
		return nil, fmt.Errorf("no gcloud credentials cached in: %q", confPath)
	}
	return &cache.Data[0].Credential, nil
}

// Gets the OAuth2 token from the current user's gcloud crendentials.
func getProjectId() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	confPath := path.Join(usr.HomeDir, ".config/gcloud/credentials")
	f, err := os.Open(confPath)
	if err != nil {
		return "", fmt.Errorf("unable to load gcloud credentials: %q", confPath)
	}
	defer f.Close()
	cache := &gcloudCredentialsCache{}
	if err := json.NewDecoder(f).Decode(cache); err != nil {
		return "", err
	}
	if len(cache.Data) == 0 {
		return "", fmt.Errorf("no gcloud credentials cached in: %q", confPath)
	}
	return cache.Data[0].ProjectId, nil
}

func NewCompute() (*compute.Service, error) {
	// Get Oauth2 token.
	creds, err := getOauthToken()
	if err != nil {
		log.Fatal(err)
	}

	// OAuth2 configuration.
	oAuth2Conf := &oauth.Config{
		ClientId:     creds.ClientId,
		ClientSecret: creds.ClientSecret,
		Scope:        compute.ComputeScope,
		RedirectURL:  "oob",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		AccessType:   "offline",
	}
	transport := &oauth.Transport{
		Config: oAuth2Conf,
		// Make the actual request using the cached token to authenticate.
		Token:     &oauth.Token{RefreshToken: creds.RefreshToken},
		Transport: http.DefaultTransport,
	}

	err = transport.Refresh()
	if err != nil {
		return nil, err
	}
	svc, err := compute.New(transport.Client())
	if err != nil {
		return nil, err
	}
	return svc, nil
}

type gceVmManager struct {
	projectId string
}

func NewGceManager() (VirtualMachineManager, error) {
	projectId, err := getProjectId()
	if err != nil {
		return nil, fmt.Errorf("unable to get project id: %v", err)
	}
	ret := &gceVmManager{
		projectId: projectId,
	}
	return ret, nil
}

func getMachineType(spec *VirtualMachineSpec) string {
	return "n1-standard-2"
}

func getZone(spec *VirtualMachineSpec) string {
	return "us-central1-a"
}

func (self *gceVmManager) NewMachine(spec *VirtualMachineSpec) (*VirtualMachineInfo, error) {
	service, err := NewCompute()
	if err != nil {
		return nil, fmt.Errorf("unable to get compute service: %v", err)
	}
	zone := getZone(spec)
	prefix := "https://www.googleapis.com/compute/v1/projects" + self.projectId
	machineType := getMachineType(spec)
	instance := &compute.Instance{
		Name:        spec.GetName(),
		Description: "virtigo instance",
		Zone:        fmt.Sprintf("%v/zones/%v", prefix, zone),
		MachineType: fmt.Sprintf("%v/machine-types/%v", machineType),
		NetworkInterfaces: []*compute.NetworkInterface{
			&compute.NetworkInterface{
				AccessConfigs: []*compute.AccessConfig{
					&compute.AccessConfig{Type: "ONE_TO_ONE_NAT"},
				},
				Network: prefix + "/networks/default",
			},
		},
	}

	opt, err := service.Instances.Insert(self.projectId, zone, instance).Do()

	if err != nil {
		return nil, fmt.Errorf("unable to create vm: %v", err)
	}
	info := &VirtualMachineInfo{
		Name: opt.Name,
	}
	return info, nil
}
