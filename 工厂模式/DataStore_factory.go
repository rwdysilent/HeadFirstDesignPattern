// Copyright 2019 pfwu. All rights reserved.
//
// @Author: pfwu
// @Date: 2019/2/1 15:40

//reference: https://matthewbrown.io/2016/01/23/factory-pattern-in-golang/

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
)

var UserNotFoundError = errors.New("User not found ")

type DataStore interface {
	Name() string
	FindUserNameById(id int64) (string, error)
}

type PostgreSQLDataStore struct {
	DSN string
	DB  sql.DB
}

func (pds *PostgreSQLDataStore) Name() string {
	return "PostgreSQLDataStore"
}

func (pds *PostgreSQLDataStore) FindUserNameById(id int64) (string, error) {
	var username string
	err := pds.DB.Query("SELECT username FROM users WHERE id=$1", id).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", UserNotFoundError
		}
		return "", err
	}
	return username, nil
}

type MemoryDataStore struct {
	sync.RWMutex
	Users map[int64]string
}

func (mds *MemoryDataStore) Name() string {
	return "MemoryDataStore"
}

func (mds *MemoryDataStore) FindUserNameById(id int64) (string, error) {
	mds.RWMutex.RLock()
	defer mds.RWMutex.RUnlock()
	username, ok := mds.Users[id]
	if !ok {
		return "", UserNotFoundError
	}
	return username, nil
}

type DataStoreFactory func(conf map[string]string) (DataStore, error)

func NewPostgreSQLDataStore(conf map[string]string) (DataStore, error) {
	dsn, ok := conf.Get("DATASTORE_POSTGRES_DSN", "")
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s is required for the postgres datastore", "DATASTORE_POSTGRES_DSN"))
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Panicf("Failed to connect to datastore: %s", err.Error())
		return nil, datastore.FailedToConnect
	}

	return &PostgreSQLDataStore{
		DSN: dsn,
		DB:  db,
	}, nil
}

func NewMemoryDataStore(conf map[string]string) (DataStore, error) {
	return &MemoryDataStore{
		Users: map[int64]string{
			1: "mnbbrown",
			0: "root",
		},
		RWMutex: sync.RWMutex{},
	}, nil
}

var dataStoreFactories = make(map[string]DataStoreFactory)

func Register(name string, factory DataStoreFactory) {
	if factory == nil {
		log.Panicf("Datastore factory %s does not exist.", name)
	}
	_, registered := dataStoreFactories[name]
	if registered {
		log.Printf("Datastore factory %s already registered. Ignoring.", name)
		return
	}
	dataStoreFactories[name] = factory
}

func init() {
	Register("postgres", NewPostgreSQLDataStore)
	Register("memory", NewMemoryDataStore)
}

func CreateDataStore(conf map[string]string) (DataStore, error) {

	// Query configuration for datastore defaulting to "memory".
	engineName := conf.Get("DATASTORE", "memory")

	engineFactory, ok := dataStoreFactories[engineName]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableDatastores := make([]string, len(dataStoreFactories))
		for k := range dataStoreFactories {
			availableDatastores = append(availableDatastores, k)
		}
		return nil, errors.New(fmt.Sprintf("Invalid Datastore name. Must be one of: %s", strings.Join(availableDatastores, ", ")))
	}

	// Run the factory with the configuration.
	return engineFactory(conf)
}

func main() {
	dataStore1, err := CreateDataStore(map[string]string{
		"DATASTORE":              "postgres",
		"DATASTORE_POSTGRES_DSN": "dbname=factoriesareamazing",
	})
	if err != nil {
		log.Fatalf("err")
	}
	fmt.Println(dataStore1)

	dataStore2, err := CreateDataStore(map[string]string{
		"DATASTORE": "memory",
	})
	if err != nil {
		log.Fatalf("err")
	}
	fmt.Println(dataStore2)
}
