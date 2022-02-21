package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	HOST     = "localhost"
	DATABASE = "UnderTown"
	USER     = "postgres"
	PASSWORD = "12345678"
)

type DBPatient struct {
	db *sql.DB
}

func (s DBPatient) create() {
	_, err := s.db.Exec("CREATE TABLE patient(id INTEGER, name VARCHAR(50), surname VARCHAR(50), patronymic VARCHAR(50), email VARCHAR(50), phoneNumber VARCHAR(50), roomId VARCHAR(50));")
	cherr(err)
	log.Println("database: ", "Finished creating table patient")
}
func (s DBPatient) delete() {
	// Избавляемся от таблицы, если таковая таблица существует
	_, err := s.db.Exec("DROP TABLE IF EXISTS patient;")
	cherr(err)
	log.Print("database: ", "Finished dropping table patient (if existed)")
}

func (s DBPatient) add(NewInfo Patient) {
	// Добавляем информацию в таблицу
	sql_statement := fmt.Sprintf("INSERT INTO patient (id, name, surname, patronymic, email, phoneNumber, roomId) VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s');", NewInfo.Id, NewInfo.Name, NewInfo.Surname, NewInfo.Patronymic, NewInfo.Email, NewInfo.PhoneNumber, NewInfo.RoomId)
	_, err := s.db.Exec(sql_statement)
	cherr(err)
	log.Println("database: ", "Inserted 1 rows of data in patient")
}
func (s DBPatient) isLog(CheckId string) bool {
	sql_statement := fmt.Sprintf("SELECT * FROM patient WHERE phoneNumber = '%s';", CheckId)
	log.Println("database: ", "CheckId data in patient")
	rows, err := s.db.Query(sql_statement)
	cherr(err)
	defer rows.Close()
	return rows.Next()
}

func (s DBPatient) getData(GetId string) Patient {
	sql_statement := fmt.Sprintf("SELECT * FROM patient WHERE phoneNumber = '%s';", GetId)
	log.Println("database: ", "Get data in doctor")
	rows, err := s.db.Query(sql_statement)
	cherr(err)
	defer rows.Close()
	if rows.Next() {
		var dP Patient
		cherr(rows.Scan(&dP.Id, &dP.Name, &dP.Surname, &dP.Patronymic, &dP.Email, &dP.PhoneNumber, &dP.RoomId))
		return dP
	}
	return Patient{}
}

func (s DBPatient) getListId() []string {
	sql_statement := fmt.Sprintf("SELECT phoneNumber FROM patient;")
	log.Println("database: ", "Get data in patient")
	rows, err := s.db.Query(sql_statement)
	cherr(err)
	defer rows.Close()
	var massPh []string
	for rows.Next() {
		var ph string
		cherr(rows.Scan(&ph))
		massPh = append(massPh, ph)
	}
	return massPh
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type DBDoctor struct {
	db *sql.DB
}

func (s DBDoctor) create() {
	_, err := s.db.Exec("CREATE TABLE doctor(id INTEGER, name VARCHAR(50), surname VARCHAR(50), patronymic VARCHAR(50), email VARCHAR(50), phoneNumber VARCHAR(50));")
	cherr(err)
	log.Println("database: ", "Finished creating table doctor")
}
func (s DBDoctor) delete() {
	// Избавляемся от таблицы, если таковая таблица существует
	_, err := s.db.Exec("DROP TABLE IF EXISTS doctor;")
	cherr(err)
	log.Print("database: ", "Finished dropping table doctor (if existed)")
}

func (s DBDoctor) add(NewInfo Doctor) {
	// Добавляем информацию в таблицу
	sql_statement := fmt.Sprintf("INSERT INTO doctor (id, name, surname, patronymic, email, phoneNumber) VALUES (%d, '%s', '%s', '%s', '%s', '%s');",
		NewInfo.Id, NewInfo.Name, NewInfo.Surname, NewInfo.Patronymic, NewInfo.Email, NewInfo.PhoneNumber)
	_, err := s.db.Exec(sql_statement)
	cherr(err)
	log.Println("database: ", "Inserted 1 rows of data in doctor")
}

func (s DBDoctor) isLog(CheckId string) bool {
	sql_statement := fmt.Sprintf("SELECT * FROM doctor WHERE phoneNumber = '%s';", CheckId)
	log.Println("database: ", "CheckId data in doctor")
	rows, err := s.db.Query(sql_statement)
	cherr(err)
	defer rows.Close()
	return rows.Next()
}

func (s DBDoctor) getIn_Name(GetId string) string {
	sql_statement := fmt.Sprintf("SELECT Name FROM doctor WHERE phoneNumber = '%s';", GetId)
	log.Println("database: ", "Get data in doctor")
	rows, err := s.db.Query(sql_statement)
	cherr(err)
	defer rows.Close()
	if rows.Next() {
		var name string
		cherr(rows.Scan(&name))
		return name
	}
	return "NILL"
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type DBRoom struct {
	db *sql.DB
}

func (s DBRoom) create() {
	_, err := s.db.Exec("CREATE TABLE room(roomId VARCHAR(50), temperature VARCHAR(50), humidity VARCHAR(50), illumination VARCHAR(50), energy VARCHAR(50), humans VARCHAR(50), door VARCHAR(50), heater VARCHAR(50), vent VARCHAR(50));")
	cherr(err)
	log.Println("database: ", "Finished creating table room")
}
func (s DBRoom) delete() {
	// Избавляемся от таблицы, если таковая таблица существует
	_, err := s.db.Exec("DROP TABLE IF EXISTS room;")
	cherr(err)
	log.Print("database: ", "Finished dropping table room (if existed)")
}

func (s DBRoom) add(NewInfo Room) {
	// Добавляем информацию в таблицу
	sql_statement := fmt.Sprintf("INSERT INTO room (roomId, temperature, humidity, illumination, energy, humans, door, heater, vent) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');",
		NewInfo.RoomId, NewInfo.Temperature, NewInfo.Humidity, NewInfo.Illumination, NewInfo.Energy,
		NewInfo.Humans, NewInfo.Door, NewInfo.Heater, NewInfo.Vent)
	_, err := s.db.Exec(sql_statement)
	cherr(err)
	log.Println("database: ", "Inserted 1 rows of data in room")
}

func (s DBRoom) getData(GetId string) Room {
	sql_statement := fmt.Sprintf("SELECT * FROM room WHERE roomId = '%s';", GetId)
	log.Println("database: ", "Get data in room")
	rows, err := s.db.Query(sql_statement)
	cherr(err)
	defer rows.Close()
	if rows.Next() {
		var dR Room
		cherr(rows.Scan(&dR.RoomId, &dR.Temperature, &dR.Humidity, &dR.Illumination, &dR.Energy, &dR.Humans, &dR.Door, &dR.Heater, &dR.Vent))
		return dR
	}
	return Room{}
}

func (s DBRoom) update(NewInfo Room) {
	sql_statement := fmt.Sprintf("UPDATE room SET roomId = '%s', temperature = '%s', humidity = '%s', illumination = '%s', energy = '%s', humans = '%s', door = '%s', heater = '%s', vent = '%s' WHERE roomId = '%s';",
		NewInfo.RoomId, NewInfo.Temperature, NewInfo.Humidity, NewInfo.Illumination, NewInfo.Energy, NewInfo.Humans, NewInfo.Door, NewInfo.Heater, NewInfo.Vent, NewInfo.RoomId)
	_, err := s.db.Exec(sql_statement)
	cherr(err)
	log.Println("database: ", "Updated 1 rows of data in room")
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type DBESP8266 struct {
	db *sql.DB
}

func (s DBESP8266) create() {
	_, err := s.db.Exec("CREATE TABLE ESP8266(id serial PRIMARY KEY, room VARCHAR(50));")
	cherr(err)
	log.Println("database: ", "Finished creating table ESP8266")
}
func (s DBESP8266) delete() {
	// Избавляемся от таблицы, если таковая таблица существует
	_, err := s.db.Exec("DROP TABLE IF EXISTS ESP8266;")
	cherr(err)
	log.Print("database: ", "Finished dropping table ESP8266 (if existed)")
}
func (s DBESP8266) check(IDesp interface{}) interface{} {
	n, err := strconv.Atoi(IDesp.(string))
	cherr(err)
	//checkIDesp  - возвращает новый ID если в базе ESP8266 нет входящего IDesp или возвращает IDesp если есть
	sql_statement := fmt.Sprintf("SELECT * FROM ESP8266 WHERE id = %d;", n)
	rows, err := s.db.Query(sql_statement)
	cherr(err)
	defer rows.Close()

	if !rows.Next() {
		sql_statement := "INSERT INTO ESP8266 (room) VALUES ('Some room');"
		_, err := s.db.Exec(sql_statement)
		cherr(err)
		rows1, err := s.db.Query("SELECT id FROM ESP8266 ORDER BY ID DESC LIMIT 1")
		cherr(err)
		for rows1.Next() {
			var id int
			err := rows1.Scan(&id)
			switch err {
			case sql.ErrNoRows:
				log.Println("database: ", "No rows were returned")
			case nil:
				return strconv.Itoa(id)
			default:
				cherr(err)
			}
		}
	}
	return IDesp
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type DBMasterCard struct {
	db *sql.DB
}

func (s DBMasterCard) create() {
	_, err := s.db.Exec("CREATE TABLE MasterCard(name VARCHAR(50));")
	cherr(err)
	log.Println("database: ", "Finished creating table MasterCard")
}
func (s DBMasterCard) delete() {
	// Избавляемся от таблицы, если таковая таблица существует
	_, err := s.db.Exec("DROP TABLE IF EXISTS MasterCard;")
	cherr(err)
	log.Print("database: ", "Finished dropping table MasterCard (if existed)")
}
func (s DBMasterCard) check(Card interface{}) bool {
	//"1d6e126e3707e5ed51c5eea083e82bf4"
	sql_statement := fmt.Sprintf("SELECT * FROM MasterCard WHERE name = '%s';", Card)
	rows, err := s.db.Query(sql_statement)
	cherr(err)
	defer cherr(rows.Close())
	return rows.Next()
}
func (s DBMasterCard) add() interface{} {
	//Добавляем информацию в таблицу
	h := newHash()
	sql_statement := fmt.Sprintf("INSERT INTO MaterCard (name) VALUES ('%s');", h)
	_, err := s.db.Exec(sql_statement)
	cherr(err)
	log.Println("database: ", "Add new MasterCard: ", h)
	return h
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type DBCard struct {
	db *sql.DB
}

func (s DBCard) create() {
	_, err := s.db.Exec("CREATE TABLE Card(name VARCHAR(50));")
	cherr(err)
	log.Println("database: ", "Finished creating table Card")
}
func (s DBCard) delete() {
	// Избавляемся от таблицы, если таковая таблица существует
	_, err := s.db.Exec("DROP TABLE IF EXISTS Card;")
	cherr(err)
	log.Print("database: ", "Finished dropping table Card (if existed)")
}

func (s DBCard) check(Card interface{}) bool {
	//Возвращат true если есть hash в таблице
	sql_statement := fmt.Sprintf("SELECT * FROM Card WHERE name = '%s';", Card)
	rows, err := s.db.Query(sql_statement)
	cherr(err)
	defer cherr(rows.Close())
	return rows.Next()
}
func (s DBCard) add(Card interface{}) {
	//Добавляем информацию в таблицу
	sql_statement := fmt.Sprintf("INSERT INTO Card (name) VALUES ('%s');", Card)
	_, err := s.db.Exec(sql_statement)
	cherr(err)
	log.Println("database: ", "Add new Card: ", Card)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type DBLogCard struct {
	db *sql.DB
}

func (s DBLogCard) create() {
	_, err := s.db.Exec("CREATE TABLE LogCard(name VARCHAR(50), esp VARCHAR(50), time TIMESTAMP);")
	cherr(err)
	log.Println("database: ", "Finished creating table LogCard")
}
func (s DBLogCard) delete() {
	// Избавляемся от таблицы, если таковая таблица существует
	_, err := s.db.Exec("DROP TABLE IF EXISTS LogCard;")
	cherr(err)
	log.Print("database: ", "Finished dropping table LogCard (if existed)")
}
func (s DBLogCard) swipe(IDesp, Card interface{}) {
	sql_statement := fmt.Sprintf("INSERT INTO LogCard (name, esp, time) VALUES ('%s', '%s', '%s');", Card.(string), IDesp.(string), time.Now().Format("2006-01-02 15:04:05.000"))
	_, err := s.db.Exec(sql_statement)
	cherr(err)
	log.Printf("database: ", "Swipe on IDesp: %s, Card: %s", IDesp.(string), Card.(string))
}
func (s DBLogCard) read() {
	var (
		name string
		esp  string
		time string
	)
	sql_statement := "SELECT * FROM LogCard;"
	rows, err := s.db.Query(sql_statement)
	cherr(err)
	defer rows.Close()
	if !rows.Next() {
		fmt.Println("No rows in LogCard")
	}
	for rows.Next() {
		err := rows.Scan(&name, &esp, &time)
		switch err {
		case nil:
			fmt.Printf("Data in LogCard row = (%s, %s, %s)\n", name, esp, time)
		default:
			cherr(err)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func DBConnect() *sql.DB {
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE)
	// Подключаемся к базе данных
	db, err := sql.Open("postgres", connectionString)
	cherr(err)
	cherr(db.Ping())
	log.Println("database: ", "Successfully created connection to database")
	return db
}

func DBTruncate(NameOfData string) {
	db := DBConnect()
	defer db.Close()
	// Очищает информацию во всех колонках, но не сами rows
	sql_statement := fmt.Sprintf("TRUNCATE %s;", NameOfData)
	_, err := db.Exec(sql_statement)
	cherr(err)
	log.Printf("database: ", "Truncate rows of data in %s", NameOfData)
}
