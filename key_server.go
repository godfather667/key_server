// key_server.go - Exposes a Rest API servicing a Key Store Database.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Person Structure stores data for each person entered into the Database.
type Person struct {
	UniqID    string `json:"uniq_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	EmailAddr string `json:"email_addr,omitempty"`
	PhoneNumb string `json:"phone_numb,omitempty"`
}

// topID - Keeps Track of the highest ID recorded
var topID int // Top ID in Database

// In-memory Representation of Database
//        - topID lives "at" uniqID zero(0) in the external Keystore
//
// DESIGN - UniqID's of deleted records are "lost", all new records
//        - are assigned topID+1 for their uniqID. The mechanics of
//        - map processing protect against unassigned keys (uniqID's)
//
// LIMITS - The Internal Copy of Database is loaded on entry and saved
//          anytime the Internal Copy is Modified.
//        - Thus Database is safe as long as the program not interrupted
//          during actual database writes. It not perfectly safe but it is
//          reasonably safe. It is also not safe if used in a concurrent
//          situation.
//
var keyStore map[int]Person // Keystore Database Structure

// check - Test for Error and Panic if not nil.
func check(s string, e error) {
	if e != nil {
		fmt.Println(s)
		panic(e)
	}
}

// GetBook - Return All Valid Records in Last Name Order
func GetBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Address List = ", keyStore)
	//	saveDatabase()
	fmt.Println("Get Address Book")
}

// Get Persons Address at "uniqID"
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)
	item := params["id"]
	fmt.Println("Item = ", item)
	ui, err := strconv.Atoi(item)
	if err != nil {
		panic(err)
	}
	fmt.Println("Requested Person = ", keyStore[ui])
	fmt.Println("Get Persons Address")
}

// Create Person in Address Book
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}
	topID += 1
	ui := strconv.Itoa(topID)
	keyStore[0] = Person{ui, "-first-", "-last-", "-email-", "-phone-"}
	np := Person{ui, p.FirstName, p.LastName, p.EmailAddr, p.PhoneNumb}
	keyStore[topID] = np
	fmt.Println("Create Person")
	saveDatabase()
	// fmt.Println(keyStore)
}

// modify Person at "uniqID"
func ModifyPerson(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}
	params := mux.Vars(r)
	fmt.Println(params)
	item := params["id"]
	//	fmt.Println("Item = ", item)
	ui, err := strconv.Atoi(item)
	if err != nil {
		panic(err)
	}

	cp := keyStore[ui]
	fmt.Println("CP = ", cp)

	if len(p.FirstName) > 0 {
		cp.FirstName = p.FirstName
	}
	if len(p.LastName) > 0 {
		cp.LastName = p.LastName
	}
	if len(p.EmailAddr) > 0 {
		cp.EmailAddr = p.EmailAddr
	}
	if len(p.PhoneNumb) > 0 {
		cp.PhoneNumb = p.PhoneNumb
	}
	np := Person{cp.UniqID, cp.FirstName, cp.LastName, cp.EmailAddr, cp.PhoneNumb}
	keyStore[ui] = np
	saveDatabase()
	//	fmt.Println("NP = ", np, "UI = ", ui)
	//  fmt.Println("KeyStore = ", keyStore[ui])
	fmt.Println("Modify Person")
}

// Delete Person at "uniqID"
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)
	item := params["id"]
	//	fmt.Println("Item = ", item)
	ui, err := strconv.Atoi(item)
	if err != nil {
		panic(err)
	}
	delete(keyStore, ui)
	saveDatabase()
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

func loadDatabase() {
	//	var xMem []Page // Temporary Blank in-memory Database
	data, err := ioutil.ReadFile("Data.db") // Load Database
	if err != nil {                         // If missing - Create
		data, err = json.Marshal(keyStore) // Marshall Database
		check("Marshalling Failed", err)
		_, err := os.Create("Data.db") // Create Database
		check("Create File Failed", err)
		writeData(data) // Write Database
		fmt.Println("No Database Found - Creating New Empty Database!")
	} else {
		err = json.Unmarshal(data, &keyStore) //Reload In-Memory Copy
		check("Unmarshal Failed", err)
	}
}

// SaveHandler helper to create and store a new page in the database
//
func saveDatabase() {
	data, err := json.Marshal(keyStore) // Marshal the database
	check("Marshalling Failed", err)    // Check for error
	writeData(data)                     // Write database to the disk
	return                              // Return
}

//
// Write Data Set to Disk
//
func writeData(data []byte) { // Write "Mashalled" data to external device
	err := ioutil.WriteFile("Data.db", data, 0644)
	check("Write File Failed", err) // Error Check -- Panic if write fails
}

func main() {
	fmt.Println("Address Book Server\n")
	router := mux.NewRouter()

	//Create Keystore map
	keyStore = make(map[int]Person)

	// Create Hidden System Record (Record at uniqID Zero(0) and iinitialize topID
	init := Person{"0", "-first-", "-last-", "-email-", "-phone-"}
	keyStore[0] = init
	topID = 0
	//
	// Initialize or Load Data.db
	loadDatabase()
	//
	// Setup API EndPoints
	router.HandleFunc("/address", GetBook).Methods("GET")
	router.HandleFunc("/address/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/address", CreatePerson).Methods("POST")
	router.HandleFunc("/address/{id}", ModifyPerson).Methods("PUT")
	router.HandleFunc("/address/{id}", DeletePerson).Methods("DELETE")
	router.HandleFunc("/ImportCSV", ImportCSV).Methods("GET")
	router.HandleFunc("/ExportCSV", ExportCSV).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))

}
