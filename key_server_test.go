// key_server_test - Test key_server restapi address book database
package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	//	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

//
// Intialize Routes, KeyStore, and external database.
//
func init() {
	routeInit()
}

//
// Error Function
//
func testCheck(e error) {
	if e != nil {
		panic(e)
	}
}

//
// Execute Request -
//
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

//
// Global Application Values
//

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

//
// Create Empty Initialized Database and KeyStore
//
func initDatabase() {
	for k := range KeyStore { // Clear KeyStore
		delete(KeyStore, k)
	}
	err := os.Remove("Data.db")
	if err != nil {
		fmt.Println("Database Initialized")

	}
	// Create Hidden System Record (Record at uniqID Zero(0) and iinitialize topID
	CreateDatabase()
}

//
// MD5 Generator
//
func getMD5(file string) []byte {
	f, err := os.Open(file)
	check("Open File Error", err)
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		check("Copy File Error", err)
	}
	md5_ret := (h.Sum(nil))
	return md5_ret
}

//
// Test Database Loader
//
func TestLoadDatabase(t *testing.T) {
	const null_db = "null"

	init_db := Person{"0", "-first-", "-last-", "-email-", "-phone-"}

	initDatabase() // Create Initialized Empty Database and KeyStore

	loadDatabase() // Load Database

	if !reflect.DeepEqual(KeyStore[0], init_db) {
		t.Error("\nExpected = ", init_db, "\nReturned = ", KeyStore[0])
	}
}

//
//  Test "TestPost" Function
//
func TestCreatePerson(t *testing.T) {
	person := &Person{
		UniqID:    "1",
		FirstName: "Charles",
		LastName:  "Smith",
		EmailAddr: "x@x.com",
		PhoneNumb: "555-555-0000",
	}
	initDatabase() // Create Initialized Empty Database and KeyStore

	jsonPerson, _ := json.Marshal(person)
	request, _ := http.NewRequest("POST", "/address", bytes.NewBuffer(jsonPerson))
	response := httptest.NewRecorder()
	CreatePerson(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
	var bn []byte // Empty Response!
	if !bytes.Equal(bn, response.Body.Bytes()) {
		t.Errorf("Body Didn't match:\n\tExpected:\t%q\n\tGot:\t%q", bn, response.Body.String())
	}
	person = &Person{
		UniqID:    "2",
		FirstName: "Mike",
		LastName:  "Jones",
		EmailAddr: "x@x.com",
		PhoneNumb: "555-555-0000",
	}
	jsonPerson, _ = json.Marshal(person)
	request, _ = http.NewRequest("POST", "/address", bytes.NewBuffer(jsonPerson))
	response = httptest.NewRecorder()
	CreatePerson(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
	if !bytes.Equal(bn, response.Body.Bytes()) {
		t.Errorf("Body Didn't match:\n\tExpected:\t%q\n\tGot:\t%q", bn, response.Body.String())
	}
	expected_md5 := []byte{17, 178, 136, 214, 219, 159, 205, 5, 139, 214, 206, 234, 33, 135, 168, 111}
	md5 := getMD5("Data.db")
	for i, v := range expected_md5 {
		if v != md5[i] {
			t.Errorf("Database MD5 not equal to expected value!")
		}
	}
}

//
//  Test "TestDelete" Function
//
func TestDeletePerson(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "http://localhost:8000/address/2", nil)
	response := executeRequest(request)

	checkResponseCode(t, http.StatusOK, response.Code)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	fmt.Println("Response.Body", response.Body.Bytes())
	if len(response.Body.Bytes()) > 0 {
		t.Error("Body returned unknown data, should be empty: ", response.Body.Bytes())
	}
	expected_md5 := []byte{218, 214, 173, 125, 156, 5, 86, 165, 46, 91, 27, 36, 187, 240, 80, 52}
	md5 := getMD5("Data.db")
	for i, v := range expected_md5 {
		if v != md5[i] {
			t.Errorf("Database MD5 not equal to expected value!")
			break
		}
	}
}
