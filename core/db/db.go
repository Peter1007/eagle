package db

import (
	"database/sql"
	"eagle/core"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v1"
	"io/ioutil"
)

type DbConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Dbname   string
	Charset  string
}

type DbConnection struct {
	connKey    string
	queryTimes int
	conn       *sql.DB
}

const (
	MAX_OPEN_CONN  = 200
	MAX_IDLE_CONN  = 100
	BAD_CONNECTION = "driver: bad connection"
)

var dbConfigs map[string]DbConfig
var dbConnections map[string]*DbConnection

func InitDb() error {
	dbConfigs = map[string]DbConfig{}

	filename := core.GetAppPath() + "/config/db.yml"

	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buffer, &dbConfigs)
	if err != nil {
		return err
	}

	dbConnections = map[string]*DbConnection{}
	for connKey, _ := range dbConfigs {
		dbConnections[connKey] = &DbConnection{connKey: connKey, queryTimes: 0, conn: nil}
		err = dbConnections[connKey].connect()
		if err != nil {
			return err
		}
	}

	return err
}

func GetConnection(connKey string) (*DbConnection, error) {
	dbConn, ok := dbConnections[connKey]
	if ok {
		return dbConn, nil
	} else {
		return nil, errors.New("Connect " + connKey + " not exist")
	}
}

func (dbConn *DbConnection) connect() (err error) {
	config := dbConfigs[dbConn.connKey]
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", config.Username, config.Password, config.Host, config.Port, config.Dbname, config.Charset)

	dbConn.conn, err = sql.Open("mysql", connStr)
	dbConn.conn.SetMaxOpenConns(MAX_OPEN_CONN)
	dbConn.conn.SetMaxIdleConns(MAX_IDLE_CONN)
	if err == nil {
		err = dbConn.conn.Ping()
	}

	return
}

func (dbConn *DbConnection) reConnect() (err error) {
	err = dbConn.Close()
	if err != nil && err.Error() != BAD_CONNECTION {
		return
	}

	dbConn.conn = nil
	dbConn.queryTimes = 0

	err = dbConn.connect()

	return
}

func (dbConn *DbConnection) Close() error {
	return dbConn.conn.Close()
}

func (dbConn *DbConnection) Select(sql string) ([]map[string]string, error) {
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	result := []map[string]string{}
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		result = append(result, record)
	}

	return result, err
}

func (dbConn *DbConnection) SelectOne(sql string) (map[string]string, error) {
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	result := map[string]string{}
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		for i, col := range values {
			if col != nil {
				result[columns[i]] = string(col.([]byte))
			}
		}

		break
	}

	return result, err
}

func (dbConn *DbConnection) Insert(sql string, args ...interface{}) (int64, error) {
	result, err := dbConn.prepareAndExec(sql, args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (dbConn *DbConnection) Update(sql string, args ...interface{}) (int64, error) {
	result, err := dbConn.prepareAndExec(sql, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (dbConn *DbConnection) Delete(sql string, args ...interface{}) (int64, error) {
	result, err := dbConn.prepareAndExec(sql, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (dbConn *DbConnection) Query(sql string) (*sql.Rows, error) {
	rows, err := dbConn.conn.Query(sql)
	if err != nil && err.Error() == BAD_CONNECTION {
		err = dbConn.reConnect()
		if err == nil {
			rows, err = dbConn.conn.Query(sql)
		}
	}

	return rows, err
}

func (dbConn *DbConnection) prepareAndExec(sql string, args ...interface{}) (sql.Result, error) {
	stmt, err := dbConn.conn.Prepare(sql)
	if err != nil && err.Error() == BAD_CONNECTION {
		err = dbConn.reConnect()
		if err != nil {
			return nil, err
		}

		stmt, err = dbConn.conn.Prepare(sql)
		if err != nil {
			return nil, err
		}
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}
