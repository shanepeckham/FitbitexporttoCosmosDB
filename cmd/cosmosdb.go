package cmd

import (
	"os"
	"gopkg.in/mgo.v2"

	"net/url"
	"fmt"
	"strings"
	"time"
	"net"
	"crypto/tls"
	"log"
)

var mongoURL = os.Getenv("MONGOURL")

// MongoDB database and collection names
var mongoDatabaseName = ""
var mongoDBSessionCopy *mgo.Session
var mongoDBSession *mgo.Session
var mongoDBCollection *mgo.Collection
var mongoDBSessionError error


func WriteActivitiesToCosmosDB (activities Activities)  {

	ConnectToMongo()

	log.Println("WriteActivitiesToCosmosDB")

	mongoDBSessionCopy = mongoDBSession.Copy()
	defer mongoDBSessionCopy.Close()

	// Get collection
	mongoDBCollection = mongoDBSessionCopy.DB(mongoDatabaseName).C("me") // change your collection name
	defer mongoDBSessionCopy.Close()

	err = mongoDBCollection.Insert(activities)
	log.Println("Inserted", activities)

	if err != nil {
		log.Fatal("Problem inserting activities for collection: ", err)
	}

}

func WriteSleepToCosmosDB (sleep SleepSummary)  {

	ConnectToMongo()

	log.Println("WriteSleepToCosmosDB")

	mongoDBSessionCopy = mongoDBSession.Copy()
	defer mongoDBSessionCopy.Close()

	// Get collection
	mongoDBCollection = mongoDBSessionCopy.DB(mongoDatabaseName).C("me")
	defer mongoDBSessionCopy.Close()

	err = mongoDBCollection.Insert(sleep)
	log.Println("Inserted", sleep)

	if err != nil {
		log.Fatal("Problem inserting sleep for collection: ", err)
	}

}

func WriteHeartrateToCosmosDB (heartrate HeartRate)  {

	ConnectToMongo()

	log.Println("WriteHeartrateToCosmosDB")

	mongoDBSessionCopy = mongoDBSession.Copy()
	defer mongoDBSessionCopy.Close()

	// Get collection
	mongoDBCollection = mongoDBSessionCopy.DB(mongoDatabaseName).C("me")
	defer mongoDBSessionCopy.Close()

	err = mongoDBCollection.Insert(heartrate)
	log.Println("Inserted", heartrate)

	if err != nil {
		log.Fatal("Problem inserting heartrate for collection: ", err)
	}

}


func ConnectToMongo() {


	if len(os.Getenv("MONGOURL")) == 0 {
		log.Print("The environment variable MONGOURL has not been set")
	} else {
		log.Print("The environment variable MONGOURL is " + os.Getenv("MONGOURL"))
	}

	url, err := url.Parse(mongoURL)
	if err != nil {
		log.Fatal(fmt.Sprintf("Problem parsing Mongo URL %s: ", url), err)
	}


	// Parse the connection string to extract components because the MongoDB driver is peculiar
	var dialInfo *mgo.DialInfo
	mongoUsername := ""
	mongoPassword := ""
	if url.User != nil {
		mongoUsername = url.User.Username()
		mongoPassword, _ = url.User.Password()
		st := fmt.Sprintf("%s", url.User)
		co := strings.Index(st, ":")
		mongoDatabaseName = st[:co]
	}
	mongoHost := url.Host

	mongoDatabase := mongoDatabaseName
	//mongoDatabase := mongoDatabaseName // can be anything
	mongoSSL := strings.Contains(url.RawQuery, "ssl=true")

	log.Printf("\tUsername: %s", mongoUsername)
	log.Printf("\tPassword: %s", mongoPassword)
	log.Printf("\tHost: %s", mongoHost)
	log.Printf("\tDatabase: %s", mongoDatabase)
	log.Printf("\tSSL: %t", mongoSSL)

	if mongoSSL {
		dialInfo = &mgo.DialInfo{
			Addrs:    []string{mongoHost},
			Timeout:  60 * time.Second,
			Database: mongoDatabase, // It can be anything
			Username: mongoUsername, // Username
			Password: mongoPassword, // Password
			DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
				return tls.Dial("tcp", addr.String(), &tls.Config{})
			},
		}
	} else {
		dialInfo = &mgo.DialInfo{
			Addrs:    []string{mongoHost},
			Timeout:  60 * time.Second,
			Database: mongoDatabase, // It can be anything
			Username: mongoUsername, // Username
			Password: mongoPassword, // Password
		}
	}

	success := false
	mongoDBSession, mongoDBSessionError = mgo.DialWithInfo(dialInfo)
	if mongoDBSessionError != nil {
		log.Fatal(fmt.Sprintf("Can't connect to mongo at [%s], go error: ", mongoURL), mongoDBSessionError)
	} else {
		success = true
	}

	if !success {
		os.Exit(1)
	}

	// SetSafe changes the session safety mode.
	// If the safe parameter is nil, the session is put in unsafe mode, and writes become fire-and-forget,
	// without error checking. The unsafe mode is faster since operations won't hold on waiting for a confirmation.
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode.
	mongoDBSession.SetSafe(nil)
}
