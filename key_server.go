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
	"strings"

	"github.com/gorilla/mux"
)

// Define Application Global Parameters
type App struct {
	Router *mux.Router
}

var a App

// Person Structure stores data for each person entered into the Database.
type Person struct {
	UniqID    string `json:"uniq_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	EmailAddr string `json:"email_addr,omitempty"`
	PhoneNumb string `json:"phone_numb,omitempty"`
}

// Keystore Database Structure Type
type ks map[int]Person

// topID - Keeps Track of the highest ID recorded
var topID int // Top ID in Database

// In-memory Representation of Keystore Database
//        - topID lives "at" uniqID zero(0) in the external Keystore.
//        - KeyStore[0] is a Reserved System Record.
//
// DESIGN - UniqID's of deleted records are "lost", all new records
//        - are assigned topID+1 for their uniqID. The mechanics of
//        - map processing protect against unassigned keys (uniqID's)
//
// LIMITS - The Internal Copy of database is loaded on entry and saved
//          anytime the Internal Copy is modified.
//
//        - Thus Database is safe as long as the program not interrupted
//          during actual database writes. It not perfectly safe, but it is
//          reasonably safe. It is also not safe if used in a concurrent
//          situation.
//
//          To insure Database Integrity Please issue "address/save".
//          That will explicitly save the current in-memory Copy and exit(1).
//
var KeyStore ks

// check - Test for Error and Panic if not nil.
func check(s string, e error) {
	if e != nil {
		fmt.Println(s)
		panic(e)
	}
}

// GetBook - Return All Valid Records in Last Name Order
func GetBook(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(KeyStore)
}

// Get Persons Address at "uniqID"
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	item := params["id"]
	ui, err := strconv.Atoi(item)
	check("String Conversion Failure", err)
	json.NewEncoder(w).Encode(KeyStore[ui].EmailAddr)
}

// Create Person in Address Book
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	check("Json Decorder Failure", err)
	topID += 1
	ui := strconv.Itoa(topID)
	KeyStore[0] = Person{ui, "-first-", "-last-", "-email-", "-phone-"}
	np := Person{ui, p.FirstName, p.LastName, p.EmailAddr, p.PhoneNumb}
	KeyStore[topID] = np
	saveDatabase()
}

// modify Person at "uniqID"
func ModifyPerson(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p Person
	err := decoder.Decode(&p)
	check("Json Decorder Failure", err)
	params := mux.Vars(r)
	item := params["id"]
	ui, err := strconv.Atoi(item)
	check("String Conversion Failure", err)

	cp := KeyStore[ui]
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
	KeyStore[ui] = np
	saveDatabase()
}

// Delete Person at "uniqID"
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	item := params["id"]
	ui, err := strconv.Atoi(item)
	check("String Conversion Failure", err)
	delete(KeyStore, ui)
	saveDatabase()
}

// Import Data Base in CSV Format
func ImportCSV(w http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("Data.csv")
	check("Read of Data.csv Failed! ", err)
	str_buf := strings.Split(string(dat), "\n")

	for k := range KeyStore { // Clear KeyStore
		delete(KeyStore, k)
	}
	init := Person{"0", "-first-", "-last-", "-email-", "-phone-"}
	KeyStore[0] = init
	topID = 1

	for _, v := range str_buf {
		cp := Person{}
		if len(v) != 0 {
			//			fmt.Println("V = ", v)
			d := strings.Split(v, ",")
			for j, f := range d {
				switch j {
				case 0:
					cp.UniqID = f
				case 1:
					cp.FirstName = f
				case 2:
					cp.LastName = f
				case 3:
					cp.EmailAddr = f
				case 4:
					cp.PhoneNumb = f
				default:
					fmt.Println("Default Error")
				}
			}
			KeyStore[topID] = cp
			topID += 1
		}
	}
	init = Person{"0", "-first-", "-last-", "-email-", "-phone-"}
	init.UniqID = strconv.Itoa(topID)
	KeyStore[0] = init
	//	fmt.Println("KS = ", KeyStore)
}

// Export Data Base in CSV Format
func ExportCSV(w http.ResponseWriter, r *http.Request) {
	f, err := os.Create("Data.csv")
	check("Failed to open 'Data.csv'", err)
	defer f.Close()
	var line string
	for _, v := range KeyStore {
		if v.FirstName != "-first-" {
			line = fmt.Sprintf("%v,%v,%v,%v,%v\n", v.UniqID, v.FirstName, v.LastName, v.EmailAddr, v.PhoneNumb)
		}
		//		fmt.Println("KS[", i, "] = ", v)
		_, err := f.WriteString(line)
		check("Write Failed for 'Data.csv'", err)
	}
}

// Save Database
//}
func SaveAddr(w http.ResponseWriter, r *http.Request) {
	saveDatabase()
}

func loadDatabase() {
	data, err := ioutil.ReadFile("Data.db") // Load Database
	if err != nil {                         // If missing - Create
		data, err = json.Marshal(KeyStore) // Marshall Database
		check("Marshalling Failed", err)
		_, err := os.Create("Data.db") // Create Database
		check("Create File Failed", err)
		writeData(data) // Write Database
		fmt.Println("Initialized Database")
	} else {
		err = json.Unmarshal(data, &KeyStore) //Reload In-Memory Copy
		check("Unmarshal Failed", err)
	}
}

// SaveHandler helper to create and store a new page in the database
//
func saveDatabase() {
	data, err := json.Marshal(KeyStore) // Marshal the database
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

//
// Create Empty Database
// Create Hidden System Record (Record at uniqID Zero(0) and iinitialize topID.
//
func CreateDatabase() {
	KeyStore = make(map[int]Person)
	init := Person{"0", "-first-", "-last-", "-email-", "-phone-"}
	KeyStore[0] = init
	topID = 0
	// Initialize or Load Data.db
	saveDatabase()
}

//
// Initialize Router and Endpoints
//

func routeInit() {
	//
	// Initialize Router
	//
	a.Router = mux.NewRouter()

	// Create Initial Database and KeyStore
	CreateDatabase()
	//
	// Setup API EndPoints
	//
	a.Router.HandleFunc("/address", GetBook).Methods("GET")
	a.Router.HandleFunc("/address", CreatePerson).Methods("POST")
	a.Router.HandleFunc("/address/{id}", GetPerson).Methods("GET")
	a.Router.HandleFunc("/address/{id}", ModifyPerson).Methods("PUT")
	a.Router.HandleFunc("/address/{id}", DeletePerson).Methods("DELETE")
	a.Router.HandleFunc("/address/import", ImportCSV).Methods("POST")
	a.Router.HandleFunc("/address/export", ExportCSV).Methods("POST")
	a.Router.HandleFunc("/address/save", SaveAddr).Methods("POST")
}

//
// Main -
//
func main() {
	fmt.Println("Address Book Server")

	// Initialize Router and Endpoints
	routeInit()

	// Listen and Serve at Port 8000
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}
