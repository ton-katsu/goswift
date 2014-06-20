Go bindings to the OpenStack Object Storage API
===============================================

This is a go client for the OpenStack Object Storage(Swift) APIv1.

[http://docs.openstack.org/api/openstack-object-storage/1.0/content/](http://docs.openstack.org/api/openstack-object-storage/1.0/content/)

> Now support swauth(AuthV1) only.
> Keystone(AuthV2) support comming soon.

> "Transfer-Encoding: Chunked" support comming soon.

Install
-------

    go get github.com/ton-katsu/goswift


Usage
-----

GoDoc:  [https://godoc.org/github.com/ton-katsu/goswift](https://godoc.org/github.com/ton-katsu/goswift)


Example
-------

#### Import module

    import github.com/ton-katsu/goswift


#### Create client


    c := goswift.Client{AuthUrl: "auth_url", AccountName: "account_name", Password: "account_key"}

    or

    c := goswift.Client{StorageUrl: "storage_url", Token: "account_token"}


#### List containers

    containers, header, err := c.ListContainers()

#### List containers with params

    containers, header, err := c.ListContainersWithParams(Params{Limit: 3, Marker: "tonkatsu"})

#### Create object metadata

    metadata := NewMetadata()
    metadata.SetMeta("X-Delete-After", "30")
    header, err := c.CreateObjectMeta("test", "test.jpg", metadata)
 

#### Delete object metadata

    metadata := NewMetadata()
    metadata.SetDeleteMeta("X-Delete-At")
    header, err := c.DeleteObjectMeta("test", "test.jpg", metadata)


#### Create object

    var jsonStr = []byte(`[{"count": 0, "bytes": 0, "name": "aaa"},
    {"count": 0, "bytes": 0, "name": "ton-katsu"},
    {"count": 0, "bytes": 0, "name": "ebi-katsu"},
    {"count": 19, "bytes": 8464267, "name": "katsu-don"}]
    `)

    ioutil.WriteFile("test.json", jsonStr, 0644)

    metadata := NewMetadata()
    metadata.SetMeta("Content-Type", "application/json")

    header, err := c.CreateObject("test", "test.json", "ton-katsu.json", metadata)


> And more API ... Check GoDoc:  [https://godoc.org/github.com/ton-katsu/goswift](https://godoc.org/github.com/ton-katsu/goswift)

Testing
-------

Set environment variables

    export SWIFT_API_USER='accountname'
    export SWIFT_API_KEY='password'
    export SWIFT_AUTH_URL='https://xxxxx.com/auth/v1.0'
    export SWIFT_USER_TOKEN='xxxxxxx'
    export SWIFT_STORAGE_URL='https://xxxxx.com/v1/xxxxx'

After that, run test.

    go test -v

