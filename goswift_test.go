package goswift

import (
	"io/ioutil"
	"os"
	"testing"
)

var AuthUrl = os.Getenv("SWIFT_AUTH_URL")
var AccountName = os.Getenv("SWIFT_API_USER")
var Password = os.Getenv("SWIFT_API_KEY")
var StorageUrl = os.Getenv("SWIFT_STORAGE_URL")
var Token = os.Getenv("SWIFT_USER_TOKEN")

// Account metadata operation
func TestShowAccountMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	header, err := c.ShowAccountMeta()
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log("X-Account-Bytes-Used", header["X-Account-Bytes-Used"])
	t.Log("X-Account-Container-Count", header["X-Account-Container-Count"])
}

func TestCreateAccountMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	metadata := NewMetadata()
	metadata.SetMeta("X-Account-Meta-Book", "saka01")
	metadata.SetMeta("X-Account-Meta-Subject", "saka02")
	header, err := c.CreateAccountMeta(metadata)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err = c.ShowAccountMeta()
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if header["X-Account-Meta-Book"][0] != "saka01" || header["X-Account-Meta-Subject"][0] != "saka02" {
		t.Errorf("Expected error: %s", "Could not set Metadata.")
	}
	t.Log("X-Account-Meta-Book", header["X-Account-Meta-Book"])
	t.Log("X-Account-Meta-Subject", header["X-Account-Meta-Subject"])
}

func TestUpdateAccountMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	metadata := NewMetadata()
	metadata.SetMeta("X-Account-Meta-Book", "saka10")
	metadata.SetMeta("X-Account-Meta-Subject", "saka20")
	header, err := c.UpdateAccountMeta(metadata)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err = c.ShowAccountMeta()
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if header["X-Account-Meta-Book"][0] != "saka10" || header["X-Account-Meta-Subject"][0] != "saka20" {
		t.Errorf("Expected error: %s", "Could not set Metadata.")
	}
	t.Log("X-Account-Meta-Book", header["X-Account-Meta-Book"])
	t.Log("X-Account-Meta-Subject", header["X-Account-Meta-Subject"])
}

func TestDeleteAccountMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	metadata := NewMetadata()
	metadata.SetDeleteMeta("X-Account-Meta-Book")
	metadata.SetDeleteMeta("X-Account-Meta-Subject")
	header, err := c.DeleteAccountMeta(metadata)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err = c.ShowAccountMeta()
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	// t.Log(header)
	if _, ok := header["X-Account-Meta-Book"]; ok != false {
		t.Errorf("Expected error: %s", "Could not delete Metadata.")
	}
	if _, ok := header["X-Account-Meta-Subject"]; ok != false {
		t.Errorf("Expected error: %s", "Could not delete Metadata.")
	}
	t.Log("X-Account-Meta-Book", header["X-Account-Meta-Book"])
	t.Log("X-Account-Meta-Subject", header["X-Account-Meta-Subject"])
}

// Container operation
func TestListContainersWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	res, _, err := c.ListContainers()
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log(len(res))
	t.Log(res[0].Name)
}

func TestListContainersWithStorageUrlAndToken(t *testing.T) {
	c := Client{StorageUrl: StorageUrl, Token: Token, SkipSecure: true}
	res, _, err := c.ListContainers()
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log(len(res))
	t.Log(res[0].Name)
}

func TestListContainersWithParamsAndInvalidToken(t *testing.T) {
	c := Client{StorageUrl: StorageUrl, Token: "invalid-token", SkipSecure: true}
	_, _, err := c.ListContainers()
	if err == nil {
		t.Errorf("Expected error: Does not invalid. %s", err)
	}
}

func TestListContainersWithParamsAndLimit4(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	res, _, err := c.ListContainersWithParams(Params{Limit: 1})
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected error: Limit count is different.")
	}
	t.Log(res)
}

func TestCreateContainersWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	header, err := c.CreateContainer("bb94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err = c.CreateContainer("9194ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log(header)
}

func TestDeleteContainersWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	header, err := c.DeleteContainer("9194ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log(header)
}

// Container metadata operation
func TestShowContainerMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	header, err := c.ShowContainerMeta("bb94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log("X-Container-Meta-Author", header["X-Container-Meta-Author"])
	t.Log("X-Container-Meta-Century", header["X-Container-Meta-Century"])
}

func TestCreateContainerMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	metadata := NewMetadata()
	metadata.SetMeta("X-Container-Meta-Author", "saka01")
	metadata.SetMeta("X-Container-Meta-Century", "saka02")
	header, err := c.CreateContainerMeta("bb94ba2", metadata)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log("X-Container-Meta-Author", header["X-Container-Meta-Author"])
	t.Log("X-Container-Meta-Century", header["X-Container-Meta-Century"])
	header, err = c.ShowContainerMeta("bb94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if header["X-Container-Meta-Author"][0] != "saka01" || header["X-Container-Meta-Century"][0] != "saka02" {
		t.Errorf("Expected error: %s", "Could not set Metadata.")
	}
	t.Log("X-Container-Meta-Author", header["X-Container-Meta-Author"])
	t.Log("X-Container-Meta-Century", header["X-Container-Meta-Century"])
}

func TestUpdateContainerMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	metadata := NewMetadata()
	metadata.SetMeta("X-Container-Meta-Author", "saka10")
	metadata.SetMeta("X-Container-Meta-Century", "saka20")
	header, err := c.UpdateContainerMeta("bb94ba2", metadata)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err = c.ShowContainerMeta("bb94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if header["X-Container-Meta-Author"][0] != "saka10" || header["X-Container-Meta-Century"][0] != "saka20" {
		t.Errorf("Expected error: %s", "Could not set Metadata.")
	}
	t.Log("X-Container-Meta-Author", header["X-Container-Meta-Author"])
	t.Log("X-Container-Meta-Century", header["X-Container-Meta-Century"])
}

func TestDeleteContainerMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	metadata := NewMetadata()
	metadata.SetDeleteMeta("X-Container-Meta-Author")
	metadata.SetDeleteMeta("X-Container-Meta-Century")
	header, err := c.DeleteContainerMeta("bb94ba2", metadata)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err = c.ShowContainerMeta("bb94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if _, ok := header["X-Container-Meta-Author"]; ok != false {
		t.Errorf("Expected error: %s", "Could not delete Metadata.")
	}
	if _, ok := header["X-Container-Meta-Century"]; ok != false {
		t.Errorf("Expected error: %s", "Could not delete Metadata.")
	}
	t.Log("X-Container-Meta-Author", header["X-Container-Meta-Author"])
	t.Log("X-Container-Meta-Century", header["X-Container-Meta-Century"])
}

// // Object operation
func TestCreateObjectWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	var jsonStr = []byte(`[{"count": 0, "bytes": 0, "name": "aaa"},
{"count": 0, "bytes": 0, "name": "ddddddd"},
{"count": 0, "bytes": 0, "name": "fewra"},
{"count": 19, "bytes": 8464267, "name": "saka-test"},
{"count": 1, "bytes": 66666, "name": "ssbench_000000"}]
`)
	err := ioutil.WriteFile("dd94ba2.json", jsonStr, 0644)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err := c.CreateObject("bb94ba2", "dd94ba2.json", "dd94ba2.json", nil)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log(header)
	os.Remove("dd94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
}

func TestCreateObjectWithAuthUrlAndAddContentType(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	var jsonStr = []byte(`[{"count": 0, "bytes": 0, "name": "aaa"},
{"count": 0, "bytes": 0, "name": "ddddddd"},
{"count": 0, "bytes": 0, "name": "fewra"},
{"count": 19, "bytes": 8464267, "name": "saka-test"},
{"count": 1, "bytes": 66666, "name": "ssbench_000000"}]
`)
	err := ioutil.WriteFile("aa94ba2.json", jsonStr, 0644)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	metadata := NewMetadata()
	metadata.SetMeta("Content-Type", "application/json")
	header, err := c.CreateObject("bb94ba2", "aa94ba2.json", "aa94ba2.json", metadata)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log(header)
}

func TestListObjectsWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	res, header, err := c.ListObjects("bb94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	for i := range res {
		t.Log(res[i].Name)
	}
	t.Log(header)
}

func TestListObjectsWithParamsAndAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	res, header, err := c.ListObjectsWithParams("bb94ba2", Params{Limit: 1})
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log(res)
	for i := range res {
		t.Log(res[i].Name)
	}
	t.Log(header)
}

func TestGetObjectWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	res, err := c.GetObject("bb94ba2", "aa94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	err = ioutil.WriteFile("aa94ba2.json", res, 0644)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	err = os.Remove("aa94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
}

// func TestCreateObjectChunkedWithAuthUrl(t *testing.T) {
// 	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
// 	err := c.CreateObject("saka-test", "saka-test.png", "madoka.png", Params{ChunkSize: 1024 * 10})
// 	if err != nil {
// 		t.Errorf("Expected error: %s", err)
// 	}
// }

func TestDeleteObjectWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	err := c.DeleteObject("bb94ba2", "aa94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
}

func TestCopyObjectWithAuthUrl(t *testing.T) {
	// create container
	c := Client{StorageUrl: StorageUrl, Token: Token, SkipSecure: true}
	header, err := c.CreateContainer("ee94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}

	// copy object
	header, err = c.CopyObject("bb94ba2", "dd94ba2.json", "ee94ba2", "ee94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log(header["X-Copied-From"][0])
	header, err = c.ShowObjectMeta("ee94ba2", "ee94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
}

func TestCreateObjectMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	metadata := NewMetadata()
	metadata.SetMeta("X-Delete-After", "30")
	header, err := c.CreateObjectMeta("bb94ba2", "dd94ba2.json", metadata)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err = c.ShowObjectMeta("bb94ba2", "dd94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if _, ok := header["X-Delete-At"]; ok != true {
		t.Errorf("Expected error: %s", "Could not create Metadata.")
	}
	t.Log("X-Delete-At", header["X-Delete-At"])
}

func TestUpdateObjectMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}

	metadata := NewMetadata()
	metadata.SetMeta("X-Delete-After", "360")
	header, err := c.UpdateObjectMeta("bb94ba2", "dd94ba2.json", metadata)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err = c.ShowObjectMeta("bb94ba2", "dd94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if _, ok := header["X-Delete-At"]; ok != true {
		t.Errorf("Expected error: %s", "Could not create Metadata.")
	}
	t.Log("X-Delete-At", header["X-Delete-At"])
}

func TestShowObjectMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	header, err := c.ShowObjectMeta("bb94ba2", "dd94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	t.Log("X-Delete-At", header["X-Delete-At"])
}

func TestDeleteObjectMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	metadata := NewMetadata()
	metadata.SetDeleteMeta("X-Delete-At")
	header, err := c.DeleteObjectMeta("bb94ba2", "dd94ba2.json", metadata)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err = c.ShowObjectMeta("bb94ba2", "dd94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if _, ok := header["X-Delete-At"]; ok == true {
		t.Errorf("Expected error: %s", "Could not create Metadata.")
	}
}

func TestCleanAll(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, SkipSecure: true}
	err := c.DeleteObject("bb94ba2", "dd94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	_, err = c.DeleteContainer("bb94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	err = c.DeleteObject("ee94ba2", "ee94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	_, err = c.DeleteContainer("ee94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
}
