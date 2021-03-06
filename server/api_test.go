package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/tinyurl/tinyurl/entity"
	"github.com/tinyurl/tinyurl/store"
)

const (
	TestPort          = "8877"
	TestAddr          = "http://0.0.0.0:8877"
	ConfigPathDefault = "../default.properties"
	TestOriginURL     = "http://test.origin.url"
	TestShortPath     = "shortpath"
)

var (
	ConfigPath  string
	storeClient entity.URLStore
	appService  *entity.ServiceProvider
)

func init() {
	os.Setenv("TINYURL_CONFIG_PATH", "../default.properties")
	ConfigPath = os.Getenv("TINYURL_CONFIG_PATH")
	if ConfigPath == "" {
		ConfigPath = ConfigPathDefault
	}

	storeClient = store.NewGeneralDBClient(ConfigPath)
	globalConfig := entity.GetGlobalConfig(ConfigPath)
	appService = &entity.ServiceProvider{
		StoreClient:  storeClient,
		KeyGenerater: entity.NewKeyGenerater(globalConfig.KeyAlgo),
		GlobalConfig: entity.GetGlobalConfig(ConfigPath),
	}
}

func newTestURL() entity.URL {
	return entity.URL{
		OriginURL: TestOriginURL,
		ShortPath: TestShortPath,
	}
}

func insertTestURL(url entity.URL) {
	appService.StoreClient.Create(&url)
}

func updateTestURL(url entity.URL) {
	appService.StoreClient.Update(&url)
}

func clearDatabase() {
	appService.StoreClient.DropDatabase()
}

func PostForm(postURL string, data url.Values) interface{} {
	resp, err := http.PostForm(postURL, data)
	if err != nil {
		logrus.Fatalf("post form data to %s error: %v", postURL, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("read response body error: %v", err)
	}

	var ret interface{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		logrus.Fatalf("unmarshal response body error: %v", err)
	}

	return ret
}

// startTestServer
func startTestServer(t *testing.T) {
	go func() {
		r := BuildEngine(appService)
		r.Run(":" + TestPort)
	}()

	t.Logf("wait 2s to start testServer...\n")
	time.Sleep(time.Second * 2)
}
