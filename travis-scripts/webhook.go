package traviswebhook

/*

Copyright 2017 Shapath Neupane (@theshapguy)

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

------------------------------------------------

Listner - written in Go because it's native web server is much more robust than Python. Plus its fun to write Go!

NOTE: Make sure you are using the right domain for travis [.com] or [.org]

*/

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/container/v1"
)

var logPrint = log.Println

type configKey struct {
	Config struct {
		Host        string `json:"host"`
		ShortenHost string `json:"shorten_host"`
		Assets      struct {
			Host string `json:"host"`
		} `json:"assets"`
		Pusher struct {
			Key string `json:"key"`
		} `json:"pusher"`
		Github struct {
			APIURL string   `json:"api_url"`
			Scopes []string `json:"scopes"`
		} `json:"github"`
		Notifications struct {
			Webhook struct {
				PublicKey string `json:"public_key"`
			} `json:"webhook"`
		} `json:"notifications"`
	} `json:"config"`
}

func payloadSignature(r *http.Request) ([]byte, error) {

	signature := r.Header.Get("Signature")
	b64, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return nil, errors.New("cannot decode signature")
	}

	return b64, nil
}

func parsePublicKey(key string) (*rsa.PublicKey, error) {

	// https://golang.org/pkg/encoding/pem/#Block
	block, _ := pem.Decode([]byte(key))

	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("invalid public key")
	}

	return publicKey.(*rsa.PublicKey), nil

}

func travisPublicKey() (*rsa.PublicKey, error) {
	// NOTE: Use """https://api.travis-ci.com/config""" for private repos.
	response, err := http.Get("https://api.travis-ci.org/config")

	if err != nil {
		return nil, errors.New("cannot fetch travis public key")
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var t configKey
	err = decoder.Decode(&t)
	if err != nil {
		return nil, errors.New("cannot decode travis public key")
	}

	key, err := parsePublicKey(t.Config.Notifications.Webhook.PublicKey)
	if err != nil {
		return nil, err
	}

	return key, nil

}

func payloadDigest(payload string) []byte {
	hash := sha1.New()
	hash.Write([]byte(payload))
	return hash.Sum(nil)
}

func respondWithError(w http.ResponseWriter, m string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)
	message := fmt.Sprintf("{\"message\": \"%s\"}", m)
	w.Write([]byte(message))
}

func respondWithSuccess(w http.ResponseWriter, m string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	message := fmt.Sprintf("{\"message\": \"%s\"}", m)
	w.Write([]byte(message))
}

func DeployHandler(w http.ResponseWriter, r *http.Request) {

	key, err := travisPublicKey()
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	signature, err := payloadSignature(r)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	payload := payloadDigest(r.FormValue("payload"))

	err = rsa.VerifyPKCS1v15(key, crypto.SHA1, payload, signature)

	if err != nil {
		respondWithError(w, errors.New("unauthorized payload").Error())
		return
	}

	resp, err := deleteCluster()

	if err != nil {
		log.Printf("Cluster delete error=%v", err)
		respondWithError(w, err.Error())
		return
	}
	if resp == nil {
		log.Printf("No cluster to delete")
		respondWithSuccess(w, "No cluster to deleted")
		return
	}

	log.Printf("Cluster deleted. Response status=" + resp.Status)
	respondWithSuccess(w, "cluster deleted. Response status="+resp.Status)
}

func deleteCluster() (*container.Operation, error) {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, container.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	containerService, err := container.New(c)
	if err != nil {
		return nil, err
	}

	// Deprecated. The Google Developers Console [project ID or project
	// number](https://support.google.com/cloud/answer/6158840).
	// This field has been deprecated and replaced by the name field.
	projectID := os.Getenv("PROJECT_ID")

	// Deprecated. The name of the Google Compute Engine
	// [zone](/compute/docs/zones#available) in which the cluster
	// resides.
	// This field has been deprecated and replaced by the name field.
	zone := os.Getenv("ZONE")

	// Deprecated. The name of the cluster to delete.
	// This field has been deprecated and replaced by the name field.
	clusterID := os.Getenv("CLUSTER_ID")
	resp, err := containerService.Projects.Zones.Clusters.List(projectID, zone).Context(ctx).Do()

	if err != nil {
		return nil, err
	}
	for _, c := range resp.Clusters {
		if c.Name == clusterID {
			return containerService.Projects.Zones.Clusters.Delete(projectID, zone, clusterID).Context(ctx).Do()
		}
	}
	return nil, nil
}
