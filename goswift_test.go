package goswift

import (
	"io/ioutil"
	"os"
	"testing"
)

var AuthUrl = os.Getenv("SWIFT_AUTH_URL")
var AccountName = os.Getenv("SWIFT_API_USER")
var Password = os.Getenv("SWIFT_API_KEY")
var TenantName = os.Getenv("SWIFT_TENANT_NAME")
var RegionName = os.Getenv("SWIFT_REGION_NAME")

// Credentials
func TestCredentials(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	err := c.Credential()
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if c.StorageUrl != "" {
		t.Log(c.StorageUrl)
	} else {
		t.Errorf("Expected error: %s", err)
	}
	if c.Token != "" {
		t.Log(c.Token)
	} else {
		t.Errorf("Expected error: %s", err)
	}
}

// Account metadata operation
func TestShowAccountMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	header, err := c.ShowAccountMeta()
	if err != nil {
		t.Errorf("Expected error: %s", err)
	} else {
		t.Log(header)
	}
	if v, ok := header["X-Account-Bytes-Used"]; ok == true {
		t.Log("X-Account-Bytes-Used", v[0])
	}
	if v, ok := header["X-Account-Container-Count"]; ok == true {
		t.Log("X-Account-Container-Count", v[0])
	}
}

func TestCreateAccountMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	if v, ok := header["X-Account-Meta-Book"]; ok == true {
		if v[0] != "saka01" {
			t.Errorf("Expected error: %s", "Could not set Metadata.")
		} else {
			t.Log("X-Account-Meta-Book", v)
		}
	}
	if v, ok := header["X-Account-Meta-Subject"]; ok == true {
		if v[0] != "saka02" {
			t.Errorf("Expected error: %s", "Could not create Metadata.")
		} else {
			t.Log("X-Account-Meta-Subject", v)
		}
	}
}

func TestUpdateAccountMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	if v, ok := header["X-Account-Meta-Book"]; ok == true {
		if v[0] != "saka10" {
			t.Errorf("Expected error: %s", "Could not set Metadata.")
		} else {
			t.Log("X-Account-Meta-Book", v)
		}
	}
	if v, ok := header["X-Account-Meta-Subject"]; ok == true {
		if v[0] != "saka20" {
			t.Errorf("Expected error: %s", "Could not create Metadata.")
		} else {
			t.Log("X-Account-Meta-Subject", v)
		}
	}

}

func TestDeleteAccountMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	} else {
		t.Log(header)
	}
	if v, ok := header["X-Account-Meta-Book"]; ok != false {
		t.Errorf("Expected error: %s", "Could not delete Metadata.")
		t.Log("X-Account-Meta-Book", v)
	}
	if v, ok := header["X-Account-Meta-Subject"]; ok != false {
		t.Errorf("Expected error: %s", "Could not delete Metadata.")
		t.Log("X-Account-Meta-Subject", v)
	}
}

// Container operation
func TestListContainersWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	res, _, err := c.ListContainers()
	if err != nil {
		t.Errorf("Expected error: %s", err)
	} else {
		t.Log(len(res))
	}
}

func TestListContainersWithParamsAndLimit4(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	res, _, err := c.ListContainersWithParams(Params{Limit: 1})
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected error: Limit count is different.")
	} else {
		t.Log(res)
	}
}

func TestCreateContainersWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	header, err := c.CreateContainer("bb94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	header, err = c.CreateContainer("9194ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	} else {
		t.Log(header)
	}
}

func TestDeleteContainersWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	header, err := c.DeleteContainer("9194ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	} else {
		t.Log(header)
	}
}

// Container metadata operation
func TestShowContainerMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	header, err := c.ShowContainerMeta("bb94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if v, ok := header["X-Container-Meta-Author"]; ok == true {
		t.Log("X-Container-Meta-Author", v[0])
	}
	if v, ok := header["X-Container-Meta-Century"]; ok == true {
		t.Log("X-Container-Meta-Century", v[0])
	}
}

func TestCreateContainerMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	if v, ok := header["X-Container-Meta-Author"]; ok == true {
		if v[0] != "saka01" {
			t.Errorf("Expected error: %s", "Could not set Metadata.")
		} else {
			t.Log("X-Container-Meta-Author", v[0])
		}
	}
	if v, ok := header["X-Container-Meta-Century"]; ok == true {
		if v[0] != "saka02" {
			t.Errorf("Expected error: %s", "Could not create Metadata.")
		} else {
			t.Log("X-Container-Meta-Century", v[0])
		}
	}
}

func TestUpdateContainerMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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

	if v, ok := header["X-Container-Meta-Author"]; ok == true {
		if v[0] != "saka10" {
			t.Errorf("Expected error: %s", "Could not set Metadata.")
		} else {
			t.Log("X-Container-Meta-Author", v[0])
		}
	}
	if v, ok := header["X-Container-Meta-Century"]; ok == true {
		if v[0] != "saka20" {
			t.Errorf("Expected error: %s", "Could not create Metadata.")
		} else {
			t.Log("X-Container-Meta-Century", v[0])
		}
	}
}

func TestDeleteContainerMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	if v, ok := header["X-Container-Meta-Author"]; ok != false {
		t.Errorf("Expected error: %s", "Could not delete Metadata.")
		t.Log("X-Container-Meta-Author", v)
	}
	if v, ok := header["X-Container-Meta-Century"]; ok != false {
		t.Errorf("Expected error: %s", "Could not delete Metadata.")
		t.Log("X-Container-Meta-Century", v)
	}
}

// // Object operation
func TestCreateObjectWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	} else {
		t.Log(header)
	}
	os.Remove("dd94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
}

func TestCreateObjectWithAuthUrlAndAddContentType(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	} else {
		t.Log(header)
	}
}

func TestListObjectsWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	res, header, err := c.ListObjects("bb94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	} else {
		t.Log(header)
	}
	for i := range res {
		t.Log(res[i].Name)
	}
}

func TestListObjectsWithParamsAndAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	res, header, err := c.ListObjectsWithParams("bb94ba2", Params{Limit: 1})
	if err != nil {
		t.Errorf("Expected error: %s", err)
	} else {
		t.Log(res)
		t.Log(header)
	}
	for i := range res {
		t.Log(res[i].Name)
	}
}

func TestGetObjectWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	err := c.DeleteObject("bb94ba2", "aa94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
}

func TestCopyObjectWithAuthUrl(t *testing.T) {
	// create container
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	header, err := c.CreateContainer("ee94ba2")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}

	// copy object
	header, err = c.CopyObject("bb94ba2", "dd94ba2.json", "ee94ba2", "ee94ba2.json")
	t.Log(header)
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
	if _, ok := header["X-Copied-From"]; ok == true {
		if header["X-Copied-From"][0] != "bb94ba2/dd94ba2.json" {
			t.Errorf("Expected error: %s", "Could not copy object.")
		}
	}
	header, err = c.ShowObjectMeta("ee94ba2", "ee94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	}
}

func TestCreateObjectMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	if v, ok := header["X-Delete-At"]; ok != true {
		t.Errorf("Expected error: %s", "Could not create Metadata.")
	} else {
		t.Log("X-Delete-At", v[0])
	}
}

func TestUpdateObjectMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	if v, ok := header["X-Delete-At"]; ok != true {
		t.Errorf("Expected error: %s", "Could not create Metadata.")
	} else {
		t.Log("X-Delete-At", v[0])
	}
}

func TestShowObjectMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
	header, err := c.ShowObjectMeta("bb94ba2", "dd94ba2.json")
	if err != nil {
		t.Errorf("Expected error: %s", err)
	} else {
		t.Log("X-Delete-At", header["X-Delete-At"])
	}
}

func TestDeleteObjectMetaWithAuthUrl(t *testing.T) {
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
	c := Client{AuthUrl: AuthUrl, AccountName: AccountName, Password: Password, TenantName: TenantName, RegionName: RegionName, SkipSecure: true}
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
