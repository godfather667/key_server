// key_server_test - Test key_server restAPI address book database
//
// This program DOES NOT demonstrate the most compact type of Testing.
// It purposely demonstrates several testing techniques and the code
// is straight forward without numerous table driven opaque functions.
//
// Further this test must be run in sequence. Omitting Tests will invalidate
// the results as the tests depend on each other. This is not the best technique
// but it produces a simple test of all operations.
//
// Productions tests should be independent of each other. In practice, this
// means there will be a lot of duplication when working with Database Functions.
//
//For a table driven testing see "db_demo" at https://github.com/godfather667
package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	_ "strconv"
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
//  Test "TestDeletePerson" Function
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

//
//  Test "TestModifyPerson" Function
//
func TestModifyPerson(t *testing.T) {
	person3 := &Person{
		UniqID:    "3",
		FirstName: "Nancy",
		LastName:  "Chow",
		EmailAddr: "x@x.com",
		PhoneNumb: "555-555-0000",
	}
	jsonPerson, _ := json.Marshal(person3)
	request, _ := http.NewRequest("PUT", "http://localhost:8000/address/3", bytes.NewBuffer(jsonPerson))
	response := executeRequest(request)

	checkResponseCode(t, http.StatusOK, response.Code)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	if len(response.Body.Bytes()) > 0 {
		t.Error("Body returned unknown data, should be empty: ", response.Body.Bytes())
	}
	expected_md5 := []byte{54, 62, 217, 186, 210, 30, 205, 213, 3, 249, 151, 217, 81, 17, 147, 213}
	md5 := getMD5("Data.db")
	for i, v := range expected_md5 {
		if v != md5[i] {
			t.Errorf("Database MD5 not equal to expected value!")
			break
		}
	}
}

//
// Test Get Person
//
func TestGetPerson(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost:8000/address/3", nil)
	response := executeRequest(request)

	checkResponseCode(t, http.StatusOK, response.Code)
	assert.Equal(t, 200, response.Code, "OK response is expected")

	expectedBytes := []byte{34, 120, 64, 120, 46, 99, 111, 109, 34, 10}
	for i, v := range response.Body.Bytes() {
		if expectedBytes[i] != v {
			t.Errorf("Expected: %s  Result = %s", string(expectedBytes), response.Body.Bytes())
			break
		}
	}
}
func TestGetBook(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://localhost:8000/address", nil)
	response := executeRequest(request)

	checkResponseCode(t, http.StatusOK, response.Code)
	assert.Equal(t, 200, response.Code, "OK response is expected")

	expectedBytes := []byte{123, 34, 48, 34, 58, 123, 34, 117, 110, 105, 113, 95, 105, 100, 34, 58,
		34, 51, 34, 44, 34, 102, 105, 114, 115, 116, 95, 110, 97, 109, 101, 34, 58, 34, 45, 102,
		105, 114, 115, 116, 45, 34, 44, 34, 108, 97, 115, 116, 95, 110, 97, 109, 101, 34, 58, 34,
		45, 108, 97, 115, 116, 45, 34, 44, 34, 101, 109, 97, 105, 108, 95, 97, 100, 100, 114, 34,
		58, 34, 45, 101, 109, 97, 105, 108, 45, 34, 44, 34, 112, 104, 111, 110, 101, 95, 110, 117,
		109, 98, 34, 58, 34, 45, 112, 104, 111, 110, 101, 45, 34, 125, 44, 34, 49, 34, 58, 123, 34,
		117, 110, 105, 113, 95, 105, 100, 34, 58, 34, 49, 34, 44, 34, 102, 105, 114, 115, 116, 95,
		110, 97, 109, 101, 34, 58, 34, 67, 104, 97, 114, 108, 101, 115, 34, 44, 34, 108, 97, 115,
		116, 95, 110, 97, 109, 101, 34, 58, 34, 83, 109, 105, 116, 104, 34, 44, 34, 101, 109, 97,
		105, 108, 95, 97, 100, 100, 114, 34, 58, 34, 120, 64, 120, 46, 99, 111, 109, 34, 44, 34,
		112, 104, 111, 110, 101, 95, 110, 117, 109, 98, 34, 58, 34, 53, 53, 53, 45, 53, 53, 53,
		45, 48, 48, 48, 48, 34, 125, 44, 34, 51, 34, 58, 123, 34, 117, 110, 105, 113, 95, 105,
		100, 34, 58, 34, 51, 34, 44, 34, 102, 105, 114, 115, 116, 95, 110, 97, 109, 101, 34, 58,
		34, 78, 97, 110, 99, 121, 34, 44, 34, 108, 97, 115, 116, 95, 110, 97, 109, 101, 34, 58,
		34, 67, 104, 111, 119, 34, 44, 34, 101, 109, 97, 105, 108, 95, 97, 100, 100, 114, 34, 58,
		34, 120, 64, 120, 46, 99, 111, 109, 34, 44, 34, 112, 104, 111, 110, 101, 95, 110, 117,
		109, 98, 34, 58, 34, 53, 53, 53, 45, 53, 53, 53, 45, 48, 48, 48, 48, 34, 125, 125, 10}
	for i, v := range response.Body.Bytes() {
		if expectedBytes[i] != v {
			t.Errorf("Expected: %s\n  Result = %s", string(expectedBytes), response.Body.Bytes())
			break
		}
	}
}

//
// Test Export Function
//
func TestExportCSV(t *testing.T) {
	request, _ := http.NewRequest("POST", "/address", nil)
	response := httptest.NewRecorder()
	ExportCSV(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
	var bn []byte // Empty Response!
	if !bytes.Equal(bn, response.Body.Bytes()) {
		t.Errorf("Body Didn't match:\n\tExpected:\t%q\n\tGot:\t%q", bn, response.Body.String())
	}
	dat, err := ioutil.ReadFile("Data.csv")
	check("Read of Data.csv Failed! ", err)

	expectedBytes := []byte{49, 44, 67, 104, 97, 114, 108, 101, 115, 44, 83, 109, 105, 116, 104, 44, 120, 64, 120, 46, 99, 111, 109, 44, 53, 53, 53,
		45, 53, 53, 53, 45, 48, 48, 48, 48, 10, 51, 44, 78, 97, 110, 99, 121, 44, 67, 104, 111, 119, 44, 120, 64, 120, 46, 99, 111, 109, 44, 53, 53, 53, 45,
		53, 53, 53, 45, 48, 48, 48, 48, 10}
	for i, v := range dat {
		if expectedBytes[i] != v {
			t.Errorf("Expected: %s\n  Result = %s", string(expectedBytes), dat)
			break
		}
	}
}

//
// Test Import Function
//
func TestImportCSV(t *testing.T) {
	OrgStore := make(map[int]Person)
	for i, _ := range KeyStore { // Make Copy of Orginal KeyStore
		OrgStore[i] = KeyStore[i]
	}
	request, _ := http.NewRequest("POST", "/address", nil)
	response := httptest.NewRecorder()
	ImportCSV(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")

	if KeyStore[1].FirstName != "Charles" ||
		KeyStore[1].LastName != "Smith" ||
		KeyStore[1].EmailAddr != "x@x.com" ||
		KeyStore[1].PhoneNumb != "555-555-0000" {
		t.Error("Expected: Charles Smith -- Got: ", KeyStore[1])
	}
	if KeyStore[2].FirstName != "Nancy" ||
		KeyStore[2].LastName != "Chow" ||
		KeyStore[2].EmailAddr != "x@x.com" ||
		KeyStore[2].PhoneNumb != "555-555-0000" {
		t.Error("Expected: Nancy Chow -- Got: ", KeyStore[2])
	}
}
