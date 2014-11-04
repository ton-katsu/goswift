package goswift

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const userAgent = "goswift/0.1"

// Keystone request json
type KeystoneV2Req struct {
	Auth Credentials `json:"auth"`
}

type Credentials struct {
	PasswordCredentials UserPass `json:"passwordCredentials"`
	TenantName          string   `json:"tenantName"`
}

type UserPass struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Keystoen respose json
type KeystoneV2Res struct {
	Access AccessInfo `json:"access"`
}

type AccessInfo struct {
	Metadata       MetadataInfo      `json:"metadata"`
	User           UserInfo          `json:"user"`
	ServiceCatalog []ServiceCatalogs `josn:"serviceCatalog"`
	Token          TokenInfo         `json:"token"`
}

type MetadataInfo struct {
	Roles   []string `json:"roles"`
	IsAdmin uint     `json:"is_admin"`
}

type UserInfo struct {
	Name        string      `json:"name"`
	Roles       []RolesInfo `json:"roles"`
	Id          string      `json:"id"`
	RolesLinkes []string    `json:"roles_links"`
	Username    string      `json:"username"`
}

type RolesInfo struct {
	Name string `json:"name"`
}

type ServiceCatalogs struct {
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	EndpointsLinks []string        `json:"endpoints_links"`
	Endpoints      []EndpointsInfo `json:"endpoints"`
}

type EndpointsInfo struct {
	PublicUrl   string `json:"publicURL"`
	Id          string `json:"id"`
	InternalURL string `json:"internalURL"`
	Region      string `json:"region"`
	AdminUrl    string `json:"adminURL"`
}

type TokenInfo struct {
	Tenant   TenantInfo `json:"tenant"`
	Id       string     `json:"id"`
	Expires  string     `json:"expires"`
	IssuedAt string     `json:"issued_at"`
}

type TenantInfo struct {
	Name        string `json:"name"`
	Id          string `json:"id"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type Client struct {
	Client      *http.Client
	AccountName string
	Password    string
	AuthUrl     string
	StorageUrl  string
	Token       string
	TenantName  string
	RegionName  string
	SkipSecure  bool
	ChunkSize   uint
}

func (c *Client) SWAuthV1() error {
	req, _ := http.NewRequest("GET", c.AuthUrl, nil)
	req.Header.Set("X-Auth-User", c.AccountName)
	req.Header.Set("X-Auth-Key", c.Password)
	req.Header.Set("User-Agent", userAgent)
	res, err := c.Client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if err := CheckResponse(res); err != nil {
		return err
	}
	c.StorageUrl = res.Header["X-Storage-Url"][0]
	c.Token = res.Header["X-Auth-Token"][0]
	return err
}

func (c *Client) KeystoneAuthV2() error {
	a := KeystoneV2Req{Credentials{UserPass{c.AccountName, c.Password}, c.TenantName}}
	b, _ := json.Marshal(a)
	req, _ := http.NewRequest("POST", c.AuthUrl, bytes.NewReader(b))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if err := CheckResponse(res); err != nil {
		return err
	}

	resbody, _ := ioutil.ReadAll(res.Body)
	var keystonev2res KeystoneV2Res
	err = json.Unmarshal(resbody, &keystonev2res)
	if err != nil {
		fmt.Println(err)
	}
	catalogList := keystonev2res.Access.ServiceCatalog
	var publicUrl string
	for i := range catalogList {
		if catalogList[i].Name == "swift" || catalogList[i].Type == "object-store" {
			for e := range catalogList[i].Endpoints {
				if catalogList[i].Endpoints[e].Region == c.RegionName {
					publicUrl = catalogList[i].Endpoints[e].PublicUrl
				}
			}
		}
	}
	c.StorageUrl = publicUrl
	c.Token = keystonev2res.Access.Token.Id
	return err
}

func (c *Client) Credential() (err error) {
	c.setClient()
	u, _ := url.Parse(strings.Trim(c.AuthUrl, "/"))
	ksver := strings.Split(u.Path, "/")[1]
	swauthver := strings.Split(u.Path, "/")[2]
	switch {
	case swauthver == "v1.0":
		err = c.SWAuthV1()
	case ksver == "v2" || ksver == "v2.0":
		err = c.KeystoneAuthV2()
	default:
		err = errors.New("Check the API version. Support to v1 or v2.")
	}
	return err
}

func (c *Client) setHeaders(req *http.Request, header http.Header) *http.Request {
	req.Header.Set("X-Auth-Token", c.Token)
	req.Header.Set("User-Agent", userAgent)
	if len(header) != 0 {
		for k, v := range header {
			req.Header.Set(k, fmt.Sprintf("%v", v[0]))
		}
	}
	return req
}

func (c *Client) setClient() {
	if c.Client == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.SkipSecure},
		}
		c.Client = &http.Client{Transport: tr}
	}
}

func (c *Client) setCredential() (err error) {
	if c.AuthUrl != "" && c.AccountName != "" && c.Password != "" {
		if c.Token == "" && c.StorageUrl == "" {
			err = c.Credential()
		}
	}
	if c.Token == "" || c.StorageUrl == "" {
		return errors.New("Check the params.")
	}
	return err
}

func (c *Client) request(method string, path string, body io.Reader, contentLength int64, header http.Header, params url.Values) ([]byte, map[string][]string, error) {
	c.setClient()
	if err := c.setCredential(); err != nil {
		return nil, nil, err
	}
	urls := fmt.Sprintf("%s/%s", strings.Trim(c.StorageUrl, "/"), path)
	if params == nil {
		params = make(url.Values)
	}
	params.Set("format", "json")
	urls += "?" + params.Encode()
	req, _ := http.NewRequest(method, urls, body)
	req = c.setHeaders(req, header)
	req.ContentLength = contentLength
	res, err := c.Client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	if err := CheckResponse(res); err != nil {
		return nil, nil, err
	}
	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	return resbody, res.Header, err
}

type Params struct {
	Limit     uint
	Marker    string
	Endmarker string
	Prefix    string
	Delimiter string
	Path      string
}

func (p *Params) setQueryParams(params url.Values) {
	if p.Limit != 0 {
		params.Set("limit", fmt.Sprintf("%v", p.Limit))
	}
	if p.Marker != "" {
		params.Set("marker", fmt.Sprintf("%v", p.Marker))
	}
	if p.Endmarker != "" {
		params.Set("end_marker", fmt.Sprintf("%v", p.Endmarker))
	}
	if p.Prefix != "" {
		params.Set("prefix", fmt.Sprintf("%v", p.Prefix))
	}
	if p.Delimiter != "" {
		params.Set("delimiter", fmt.Sprintf("%v", p.Delimiter))
	}
	if p.Endmarker != "" {
		params.Set("path", fmt.Sprintf("%v", p.Path))
	}
}

type Metadata http.Header

func NewMetadata() Metadata {
	return make(Metadata)
}

func (m *Metadata) SetMeta(key, value string) {
	h := http.Header(*m)
	h.Set(key, value)
}

func (m *Metadata) SetDeleteMeta(key string) {
	h := http.Header(*m)
	var buffer bytes.Buffer
	split := strings.Split(strings.ToLower(key), "-")
	for j := range split {
		if split[j] == "x" {
			buffer.WriteString("X-Remove")
		} else {
			buffer.WriteString("-" + strings.Title(split[j]))
		}
	}
	h.Set(buffer.String(), "x")
}

// Accounts metadata operation
func (c *Client) ShowAccountMeta() (http.Header, error) {
	_, header, err := c.request("HEAD", "", nil, 0, nil, nil)
	return header, err
}

func (c *Client) CreateAccountMeta(metadata Metadata) (http.Header, error) {
	_, header, err := c.request("POST", "", nil, 0, http.Header(metadata), nil)
	return header, err
}

func (c *Client) UpdateAccountMeta(metadata Metadata) (http.Header, error) {
	return c.CreateAccountMeta(metadata)
}

func (c *Client) DeleteAccountMeta(metadata Metadata) (http.Header, error) {
	return c.CreateAccountMeta(metadata)
}

// Containers operation
type Container struct {
	Count uint
	Bytes uint
	Name  string
}

func (c *Client) ListContainers() ([]Container, http.Header, error) {
	return c.ListContainersWithParams(Params{})
}

func (c *Client) ListContainersWithParams(p Params) ([]Container, http.Header, error) {
	params := make(url.Values)
	p.setQueryParams(params)
	var container []Container
	body, header, err := c.request("GET", "", nil, 0, nil, params)
	if body != nil {
		json.Unmarshal(body, &container)
	}
	return container, header, err
}

func (c *Client) CreateContainer(containerName string) (http.Header, error) {
	_, header, err := c.request("PUT", containerName, nil, 0, nil, nil)
	return header, err
}

func (c *Client) DeleteContainer(containerName string) (http.Header, error) {
	_, header, err := c.request("DELETE", containerName, nil, 0, nil, nil)
	return header, err
}

// Containers metadata operation
func (c *Client) ShowContainerMeta(containerName string) (http.Header, error) {
	_, header, err := c.request("HEAD", containerName, nil, 0, nil, nil)
	return header, err
}

func (c *Client) CreateContainerMeta(containerName string, metadata Metadata) (http.Header, error) {
	_, header, err := c.request("POST", containerName, nil, 0, http.Header(metadata), nil)
	return header, err
}

func (c *Client) UpdateContainerMeta(containerName string, metadata Metadata) (http.Header, error) {
	return c.CreateContainerMeta(containerName, metadata)
}

func (c *Client) DeleteContainerMeta(containerName string, metadata Metadata) (http.Header, error) {
	return c.CreateContainerMeta(containerName, metadata)
}

// Objects operation
type Object struct {
	Hash         string
	LastModified string `json:"last_modified"`
	Bytes        uint
	Name         string
	ContentType  string `json:"content_type"`
}

func (c *Client) ListObjects(containerName string) ([]Object, http.Header, error) {
	return c.ListObjectsWithParams(containerName, Params{})
}

func (c *Client) ListObjectsWithParams(containerName string, p Params) ([]Object, http.Header, error) {
	params := make(url.Values)
	p.setQueryParams(params)
	var object []Object
	body, header, err := c.request("GET", containerName, nil, 0, nil, params)
	if body != nil {
		json.Unmarshal(body, &object)
	}
	return object, header, err
}

func (c *Client) GetObject(containerName string, objectName string) ([]byte, error) {
	objectPath := fmt.Sprintf("%s/%s", containerName, objectName)
	resbody, _, err := c.request("GET", objectPath, nil, 0, nil, nil)
	return resbody, err
}

func (c *Client) CreateObject(containerName string, objectName string, contentName string, metadata Metadata) (http.Header, error) {
	objectPath := fmt.Sprintf("%s/%s", containerName, objectName)
	fi, err := os.Stat(contentName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s", err))
	}
	contentLength := fi.Size()
	b, err := ioutil.ReadFile(contentName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s", err))
	}
	body := bytes.NewReader(b)
	// if c.ChunkSize != 0 {
	// 	return nil
	// }
	_, header, err := c.request("PUT", objectPath, body, contentLength, http.Header(metadata), nil)
	return header, err
}

func (c *Client) DeleteObject(containerName string, objectName string) error {
	objectPath := fmt.Sprintf("%s/%s", containerName, objectName)
	_, _, err := c.request("DELETE", objectPath, nil, 0, nil, nil)
	return err
}

func (c *Client) CopyObject(fromContainerName string, fromObjectName string, toContainerName string, toObjectName string) (http.Header, error) {
	toObjectPath := fmt.Sprintf("%s/%s", toContainerName, toObjectName)
	fromObjectPath := fmt.Sprintf("%s/%s", fromContainerName, fromObjectName)
	metadata := NewMetadata()
	metadata.SetMeta("Destination", toObjectPath)
	_, header, err := c.request("COPY", fromObjectPath, nil, 0, http.Header(metadata), nil)
	return header, err
}

// Objects metadata operation
func (c *Client) ShowObjectMeta(containerName string, objectName string) (http.Header, error) {
	objectPath := fmt.Sprintf("%s/%s", containerName, objectName)
	_, header, err := c.request("HEAD", objectPath, nil, 0, nil, nil)
	return header, err
}

func (c *Client) CreateObjectMeta(containerName string, objectName string, metadata Metadata) (http.Header, error) {
	objectPath := fmt.Sprintf("%s/%s", containerName, objectName)
	_, header, err := c.request("POST", objectPath, nil, 0, http.Header(metadata), nil)
	return header, err
}

func (c *Client) UpdateObjectMeta(containerName string, objectName string, metadata Metadata) (http.Header, error) {
	return c.CreateObjectMeta(containerName, objectName, metadata)
}

func (c *Client) DeleteObjectMeta(containerName string, objectName string, metadata Metadata) (http.Header, error) {
	return c.CreateObjectMeta(containerName, objectName, metadata)
}

// Error contains an error response from the server.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Body    string
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("Swift API: Error %d: %v", e.Code, e.Message)
	}
	return fmt.Sprintf("Swift API: got HTTP response code %d with body: %v", e.Code, e.Body)
}

type errorReply struct {
	Error *Error `json:"error"`
}

// CheckResponse returns an error (of type *Error) if the response
// status code is not 2xx.
func CheckResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	slurp, err := ioutil.ReadAll(res.Body)
	if err == nil {
		jerr := new(errorReply)
		err = json.Unmarshal(slurp, jerr)
		if err == nil && jerr.Error != nil {
			if jerr.Error.Code == 0 {
				jerr.Error.Code = res.StatusCode
			}
			jerr.Error.Body = string(slurp)
			return jerr.Error
		}
	}
	return &Error{
		Code: res.StatusCode,
		Body: string(slurp),
	}
}
