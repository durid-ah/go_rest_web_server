package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	// handles the sqlite stuff
	_ "github.com/mattn/go-sqlite3"
)

var initialized uint32
var instance *ContactRepo
var mu sync.Mutex


// Contact : The model of each item in the database table
type Contact struct {
	ID			int		
	FirstName 	string	`json:"firstName"`
	LastName 	string	`json:"lastName"`
	PhoneNumber string	`json:"phoneNumber"`
	Email 		string	`json:"email"`
}

// ContactRepo : This is the repository layer of the contact database
type ContactRepo struct {
	database *(sql.DB)
	insertItem *(sql.Stmt)
}

// GetInstance : Get a contact repo singleton instance that is thread safe
func GetInstance() *ContactRepo {

    if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}

    mu.Lock()
    defer mu.Unlock()

    if initialized == 0 {
		instance = &ContactRepo{}
        instance.Initialize()
        atomic.StoreUint32(&initialized, 1)
    }

    return instance
}

// Initialize : Initializes the repository 
func (repo *ContactRepo) Initialize() {
	var err error
	repo.database, err = sql.Open("sqlite3", "./resources/contact.db")

	if (err != nil) {
		fmt.Println(err)
	}

	statement, err := repo.database.Prepare(
		"CREATE TABLE IF NOT EXISTS contact (" +
		"id INTEGER PRIMARY KEY," +
		"firstName TEXT," +
		"lastName TEXT," +
		"phoneNumber TEXT," +
		"email TEXT)")
	
	if err != nil {
		fmt.Println(err)
	} else {
		statement.Exec()
	}
	

	repo.insertItem, _ = repo.database.Prepare(
		"INSERT INTO contact (firstName, lastName, phoneNumber, email) VALUES (?,?,?,?)")
}

// InsertContact : insert a new contact in the database
func (repo *ContactRepo) InsertContact(contact *Contact) {
	repo.insertItem.Exec(
		contact.FirstName,
		contact.LastName,
		contact.PhoneNumber,
		contact.Email,
	)
}

// GetAllContacts : gets the user's contacts from the database
func (repo *ContactRepo) GetAllContacts() *[]*Contact {
	
	contactList := make([] *Contact, 0)

	rows, _ := repo.database.Query(
		"SELECT id, firstName, lastName, phoneNumber, email FROM contact")

	var id int
    var firstname string
	var lastname string
	var email string
	var phoneNumber string

	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname, &phoneNumber, &email)
		contact := &Contact { 
			ID: id,
			FirstName: firstname,
			LastName: lastname, 
			Email: email,
			PhoneNumber: phoneNumber }
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
		contactList = append(contactList, contact)
	}

	return &contactList
}

