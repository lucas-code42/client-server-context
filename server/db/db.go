package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/desafio/sever/api/model"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite3Db struct {
	db *sql.DB
}

func InitSqlite3(conn *sql.DB) *Sqlite3Db {
	return &Sqlite3Db{db: conn}
}

func ConnectSqlite3() (*Sqlite3Db, error) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Println("could not connect with sqlite3")
		return &Sqlite3Db{}, err
	}

	if err = db.Ping(); err != nil {
		log.Println("error to ping sqlite3")
		return &Sqlite3Db{}, err
	}

	sqliteConnection := InitSqlite3(db)
	if err = sqliteConnection.CreateTable(); err != nil {
		return &Sqlite3Db{}, err
	}

	return sqliteConnection, nil
}

func (s *Sqlite3Db) CreateTable() error {
	statement := `
		CREATE TABLE IF NOT EXISTS dolarPrice(
			code VARCHAR(32),
			codein VARCHAR(32),
			name VARCHAR(255),
			high VARCHAR(32),
			low VARCHAR(32),
			pctChange VARCHAR(32),
			bid VARCHAR(32),
			ask VARCHAR(32),
			timestamp VARCHAR(32),
			create_date VARCHAR(32)
		)
	`
	_, err := s.db.Exec(statement)
	if err != nil {
		log.Printf("could not create database")
		log.Fatal(err)
		return err
	}
	return nil
}

func (s *Sqlite3Db) Save(data model.Awesomeapi) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	query := "INSERT INTO dolarPrice(code, codein, name, high, low, pctChange, bid, ask, timestamp, create_date) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	result, resultErr := stmt.ExecContext(
		ctx,
		data.Code,
		data.Codein,
		data.Name,
		data.High,
		data.Low,
		data.PctChange,
		data.Bid,
		data.Ask,
		data.Timestamp,
		data.CreateDate,
	)
	if resultErr != nil {
		return resultErr
	}

	_, err = result.RowsAffected()
	if err != nil {
		return resultErr
	}

	select {
	case <-ctx.Done():
		log.Println("could not persit data on context time")
		return resultErr
	case <-time.After(10 * time.Millisecond):
		log.Println("save data successfully")
		return nil
	}

}
