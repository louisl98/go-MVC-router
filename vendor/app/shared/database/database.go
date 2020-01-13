package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
)

var (
	// SQL wrapper
	SQL *sqlx.DB
	// Database info
	databases Info
)

// Type is the type of database from a Type* constant
type Type string

// TypeMySQL is MySQL
const TypeMySQL Type = "MySQL"

// Info contains the database configurations
type Info struct {
	// Database type
	Type Type
	// MySQL info if used
	MySQL MySQLInfo
}

// MySQLInfo is the details for the database connection
type MySQLInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

// DSN returns the Data Source Name
func DSN(ci MySQLInfo) string {
	// Example: root:@tcp(localhost:3306)/test
	return ci.Username +
		":" +
		ci.Password +
		"@tcp(" +
		ci.Hostname +
		":" +
		fmt.Sprintf("%d", ci.Port) +
		")/" +
		ci.Name + ci.Parameter
}

// Connect to the database
func Connect(d Info) {
	var err error
	// Store the config
	databases = d
	// Connect to MySQL
	if SQL, err = sqlx.Connect("mysql", DSN(d.MySQL)); err != nil {
		log.Println("SQL Driver Error", err)
	}
	// Check if is alive
	if err = SQL.Ping(); err != nil {
		log.Println("Database Error", err)
	}
}

// ReadConfig returns the database information
func ReadConfig() Info {
	return databases
}
