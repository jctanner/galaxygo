package galaxy_database

import (
    "fmt"
    "os"

    "database/sql"
    _ "github.com/lib/pq"

    "github.com/jctanner/galaxygo/pkg/galaxy_settings"
)


var settings = galaxy_settings.NewGalaxySettings()


func OpenDatabaseConnection() (*sql.DB, error) {
    /*
    dbhost := os.Getenv("DATABASE_HOST")
    dbname := os.Getenv("DATABASE_NAME")
    dbuser := os.Getenv("DATABASE_USER")
    dbpass := os.Getenv("DATABASE_PASSWORD")
    */

    dbhost := settings.Database_host
    dbname := settings.Database_name
    dbuser := settings.Database_username
    dbpass := settings.Database_password

    // Connection string to the PostgreSQL database
    connStr := "postgres://" + dbuser + ":" + dbpass + "@" + dbhost + ":5432/" + dbname + "?sslmode=disable"

    // Open a connection to the database
    db, err := sql.Open("postgres", connStr)
    if err != nil {
    fmt.Println("Error opening connection:", err)
        return nil, err
    }

    return db, nil
}


func ExecuteQueryWithDatabase(qs string, db *sql.DB) ([]map[string]interface{}, error) {
    rows, err := db.Query(qs)
    if err != nil {
        fmt.Println("Error executing query:", err)
        return nil, err
    }
    defer rows.Close()

    // Get the column names of the query result
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    //fmt.Println(columns)

    // Create a slice to hold the rows as maps
    var results []map[string]interface{}

    // Iterate over the rows
    for rows.Next() {
        // Create a map to hold the values of the current row
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        for i := range columns {
            valuePtrs[i] = &values[i]
        }

        // Scan the values of the current row into the map
        err := rows.Scan(valuePtrs...)
        if err != nil {
            return nil, err
        }

        // Create a map to hold the values of the current row
        rowData := make(map[string]interface{})

        // Copy the values from the slice to the map
        for i, column := range columns {
            rowData[column] = values[i]
        }

        //fmt.Println(rowData)

        // Append the map to the slice
        results = append(results, rowData)
    }

    return results, nil

}

func ExecuteQuery(qs string) ([]map[string]interface{}, error) {

    dbhost := os.Getenv("DATABASE_HOST")
    dbname := os.Getenv("DATABASE_NAME")
    dbuser := os.Getenv("DATABASE_USER")
    dbpass := os.Getenv("DATABASE_PASSWORD")

    // Connection string to the PostgreSQL database
    connStr := "postgres://" + dbuser + ":" + dbpass + "@" + dbhost + ":5432/" + dbname + "?sslmode=disable"

    // Open a connection to the database
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        fmt.Println("Error opening connection:", err)
        return nil, err
    }
    defer db.Close()

    /*
    rows, err := db.Query(qs)
    if err != nil {
        fmt.Println("Error executing query:", err)
        return nil, err
    }
    defer rows.Close()

    // Get the column names of the query result
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    //fmt.Println(columns)

    // Create a slice to hold the rows as maps
    var results []map[string]interface{}

    // Iterate over the rows
    for rows.Next() {
        // Create a map to hold the values of the current row
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        for i := range columns {
            valuePtrs[i] = &values[i]
        }

        // Scan the values of the current row into the map
        err := rows.Scan(valuePtrs...)
        if err != nil {
            return nil, err
        }

        // Create a map to hold the values of the current row
        rowData := make(map[string]interface{})

        // Copy the values from the slice to the map
        for i, column := range columns {
            rowData[column] = values[i]
        }

        //fmt.Println(rowData)

        // Append the map to the slice
        results = append(results, rowData)
    }

    return results, nil
    */

    return ExecuteQueryWithDatabase(qs, db)
}

