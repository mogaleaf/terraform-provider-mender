package api

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func launchServer(t *testing.T) {
	http.HandleFunc("/api/management/v1/useradm/auth/login", func(writer http.ResponseWriter, request *http.Request) {
		assert.EqualValues(t, "Basic dGVzdDpjZWNpbGU=", request.Header.Get("Authorization"))
		writer.Write([]byte("token"))
	})

	http.HandleFunc("/api/management/v1/deployments/artifacts", func(writer http.ResponseWriter, request *http.Request) {
		assert.EqualValues(t, "Bearer token", request.Header.Get("Authorization"))
		writer.WriteHeader(201)
		writer.Write([]byte("abc/def"))
	})
	http.ListenAndServe("localhost:8090", nil)
}

func TestClient_Login(t *testing.T) {
	go launchServer(t)
	c := New("http://localhost:8090", "test", "cecile")
	err := c.Login()
	assert.Nil(t, err)
	assert.True(t, c.(*client).isAuthenticated)
	assert.EqualValues(t, c.(*client).token, "token")
	println(c.(*client).token)
}

func TestClient_UploadArtifact(t *testing.T) {
	go launchServer(t)
	c := New("http://localhost:8090", "test", "cecile")
	err := c.Login()
	assert.Nil(t, err)
	assert.True(t, c.(*client).isAuthenticated)
	dat, err := os.ReadFile("test.mender")
	assert.Nil(t, err)
	fmt.Println(c.UploadArtifact(dat))
}
