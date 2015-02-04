package main

import (
	"reflect"
	"testing"
)

func TestSendOneUser(t *testing.T) {
	// create a user to send
	var person_to_upload = new(PersonUpload)
	var user = new(Person)
	user.Name = "Test User"
	user.Email = "testuser@hotmail.com"
	person_to_upload.User = *user
	person_to_upload.New = true
	// send user to database
	uploaded_person := send_user(*person_to_upload)
	// compare the sent person to the received person
	if !(reflect.DeepEqual(person_to_upload.User.Name,
		uploaded_person.User.Name) && reflect.DeepEqual(
		person_to_upload.User.Email, uploaded_person.User.Email)) {
		t.Error("Failed when sending one person.")
	}
}

func TestTwoUsersSameNameDifferentEmail(t *testing.T) {
	// create users to send
	var person_to_upload = new(PersonUpload)
	person_to_upload.New = true
	var person_to_upload2 = new(PersonUpload)
	person_to_upload2.New = true
	var user = new(Person)
	user.Name = "Same Name Different Email"
	for i := 0; i < 2; i++ {
		if i == 0 {
			user.Email = "samenamedifferentemail@hotmail.com"
			person_to_upload.User = *user
		} else {
			user.Email = "samenamedifferentemail@gmail.com"
			person_to_upload2.User = *user
		}
	}
	// send users to database
	uploaded_person := send_user(*person_to_upload)
	uploaded_person2 := send_user(*person_to_upload2)
	// both users should go through
	if !(reflect.DeepEqual(person_to_upload.User.Name, uploaded_person.User.Name) && reflect.DeepEqual(person_to_upload.User.Email,
		uploaded_person.User.Email) && reflect.DeepEqual(person_to_upload2.User.Name,
		uploaded_person2.User.Name) && reflect.DeepEqual(person_to_upload2.User.Email, uploaded_person2.User.Email)) {
		t.Error("Failed when sending two users with the same name but different emails.")
	}
}

func TestUserWithoutName(t *testing.T) {
	// create user to send
	var person_to_upload = new(PersonUpload)
	person_to_upload.New = true
	var user = new(Person)
	user.Email = "testusernoname@yahoo.ca"
	person_to_upload.User = *user
	// grab the number of rows in the persons table
	var db = initDb()
	var num_of_users = new(Table_Length)
	db.Table("persons").Count(&num_of_users.Length)
	// send user to database
	send_user(*person_to_upload)
	// grab the number of rows in the persons table after attempted upload
	var num_of_users_after_upload = new(Table_Length)
	db.Table("persons").Count(&num_of_users_after_upload.Length)
	// person should not get uploaded
	if (num_of_users.Length != num_of_users_after_upload.Length) {
		t.Error("A user without a name was uploaded.")
	}
}

func TestUserWithoutEmail(t *testing.T) {
	// create user to send
	var person_to_upload = new(PersonUpload)
	person_to_upload.New = true
	var user = new(Person)
	user.Name = "Test User No Email"
	person_to_upload.User = *user
	// grab the number of rows in the persons table
	var db = initDb()
	var num_of_users = new(Table_Length)
	db.Table("persons").Count(&num_of_users.Length)
	// send user to database
	send_user(*person_to_upload)
	// grab the number of rows in the persons table after attempted upload
	var num_of_users_after_upload = new(Table_Length)
	db.Table("persons").Count(&num_of_users_after_upload.Length)
	// person should not get uploaded
	if (num_of_users.Length != num_of_users_after_upload.Length) {
		t.Error("A user without an email was uploaded.")
	}
}

func TestSendTwoUsersSameInfo(t *testing.T) {
	// create a user to send
	var person_to_upload = new(PersonUpload)
	var user = new(Person)
	user.Name = "Test Duplicate User"
	user.Email = "dupeuser@gmail.com"
	person_to_upload.User = *user
	person_to_upload.New = true
    var person_to_upload2 = new(PersonUpload)
    person_to_upload2.User = *user
    person_to_upload2.New = true

	// send user to database
	uploaded_person := send_user(*person_to_upload)
    // get the number of rows in the persons table
    var db = initDb()
    var num_of_users = new(Table_Length)
    db.Table("persons").Count(&num_of_users.Length)
	// send user to database again
	send_user(*person_to_upload2)
    // get the number of rows in the persons table after upload
    var num_of_users_after_upload = new(Table_Length)
    db.Table("persons").Count(&num_of_users_after_upload.Length)

	// uploaded_person should contain test duplicate user and dupeuser@gmail.com
	// num_of_users and num_of_users_after_upload should have the same length value
	if !(reflect.DeepEqual(person_to_upload.User.Name, uploaded_person.User.Name) && reflect.DeepEqual(
		person_to_upload.User.Email, uploaded_person.User.Email)) || (num_of_users.Length != num_of_users_after_upload.Length) {
		t.Error("Sent two users with the same name and email and both went through.")
	}
}

func TestTwoUsersDifferentNameSameEmail(t *testing.T) {
    // create users to send
    var person_to_upload = new(PersonUpload)
    person_to_upload.New = true
    var person_to_upload2 = new(PersonUpload)
    person_to_upload2.New = true
    var user = new(Person)
    user.Email = "diffnamesameemail@mail.utoronto.ca"
    for i := 0; i < 2; i++ {
        if i == 0 {
            user.Name = "Diff Name Same Email"
            person_to_upload.User = *user
        } else {
            user.Name = "Diff Name Same Email THIS USER SHOULD NOT GO THROUGH"
            person_to_upload2.User = *user
        }
    }
    // send users to database
    uploaded_person := send_user(*person_to_upload)
    // get number of rows in persons table after first upload
    var db = initDb()
    var num_of_users = new(Table_Length)
    db.Table("persons").Count(&num_of_users.Length)
    // upload second user and get number of rows after upload
    send_user(*person_to_upload2)
    var num_of_users_after_upload = new(Table_Length)
    db.Table("persons").Count(&num_of_users_after_upload.Length)
    // only the first user should go through
    if !(reflect.DeepEqual(person_to_upload.User.Name, uploaded_person.User.Name) && reflect.DeepEqual(person_to_upload.User.Email,
        uploaded_person.User.Email)) || (num_of_users.Length != num_of_users_after_upload.Length) {
        t.Error("Failed when sending two users with different names but same emails.")
    }
}
