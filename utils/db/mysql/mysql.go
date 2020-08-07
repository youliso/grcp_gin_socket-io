package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"time"
)

var mDb = make(map[string]*sql.DB)

// InitMySQLPool func init DB pool
func InitMysql(DbName, uri string, max, min int) {
	conn, err := sql.Open("mysql", uri)
	if err != nil {
		panic(err.Error())
	}
	//设置连接池
	conn.SetMaxOpenConns(max)
	conn.SetMaxIdleConns(min)
	conn.SetConnMaxLifetime(10 * time.Minute)
	if err := conn.Ping(); err != nil {
		println(err.Error())
	}
	mDb[DbName] = conn
}

// Close pool
func Close(DbName string) error {
	return mDb[DbName].Close()
}

// Query via pool
func Query(DbName, queryStr string, EntityType reflect.Type, args ...interface{}) ([]map[string]interface{}, error) {
	columnToField := make(map[string]string)
	types := EntityType
	for i := 0; i < types.NumField(); i++ {
		typ := types.Field(i)
		tag := typ.Tag
		if len(tag) > 0 {
			column := tag.Get("column")
			name := typ.Name
			columnToField[column] = name
		}
	}
	rows, err := mDb[DbName].Query(queryStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	rowsMap := make([]map[string]interface{}, 0, 10)
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		obj := reflect.New(EntityType).Interface()
		typ := reflect.ValueOf(obj).Elem()
		for i, col := range values {
			if col != nil {
				name := columnToField[columns[i]]
				field := typ.FieldByName(name)
				switch field.Kind() {
				case reflect.String:
					rowMap[name] = string(col.([]byte))
				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					v, _ := strconv.ParseUint(string(col.([]byte)), 10, 0)
					rowMap[name] = v
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v, _ := strconv.ParseInt(string(col.([]byte)), 10, 0)
					rowMap[name] = v
				case reflect.Float32:
					v, _ := strconv.ParseFloat(string(col.([]byte)), 32)
					rowMap[name] = v
				case reflect.Float64:
					v, _ := strconv.ParseFloat(string(col.([]byte)), 64)
					rowMap[name] = v
				case reflect.Struct:
					switch field.Type().String() {
					case "time.Time":
						rowMap[columns[i]] = col
					}
				}
			}
		}
		rowsMap = append(rowsMap, rowMap)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rowsMap, nil
}

func execute(DbName, sqlStr string, args ...interface{}) (sql.Result, error) {
	return mDb[DbName].Exec(sqlStr, args...)
}

// Update via pool
func Update(DbName, updateStr string, args ...interface{}) (int64, error) {
	result, err := execute(DbName, updateStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

// Insert via pool
func Insert(DbName, insertStr string, args ...interface{}) (int64, error) {
	result, err := execute(DbName, insertStr, args...)
	if err != nil {
		return 0, err
	}
	lastid, err := result.LastInsertId()
	return lastid, err
}

// Delete via pool
func Delete(DbName, deleteStr string, args ...interface{}) (int64, error) {
	result, err := execute(DbName, deleteStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

// SQLConnTransaction is for transaction connection
type SQLConnTransaction struct {
	SQLTX *sql.Tx
}

// Begin transaction
func Begin(DbName string) (*SQLConnTransaction, error) {
	var oneSQLConnTransaction = &SQLConnTransaction{}
	var err error
	if pingErr := mDb[DbName].Ping(); pingErr == nil {
		oneSQLConnTransaction.SQLTX, err = mDb[DbName].Begin()
	}
	return oneSQLConnTransaction, err
}

// Rollback transaction
func (t *SQLConnTransaction) Rollback() error {
	return t.SQLTX.Rollback()
}

// Commit transaction
func (t *SQLConnTransaction) Commit() error {
	return t.SQLTX.Commit()
}

// Query via transaction
func (t *SQLConnTransaction) Query(queryStr string, EntityType reflect.Type, args ...interface{}) ([]map[string]interface{}, error) {
	columnToField := make(map[string]string)
	types := EntityType
	for i := 0; i < types.NumField(); i++ {
		typ := types.Field(i)
		tag := typ.Tag
		if len(tag) > 0 {
			column := tag.Get("column")
			name := typ.Name
			columnToField[column] = name
		}
	}
	rows, err := t.SQLTX.Query(queryStr, args...)
	if err != nil {
		panic(err.Error())
		return nil, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	rowsMap := make([]map[string]interface{}, 0, 10)
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		rowMap := make(map[string]interface{})
		obj := reflect.New(EntityType).Interface()
		typ := reflect.ValueOf(obj).Elem()
		for i, col := range values {
			if col != nil {
				name := columnToField[columns[i]]
				field := typ.FieldByName(name)
				switch field.Kind() {
				case reflect.String:
					rowMap[name] = string(col.([]byte))
				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					v, _ := strconv.ParseUint(string(col.([]byte)), 10, 0)
					rowMap[name] = v
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v, _ := strconv.ParseInt(string(col.([]byte)), 10, 0)
					rowMap[name] = v
				case reflect.Float32:
					v, _ := strconv.ParseFloat(string(col.([]byte)), 32)
					rowMap[name] = v
				case reflect.Float64:
					v, _ := strconv.ParseFloat(string(col.([]byte)), 64)
					rowMap[name] = v
				case reflect.Struct:
					switch field.Type().String() {
					case "time.Time":
						rowMap[columns[i]] = col
					}
				}
			}
		}
		rowsMap = append(rowsMap, rowMap)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rowsMap, nil
}

func (t *SQLConnTransaction) execute(sqlStr string, args ...interface{}) (sql.Result, error) {
	return t.SQLTX.Exec(sqlStr, args...)
}

// Update via transaction
func (t *SQLConnTransaction) Update(updateStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(updateStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

// Insert via transaction
func (t *SQLConnTransaction) Insert(insertStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(insertStr, args...)
	if err != nil {
		return 0, err
	}
	lastid, err := result.LastInsertId()
	return lastid, err

}

// Delete via transaction
func (t *SQLConnTransaction) Delete(deleteStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(deleteStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}
