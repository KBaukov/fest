package db

import (
	_ "encoding/json"
	"fest/config"
	"fmt"
	_ "math"
	_ "reflect"
	_ "sort"
	_ "strings"
	"time"

	//"errors"
	"github.com/jmoiron/sqlx"
	"log"

	_ "github.com/lib/pq"
)

var (
	cfg             config.Configuration
	sessionDuration string
)

const (
	getSessionQuery       = `select * from sessions where device_id=$1 and login = $2;`
	ceateSessionQuery     = `insert into sessions ( device_id, login, token, expier_at ) values ($1,$2,$3, (now() + INTERVAL '%s') );`
	updateSessionQuery    = `update sessions set expier_at= (now() + INTERVAL '%s'), token=$1  where device_id =$2;`
	deleteSessionQuery    = `delete from sessions where device_id=$1;`
	selectDieSessionQuery = `select device_id from sessions where expier_at <= now();`
	deleteDieSessionQuery = `delete from sessions where expier_at <= now();`
	creatUserQuery        = `insert into users ( name, email ) values ($1,$2);`

	authQuery   = `select token from users where email = $1;`
	tokenUpdate = `update users set token = $1 where email = $2;`
)

// database структура подключения к базе данных
type Database struct {
	Conn *sqlx.DB
}

func init() {
	cfg, _ = config.LoadConfig("config.json")
	sessionDuration = cfg.SessionSettings.SessionDuration
}

// dbService представляет интерфейс взаимодействия с базой данных
type DbService interface {
	//GetLastId(table string) (int, error)
	//GetSeats() ([]ent.Seat, error)
	//GetSeatStates() ([]ent.SeatState, error)
	//GetSeatsInfo() ([]ent.SeatInfo, error)
	//SetSeatStates() error
	//GetEventTarifs() ([]ent.Tafif, error)
	//CheckSeatStatess() (bool, error)
	//ReserveSeat() (bool, error)
	//UnReserveSeat() (bool, error)
	//CreateOrder() (string, string, float32, []*ent.Item, error)
	//GetNewOrderNumber() (string, error)
	//CalculateOrderAmount() (float32, []*ent.Item, error)
	//ClearExpiredReserves() ([]string, error)
	//
	//OrderLog() (bool, error)
	//
	//SendTickets() (bool, error)
	//GetLastTicketNumber() (string, error)
}

// newDB открывает соединение с базой данных
func NewDB(connectionString string) (Database, error) {
	dbConn, err := sqlx.Open("postgres", connectionString)
	log.Println(connectionString)
	return Database{Conn: dbConn}, err
}

// #################################################################
func (db Database) CreateUser(name string, email string) (bool, error) {

	stmt, err := db.Conn.Prepare(creatUserQuery)
	if err != nil {
		log.Println("Error while create user:  = %v", err)
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, email)
	if err != nil {
		log.Println("Error while create user:  = %v", err)
		return false, err
	}

	return true, nil
}

func (db Database) AuthUser(login string) (bool, error) {
	var t string
	stmt, err := db.Conn.Prepare(authQuery)
	if err != nil {
		log.Println("Error while auth user:  = %v", err)
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(login)
	if err != nil {
		log.Println("Error while auth user:  = %v", err)
		return false, err
	}
	for rows.Next() {

		err = rows.Scan(&t)
		if err != nil {
			log.Println("Error while auth user:  = %v", err)
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (db Database) CreateSession(deviceId string, login string, token string) (bool, error) {

	stmt, err := db.Conn.Prepare(fmt.Sprintf(ceateSessionQuery, sessionDuration))
	if err != nil {
		log.Println("Error while create session:  = %v", err)
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(deviceId, login, token)
	if err != nil {
		log.Println("Error while create session:  = %v", err)
		return false, err
	}

	return true, nil
}

func (db Database) CheckSession(deviceId string, login string, token string) (bool, error) {

	stmt, err := db.Conn.Prepare(getSessionQuery)
	if err != nil {
		log.Println("Error while check session:  = %v", err)
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(deviceId)
	if err != nil {
		log.Println("Error while create session:  = %v", err)
		return false, err
	}
	for rows.Next() {
		var (
			dId    string
			ll     string
			tt     string
			exTime time.Time
		)
		err = rows.Scan(&dId, &ll, &tt, &exTime)
		if err != nil {
			log.Println("Error while check session:  = %v", err)
			return false, err
		}
		if ll != login || tt != token || exTime.Before(time.Now()) {
			//DELETE seesion object
			stmt, err = db.Conn.Prepare(deleteSessionQuery)
			if err != nil {
				log.Println("Error while delete session:  = %v", err)
				return false, err
			}
			_, err = stmt.Exec(deviceId)
			if err != nil {
				log.Println("Error while delete session:  = %v", err)
				return false, err
			}
			return false, nil
		}
		return true, nil
	}

	return false, nil
}

func (db Database) UpdateSession(deviceId string, token string, login string, isCheck bool) (bool, error) {

	var (
		dId    string
		ll     string
		tt     string
		exTime time.Time
	)

	if isCheck {
		stmt, err := db.Conn.Prepare(getSessionQuery)
		if err != nil {
			log.Println("Error while check session:  = %v", err)
			return false, err
		}
		defer stmt.Close()

		rows, err := stmt.Query(deviceId)
		if err != nil {
			log.Println("Error while create session:  = %v", err)
			return false, err
		}
		for rows.Next() {

			err = rows.Scan(&dId, &ll, &tt, &exTime)
			if err != nil {
				log.Println("Error while check session:  = %v", err)
				return false, err
			}
			break
		}
	}

	stmt, err := db.Conn.Prepare(fmt.Sprintf(updateSessionQuery, sessionDuration))
	if err != nil {
		log.Println("Error while update session:  = %v", err)
		return false, err
	}
	_, err = stmt.Exec(token, deviceId)
	if err != nil {
		log.Println("Error while update session:  = %v", err)
		return false, err
	}
	return true, nil
}

func (db Database) ClearDieSession() ([]string, error) {
	devices := make([]string, 0)
	stmt, err := db.Conn.Prepare(selectDieSessionQuery)
	if err != nil {
		log.Println("Error while get Die session:  = %v", err)
		return devices, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Println("Error while create session:  = %v", err)
		return devices, err
	}
	for rows.Next() {
		var dId string

		err = rows.Scan(&dId)
		if err != nil {
			log.Println("Error while get die session:  = %v", err)
			return devices, err
		}
		devices = append(devices, dId)
	}
	stmt, err = db.Conn.Prepare(deleteDieSessionQuery)
	if err != nil {
		log.Println("Error while delete die session:  = %v", err)
		return devices, err
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Println("Error while delete die session:  = %v", err)
		return devices, err
	}
	return devices, nil
}

//func (db Database) GetMaxMinZoneTarifs() ([]*ent.MaxMinTafif, error) {
//	stt := make([]*ent.MaxMinTafif, 0)
//	stmt, err := db.Conn.Prepare(selectZoneMaxMin)
//	if err != nil {
//		log.Println("Error while max t values:  = %v", err)
//		return stt, err
//	}
//	defer stmt.Close()
//	rows, err := stmt.Query()
//	if err != nil {
//		log.Println("Error while max t values:  = %v", err)
//		return stt, err
//	}
//	for rows.Next() {
//		var (
//			zone int
//			name string
//			max  int
//			min  int
//		)
//
//		err = rows.Scan(&zone, &name, &max, &min)
//		st := ent.MaxMinTafif{zone, name, max, min}
//		stt = append(stt, &st)
//	}
//
//	return stt, nil
//}
