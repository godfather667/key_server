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
	person1 := &Person{
		UniqID:    "1",
		FirstName: "Charles",
		LastName:  "Smith",
		EmailAddr: "x@x.com",
		PhoneNumb: "555-555-0000",
	}
	person2 := &Person{
		UniqID:    "2",
		FirstName: "Mike",
		LastName:  "Jones",
		EmailAddr: "x@x.com",
		PhoneNumb: "555-555-0000",
	}
	person3 := &Person{
		UniqID:    "3",
		FirstName: "Mike",
		LastName:  "Jones",
		EmailAddr: "x@x.com",
		PhoneNumb: "555-555-0000",
	}

	initDatabase() // Create Initialized Empty Database and KeyStore

	// Create Three Records
	performPost(person1, t)
	performPost(person2, t)
	performPost(person3, t)

	// Check Results
	expected_md5 := []byte{29, 74, 115, 183, 216, 113, 60, 12, 111, 163, 205, 217, 2, 53, 88, 82}
	md5 := getMD5("Data.db")
	for i, v := range expected_md5 {
		if v != md5[i] {
			t.Errorf("Database MD5 not equal to expected value!")
		}
	}
}

//
// Perform Post Function
//
func performPost(p *Person, t *testing.T) {
	jsonPerson, _ := json.Marshal(p)
	request, _ := http.NewRequest("POST", "/address", bytes.NewBuffer(jsonPerson))
	response := httptest.NewRecorder()
	CreatePerson(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
	var bn []byte // Empty Response!
	if !bytes.Equal(bn, response.Body.Bytes()) {
		t.Errorf("Body Didn't match:\n\tExpected:\t%q\n\tGot:\t%q", bn, response.Body.String())
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

	if len(response.Body.Bytes()) > 0 {
		t.Error("Body returned unknown data, should be empty: ", response.Body.Bytes())
	}
	expected_md5 := []byte{223, 94, 58, 240, 111, 32, 228, 168, 20, 2, 227, 98, 25, 96, 147, 34}
	md5 := getMD5("Data.db")
	for i, v := range expected_md5 {
		if v != md5[i] {
			t.Errorf("Database MD5 not equal to expected value!")
			break
		}
	}
}
