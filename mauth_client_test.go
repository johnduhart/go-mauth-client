package go_mauth_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestFullURLWithRelative(t *testing.T) {
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.CreateClient("https://innovate.mdsol.com")
	expected := "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ := client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen")

	}
	// now, with a trailing slash
	client, _ = mauth_app.CreateClient("https://innovate.mdsol.com/")
	expected = "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ = client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen: ", actual)

	}
}

func TestFullURLWithRelativeAndParams(t *testing.T) {
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.CreateClient("https://innovate.mdsol.com")
	expected := "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ := client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen")

	}
	// now, with a trailing slash
	client, _ = mauth_app.CreateClient("https://innovate.mdsol.com/")
	expected = "https://innovate.mdsol.com/api/v2/users.json"
	actual, _ = client.fullURL("/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen: ", actual)

	}
}

func TestFullURLWithActualURL(t *testing.T) {
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.CreateClient("https://innovate.mdsol.com")
	expected := "https://balance-innovate.mdsol.com/api/v2/users.json"
	actual, _ := client.fullURL("https://balance-innovate.mdsol.com/api/v2/users.json")
	if actual != expected {
		t.Error("Expected URL not seen")

	}
}

func TestCreateClient(t *testing.T) {
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.CreateClient("https://innovate.mdsol.com")
	if client.baseUrl != "https://innovate.mdsol.com" {
		t.Error("Base URL has changed")
	}
	if client.mauthApp.AppId != app_id {
		t.Error("App ID has changed")
	}
}

func TestCreateClientBadURL(t *testing.T) {
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	_, err := mauth_app.CreateClient("some_nonsense")
	if err == nil {
		t.Error("Bad URL should fail")
	}
}

func hasMWSHeader(r *http.Request) bool {
	for header := range r.Header {
		if header == "X-Mws-Authentication" {
			return true
		}
	}
	return false
}

// Test the Get call
func TestMAuthClient_Get(t *testing.T) {
	var verb string
	var req_url string
	has_mws_header := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req_url = r.URL.String()
		verb = r.Method
		has_mws_header = hasMWSHeader(r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"fake twitter json string"}`)
	}))
	defer server.Close()
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.CreateClient(server.URL)
	// Make the Get call
	_, err := client.Get("/api/v2/users.json")
	if err != nil {
		t.Error("Get Failed: ", err)
	}
	if verb != "GET" {
		t.Error("Expected GET, got ", verb)
	}
	if !has_mws_header {
		t.Error("Expected header not present")
	}
}

// Test the Delete call
func TestMAuthClient_Delete(t *testing.T) {
	var verb string
	var url string
	has_mws_header := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url = r.URL.String()
		verb = r.Method

		has_mws_header = hasMWSHeader(r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"fake twitter json string"}`)
	}))
	defer server.Close()
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.CreateClient(server.URL)
	// Make the Get call
	_, err := client.Delete("/api/v2/users.json")
	if err != nil {
		t.Error("Delete Failed: ", err)
	}
	if verb != "DELETE" {
		t.Error("Expected DELETE, got ", verb)
	}
	if !has_mws_header {
		t.Error("Expected header not present")
	}
}

// Test the Post call
func TestMAuthClient_Post(t *testing.T) {
	var verb string
	var url string
	has_mws_header := false
	var content_type string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url = r.URL.String()
		verb = r.Method

		has_mws_header = hasMWSHeader(r)
		for header, value := range r.Header {
			if header == "Content-Type" {
				content_type = strings.Join(value, "")
			}
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"fake twitter json string\"}")
	}))
	defer server.Close()
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.CreateClient(server.URL)
	// Make the Get call
	response, err := client.Post("/api/v2/users.json", `{"uuid":"1234-1234"}`)
	if err != nil {
		t.Error("Post Failed: ", err)
	}
	if verb != "POST" {
		t.Error("Expected POST, got ", verb)
	}
	if !has_mws_header {
		t.Error("Expected header not present")
	}
	if content_type != "application/json" {
		t.Error("Expected Content-type not set")
	}
	content, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if strings.Compare(string(content[:]),
		"{\"fake twitter json string\"}") != 0 {
		t.Error("Unexpected response body: ", string(content[:]))
	}
}

// Test the Put call
func TestMAuthClient_Put(t *testing.T) {
	var verb string
	var url string
	has_mws_header := false
	var content_type string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url = r.URL.String()
		verb = r.Method

		has_mws_header = hasMWSHeader(r)
		for header, value := range r.Header {
			if header == "Content-Type" {
				content_type = strings.Join(value, "")
			}
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"fake twitter json string\"}")
	}))
	defer server.Close()
	mauth_app, _ := LoadMauth(app_id, filepath.Join("test", "private_key.pem"))
	client, _ := mauth_app.CreateClient(server.URL)
	// Make the Get call
	response, err := client.Put("/api/v2/users.json", `{"uuid":"1234-1234"}`)
	if err != nil {
		t.Error("Post Failed: ", err)
	}
	if verb != "PUT" {
		t.Error("Expected PUT, got ", verb)
	}
	if !has_mws_header {
		t.Error("Expected header not present")
	}
	if content_type != "application/json" {
		t.Error("Expected Content-type not set")
	}
	content, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if strings.Compare(string(content[:]),
		"{\"fake twitter json string\"}") != 0 {
		t.Error("Unexpected response body: ", string(content[:]))
	}
}

// Example of creating a MAuth Client
func ExampleMAuthApp_CreateClient() {
	// given an APP_UUID
	var appUUID = "7D0B2A90-0825-4AD8-9C1F-E9851795D428"
	// and a path to a KeyFile
	var keyPath = filepath.Join("test", "private_key.pem")
	// create a MAuth client
	var client *MAuthApp
	client, err := LoadMauth(appUUID, keyPath)
	if err != nil {
		log.Fatal("Unable to create client: ", err)
	}
	// Define a base URL
	var baseURL = "https://innovate.imedidata.com"
	var mauthClient *MAuthClient
	mauthClient, err = client.CreateClient(baseURL)
	if err != nil {
		log.Fatal("Unable to create MAuth Client: ", err)
	}
	println("Successfully created MAuth Client for APP: ", mauthClient.mauthApp.AppId)
}

// Example of creating a MAuth Client and making a Get Request
func ExampleMAuthClient_Get() {
	// Get information on a User
	// http://developer.imedidata.com/desktop/ActionTopics/Users/Listing_User_Account_Details.htm

	// given an APP_UUID
	var appUUID = "7D0B2A90-0825-4AD8-9C1F-E9851795D428"
	// and a path to a KeyFile
	var keyPath = filepath.Join("test", "private_key.pem")
	// create a MAuth client
	var client *MAuthApp
	client, err := LoadMauth(appUUID, keyPath)
	if err != nil {
		log.Fatal("Unable to create client: ", err)
	}
	// Define a base URL
	var baseURL = "https://innovate.imedidata.com"

	// Define and create the Client
	var mauthClient *MAuthClient
	mauthClient, err = client.CreateClient(baseURL)
	if err != nil {
		log.Fatal("Unable to create MAuth Client: ", err)
	}
	// This is made-up
	var userUuid = "347942BF-9915-405D-BB20-6196597F3BE3"
	response, err := mauthClient.Get("api/v2/users/" + userUuid + ".json")
	println("Got a status code of", response.StatusCode, "for request for User UUID", userUuid)
}

func ExampleMAuthClient_Post() {
	// Creating a Study Using a MAuth Client
	// http://developer.imedidata.com/desktop/ActionTopics/Studies/Creating_Studies.htm

	// given an APP_UUID
	var appUUID = "7D0B2A90-0825-4AD8-9C1F-E9851795D428"
	// and a path to a KeyFile
	var keyPath = filepath.Join("test", "private_key.pem")
	// create a MAuth client
	var client *MAuthApp
	client, err := LoadMauth(appUUID, keyPath)
	if err != nil {
		log.Fatal("Unable to create client: ", err)
	}
	// Define a base URL
	var baseURL = "https://innovate.imedidata.com"

	// Define and create the Client
	var mauthClient *MAuthClient
	mauthClient, err = client.CreateClient(baseURL)
	if err != nil {
		log.Fatal("Unable to create MAuth Client: ", err)
	}

	// Define the constituent entity references
	var studyGroupUUID = "347942BF-9915-405D-BB20-6196597F3BE3"
	var studyUUID = "C3C79E4A-4BFD-4A72-89E9-724A4E6A9D95"

	// This is a slimmed down version of the structure from the reference above
	type studyDefinition struct {
		Number           int    `json:"number"`
		Name             string `json:"name"`
		IsProduction     bool   `json:"is_production"`
		TherapeticArea   string `json:"therapeutic_area"`
		FullDescription  string `json:"full_description"`
		CompoundCode     string `json:"compound_code"`
		DrugDevice       string `json:"drug_device"`
		Title            string `json:"title"`
		UUID             string
		Protocol         string `json:"protocol"`
		ParentUUID       string `json:"parent_UUID"`
		EnrollmentTarget int    `json:"enrollment_target"`
		OID              string `json:"oid"`
	}

	// Create an instance of the new study
	study := &studyDefinition{
		Number:           1,
		Name:             "ABC1234",
		IsProduction:     true,
		TherapeticArea:   "Endocrine",
		FullDescription:  "Some Sample Study",
		CompoundCode:     "Mediflex",
		DrugDevice:       "Drug",
		Title:            "A sample Endocrine Study",
		UUID:             studyUUID,
		Protocol:         "ABC1234",
		ParentUUID:       "",
		EnrollmentTarget: 150,
		OID:              "ABC1234",
	}
	data, _ := json.Marshal(study)

	// POST www.imedidata.com/api/v2/study_groups/[study group uuid]/studies.json
	response, err := mauthClient.Post("api/v2/study_groups/"+studyGroupUUID+"/studies.json",
		string(data))
	println("Got a status code of", response.StatusCode, "for request to create Study", studyUUID)
}
