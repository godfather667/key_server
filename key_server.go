// key_server.go - Exposes a Rest API servicing a Key Store Database.
package main

import (
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Person Structure stores data for each person entered into the Database.
type person struct {
	uniqID    int    `json:"uniq_id"`
	firstName string `json:"first_name"`
	lastName  string `json:"last_name"`
	emailAddr string `json:"email_addr"`
	phoneNumb string `json:"phone_numb"`
}

// topID - Keeps Track of the highest ID recorded
var topID int // Top ID in Database

// In-memory Representation of Database
//        - topID lives "at" uniqID zero(0) in the external Keystore
//
// DESIGN - UniqID's of deleted records are "lost", all new records
//        - are assigned topID+1 for their uniqID. The mechanics of
//        - map processing protect against unassigned keys (uniqID's)
var keyStore map[int]person // Keystore Database Structure

// GetBook - Return All Valid Records in Last Name Order
func GetBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Address Book")
}

// Get Persons Address at "uniqID"
func GetPerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Persons Address")
}

// Create Person in Address Book
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Person")
}

// modify Person at "uniqID"
func ModifyPerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Modify Person")
}

// Delete Person at "uniqID"
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Person")
}

// Import Data Base in CSV Format
func ImportCSV(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Import CSV File")
}

// Export Data Base in CSV Format
func ExportCSV(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Export CSV File")
}

func main() {
	fmt.Println("Address Book Server")
	router := mux.NewRouter()

	//Create Keystore map
	keyStore = make(map[int]person)

	// Create Hidden System Record (Record at uniqID Zero(0) and iinitialize topID
	init := person{0, "-first-", "-last-", "-email-", "-phone-"}
	keyStore[0] = init
	topID = 0
	//
	// Setup API EndPoints
	router.HandleFunc("/address", GetBook).Methods("GET")
	router.HandleFunc("/address/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/address/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/address/{id}", ModifyPerson).Methods("PUT")
	router.HandleFunc("/address/{id}", DeletePerson).Methods("DELETE")
	router.HandleFunc("/ImportCSV", ImportCSV).Methods("GET")
	router.HandleFunc("/ExportCSV", ExportCSV).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))

}
