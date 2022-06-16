//go:build test
// +build test

// this file will not be included unless it's built for testing
package db

// this function is very dangerous because it drops the entire database
// (which is a functionality we'd like to have for running tests)
func DropDB(dbName string) {
	mdbInstance.Client.Database(dbName).Drop(mdbInstance.Ctx)
}
