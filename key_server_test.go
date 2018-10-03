// key_server_test - Test key_server restapi address book database
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// Error Function
//
func testCheck(e error) {
	if e != nil {
		panic(e)
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

	saveDatabase() // Create Empty Database
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
func TestPost(t *testing.T) {
	person := &Person{
		UniqID:    "1",
		FirstName: "Charles",
		LastName:  "Smith",
		EmailAddr: "x@x.com",
		PhoneNumb: "555-555-0000",
	}
	initDatabase() // Create Initialized Empty Database and KeyStore

	loadDatabase() // Load Database

	jsonPerson, _ := json.Marshal(person)
	request, _ := http.NewRequest("POST", "/address", bytes.NewBuffer(jsonPerson))
	response := httptest.NewRecorder()
	CreatePerson(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
	var bn []byte
	if !bytes.Equal(bn, response.Body.Bytes()) {
		t.Errorf("Body Didn't match:\n\tExpected:\t%q\n\tGot:\t%q", bn, response.Body.String())
	}
}
