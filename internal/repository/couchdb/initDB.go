package couchdb

import (
    "context"
    "log"

    kivik "github.com/go-kivik/kivik/v3"
    _ "github.com/go-kivik/couchdb/v3"
)

// InitDB initializes a connection to the CouchDB instance and creates a database if it doesn't exist
// i use parameter dbName to create a database
func InitDB(dbName string) *kivik.DB {
    client, err := kivik.New("couch", "http://admin:password@localhost:5984/")
    if err != nil {
        log.Fatalf("Failed to connect to CouchDB: %s", err)
    }

    exists, err := client.DBExists(context.TODO(), dbName)
    if err != nil {
        log.Fatalf("Failed to check if database exists: %s", err)
    }

    if !exists {
        err = client.CreateDB(context.TODO(), dbName)
        if err != nil {
            log.Fatalf("Failed to create database: %s", err)
        }
        log.Printf("Database %s created successfully", dbName)
    } else {
        log.Printf("Database %s already exists", dbName)
    }

    db := client.DB(context.TODO(), dbName)
    if err := db.Err(); err != nil {
        log.Fatalf("Failed to connect to the database: %s", err)
    }

    log.Println("Connected to CouchDB and database:", dbName)
    return db
}
