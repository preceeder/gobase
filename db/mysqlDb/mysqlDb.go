package mysqlDb

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"strings"
)

var MysqlDb *Sdb

type MysqlConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Password    string `json:"password"`
	User        string `json:"user"`
	Db          string `json:"db"`
	MaxOpenCons int    `json:"maxOpenCons"`
	MaxIdleCons int    `json:"MaxIdleCons"`
}

// 使用 viper读取的配置初始化
func InitMysqlWithViperConfig(config viper.Viper) {
	redisConfig := readMysqlConfig(config)
	initMySQL(redisConfig)
}
func readMysqlConfig(config viper.Viper) map[string]MysqlConfig {
	mysql := config.Sub("mysql")
	if mysql == nil {
		fmt.Printf("mysqlDb config is nil")
		os.Exit(1)
	}
	var subMysql map[string]MysqlConfig
	err := mysql.Unmarshal(&subMysql)
	if err != nil {
		fmt.Printf("mysqlDb config read error: " + err.Error())
		os.Exit(1)
	}
	return subMysql
}

func InitMysqlWithStruct(config map[string]MysqlConfig) {
	initMySQL(config)
}

// 初始化数据库
func initMySQL(config map[string]MysqlConfig) {
	//mysqlConfig := readMysqlConfig(config)

	MysqlDb = &Sdb{
		Db:        make(Mdb),
		DefaultDb: "default",
	}

	for key, v := range config {
		//dsn := "root:password@tcp(127.0.0.1:3306)/database"
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", v.User, v.Password, v.Host, v.Port, v.Db)
		slog.Info("链接数据库", "db", dsn)
		db := sqlx.MustConnect("mysql", dsn)
		db.SetMaxOpenConns(v.MaxOpenCons)
		db.SetMaxIdleConns(v.MaxIdleCons)
		MysqlDb.Db[key] = db
	}

	return
}

type Mdb map[string]*sqlx.DB

type Sdb struct {
	Db        Mdb
	DefaultDb string
}

func MysqlPoolClose() {
	if MysqlDb == nil {
		return
	}
	for _, ml := range MysqlDb.Db {
		slog.Info("close mysqlDb")
		ml.Close()
	}
}

func (s Sdb) defaultDb() string {
	return s.DefaultDb
}

func (s Sdb) Get(dest interface{}, query string, args ...interface{}) (err error) {
	err = s.Db[s.DefaultDb].Get(dest, query, args...)
	return
}

func (s Sdb) getDb(params map[string]any) (db string) {
	if adb, ok := params["DB"]; !ok {
		db = s.DefaultDb
	} else {
		db, _ = adb.(string)
		// 不存在就是一个 空操作
		delete(params, "DB")
	}
	return
}

func (s Sdb) getTableName(ctx utils.Context, params map[string]any) (tableName string) {
	if tn, ok := params["tableName"]; ok {
		tableName = tn.(string)
		delete(params, "tableName")
	} else {
		slog.Error("mysqlDb not table name", "requestId", ctx.RequestId, "params", params)
		panic("not table name")
	}
	return
}

// 参数解析
func (s Sdb) sqlPares(ctx utils.Context, osql string, params map[string]any) (sql string, args []any, db string) {
	db = s.getDb(params)
	sql, args, err := sqlx.Named(osql, params)
	if err != nil {
		slog.Error("sqlx.Named error :"+err.Error(), "requestId", ctx.RequestId)
		panic("sqlx.Named error :" + err.Error())
	}
	sql, args, err = sqlx.In(sql, args...)
	if err != nil {
		slog.Error("sqlx.In error :"+err.Error(), "requestId", ctx.RequestId)
		panic("sqlx.In error :" + err.Error())
	}
	sql = s.Db[db].Rebind(sql)
	return sql, args, db
}

func (s Sdb) Select(ctx utils.Context, sqlStr string, params map[string]any, row any) bool {
	q, args, db := s.sqlPares(ctx, sqlStr, params)
	err := s.Db[db].Get(row, q, args...)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false
	case err != nil:
		slog.Error("mysqlDb Query failed", "error", err, "sql", sqlStr, "data", params, "requestId", ctx.RequestId)
		return false
	}
	return true
}

func (s Sdb) Fetch(ctx utils.Context, sqlStr string, params map[string]any, row any) bool {
	q, args, db := s.sqlPares(ctx, sqlStr, params)
	rs, err := s.Db[db].Queryx(q, args...)
	if err != nil {
		slog.Error("mysqlDb Fetch failed", "error", err, "sql", sqlStr, "data", params, "requestId", ctx.RequestId)
		return false
	}
	defer rs.Close()
	err = sqlx.StructScan(rs, row)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false
	case err != nil:
		slog.Error("mysqlDb Fetch StructScan failed", "error", err, "sql", sqlStr, "data", params, "requestId", ctx.RequestId)
		return false
	}
	return true
}

// sqlStr="select * from t_user where userId=?" agrs: []any{"2222222"}
func (s Sdb) FetchByArgs(ctx utils.Context, sqlStr string, args []any, db string, row any) bool {
	if db == "" {
		db = s.DefaultDb
	}
	err := sqlx.Select(s.Db[db], row, sqlStr, args...)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false
	case err != nil:
		slog.Error("mysqlDb Fetch StructScan failed", "error", err, "sql", sqlStr, "data", args, "requestId", ctx.RequestId)
		return false
	}
	return true
}

// map[string]any{"DB": "default", "tableName": "t_user",  "Set":map[string]any{"nick": "nihao"}, "Where":map[string]any{"userId": "1111"}}

func (s Sdb) Update(ctx utils.Context, params map[string]any, tx ...*sqlx.Tx) int64 {
	tableName := s.getTableName(ctx, params)

	var sqlStr string = "update " + tableName
	var tempParams = make(map[string]any, 0)
	setValues := params["Set"].(map[string]any)
	setL := make([]string, 0)
	for k, v := range setValues {
		tpv := ""
		if vt, ok := v.([]string); ok {
			tpv = k + " = " + vt[0]
		} else {
			tpv = k + "=" + " :" + k
			tempParams[k] = v
		}
		setL = append(setL, tpv)
	}

	whereValues := params["Where"].(map[string]any)
	wvL := make([]string, 0)
	var tpv string
	for k, v := range whereValues {
		if _, ok := v.([]string); ok {
			tpv = "`" + k + "`" + " in" + "( :" + k + " )"
		} else if _, ok = v.([]int); ok {
			tpv = "`" + k + "`" + " in" + " (:" + k + ")"
		} else {
			tpv = "`" + k + "`" + "=" + " :" + k
		}
		tempParams[k] = v
		wvL = append(wvL, tpv)
	}
	sqlStr = sqlStr + " set " + strings.Join(setL, ", ") + " where " + strings.Join(wvL, " and ")
	q, args, db := s.sqlPares(ctx, sqlStr, tempParams)
	var rs sql.Result
	var err error
	if len(tx) > 0 && tx[0] != nil {
		rs, err = tx[0].Exec(q, args...)
	} else {
		rs, err = s.Db[db].Exec(q, args...)
	}

	if err != nil {
		paS, _ := sonic.MarshalString(args)
		slog.Error("mysqlDb update failed", "error", err, "sql", q, "data", paS, "requestId", ctx.RequestId)
		return -1
	}
	aft, _ := rs.RowsAffected()
	return aft
}

//map[string]any{"DB": "", "tableName":"t_user", "name": "nick", "id": 1}

func (s Sdb) Insert(ctx utils.Context, params map[string]any, tx ...*sqlx.Tx) sql.Result {
	db := s.getDb(params)
	tableName := s.getTableName(ctx, params)
	var sqlStr string = "insert into " + tableName
	var attrs = []string{}
	var attrValues = []string{}
	for k, _ := range params {
		attrs = append(attrs, "`"+k+"`")
		value := ":" + k
		attrValues = append(attrValues, value)
	}
	attrString := strings.Join(attrs, ", ")
	attrValuesString := strings.Join(attrValues, ", ")
	sqlStr = sqlStr + "( " + attrString + ")" + "  values( " + attrValuesString + " )"

	var rs sql.Result
	var err error
	if len(tx) > 0 && tx[0] != nil {
		rs, err = tx[0].NamedExec(sqlStr, params)
	} else {
		rs, err = s.Db[db].NamedExec(sqlStr, params)
	}
	if err != nil {
		slog.Error("mysqlDb insert failed", "error", err.Error(), "sql", sqlStr, "data", params, "requestId", ctx.RequestId)
		return nil
	}

	return rs
}

// map[string]any{"DB": "", "tableName":"t_user",  "Set":[{"name": "nick", "id": 1}] || {{"name": "nick", "id": 1}}}
// 这个支持插入单挑数据  或多条数据
func (s Sdb) InsertMany(ctx utils.Context, params map[string]any, tx ...*sqlx.Tx) sql.Result {
	db := s.getDb(params)
	tableName := s.getTableName(ctx, params)
	var sqlStr string = "insert into " + tableName
	var attrs = []string{}
	allValues := params["Set"] // []map[string]any
	finalVStr := ""
	if ap, ok := allValues.(map[string]any); ok {
		var attrValues = []string{}
		for k, _ := range ap {
			attrs = append(attrs, "`"+k+"`")
			value := ":" + k
			attrValues = append(attrValues, value)
		}
		finalVStr = "(" + strings.Join(attrValues, ", ") + ")"

	} else if aps, ok := allValues.([]any); ok {
		for _, ap := range aps {
			var attrValues = []string{}
			insertData, ok := ap.(map[string]any)
			if !ok {
				slog.Error("mysql Parameter error", "data", allValues, "need", "[]any", "requestId", ctx.RequestId)
				break
			}
			for k, _ := range insertData {
				attrs = append(attrs, "`"+k+"`")
				value := ":" + k
				attrValues = append(attrValues, value)
			}
			finalVStr = "(" + strings.Join(attrValues, ", ") + ")"
			break
		}
	}

	attrString := strings.Join(attrs, ", ")
	sqlStr = sqlStr + "( " + attrString + ")" + "  values " + finalVStr

	var rs sql.Result
	var err error
	if len(tx) > 0 && tx[0] != nil {
		rs, err = tx[0].NamedExec(sqlStr, allValues)
	} else {
		rs, err = s.Db[db].NamedExec(sqlStr, allValues)
	}
	if err != nil {
		slog.Error("mysqlDb insert failed", "error", err.Error(), "sql", sqlStr, "data", allValues, "requestId", ctx.RequestId)
		return nil
	}

	return rs
}

//map[string]any{"DB": "", "tableName":"t_user", "Set":map[string]any, "Update":map[string]any}

func (s Sdb) InsertUpdate(ctx utils.Context, params map[string]any, tx ...*sqlx.Tx) sql.Result {
	db := s.getDb(params)
	tableName := s.getTableName(ctx, params)
	var sqlStr string = "insert into " + tableName

	setValues := params["Set"].(map[string]any)
	var attrs = []string{}
	var attrValues = []string{}
	for k, _ := range setValues {
		attrs = append(attrs, "`"+k+"`")
		value := ":" + k
		attrValues = append(attrValues, value)
	}
	attrString := strings.Join(attrs, ", ")
	attrValuesString := strings.Join(attrValues, ", ")
	sqlStr = sqlStr + "( " + attrString + ")" + "  values( " + attrValuesString + " )"

	var UpdateL = make([]string, 0)
	if uValues, ok := params["Update"].(map[string]any); ok {
		for k, v := range uValues {
			tpv := ""
			if vt, ok := v.([]string); ok {
				tpv = k + " = " + vt[0]
			} else {
				tpv = k + "=values(" + k + ")"
			}
			UpdateL = append(UpdateL, tpv)
		}
	} else if uValues, ok := params["Update"].([]string); ok {
		for _, name := range uValues {
			tpv := ""
			tpv = name + "=values(" + name + ")"
			UpdateL = append(UpdateL, tpv)
		}
	}
	sqlStr += " on duplicate key update " + strings.Join(UpdateL, ",")
	var rs sql.Result
	var err error
	if len(tx) > 0 && tx[0] != nil {
		rs, err = tx[0].NamedExec(sqlStr, setValues)
	} else {
		rs, err = s.Db[db].NamedExec(sqlStr, setValues)
	}
	if err != nil {
		slog.Error("mysqlDb insert failed", "error", err.Error(), "sql", sqlStr, "data", params, "requestId", ctx.RequestId)
		return nil
	}

	return rs
}

func (s Sdb) Execute(ctx utils.Context, sqlStr string, params map[string]any, tx ...*sqlx.Tx) sql.Result {
	//不能做查询， 这里是没有返回结果的
	q, args, db := s.sqlPares(ctx, sqlStr, params)
	var rs sql.Result
	var err error
	if len(tx) > 0 && tx[0] != nil {
		rs, err = tx[0].Exec(q, args...)
	} else {
		rs, err = s.Db[db].Exec(q, args...)
	}
	if err != nil {
		slog.Error("mysqlDb Execute failed", "error", err, "sql", q, "data", params, "requestId", ctx.RequestId)
		return nil
	}
	return rs
}

// []map[string]any{{"sql": "", "params":{}}}

//func (s Sdb) DefaultTransactions(sqlObj []map[string]any) {
//	s.Transactions("default", sqlObj)
//}

//func (s Sdb) Transactions(db string, sqlObj []map[string]any) {
//
//	tx, err := s.Db[db].BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
//	var IsCommit bool = false
//	defer func() {
//		// 无条件提交
//		if !IsCommit {
//			tx.Commit()
//		}
//	}()
//	if err != nil {
//		panic("mysqlDb tx error: " + err.Error())
//	}
//	for _, sqlJ := range sqlObj {
//		sqlStr, ok := sqlJ["sql"]
//		if !ok {
//			panic("Transactions not sqlstr ")
//		}
//		sqlParms, ok := sqlJ["params"]
//		var q string
//		var args []any
//		if ok {
//			tmp := sqlParms.(map[string]any)
//			q, args, _ = s.sqlPares(sqlStr.(string), tmp)
//		} else {
//			q, args, _ = s.sqlPares(sqlStr.(string), nil)
//		}
//		_, execErr := tx.Exec(q, args...)
//		if execErr != nil {
//			_ = tx.Rollback()
//			IsCommit = true
//			slog.Error("mysqlDb Transactions execute failed", "error", execErr.Error())
//			return
//		}
//	}
//	if err := tx.Commit(); err != nil {
//		slog.Error("mysqlDb Transactions commit failed", "error", err.Error())
//
//	}
//	IsCommit = true
//
//}

func (s Sdb) Transaction(ctx utils.Context, queryObj func(Sdb, *sqlx.Tx)) (err error) {

	beginx, err := s.Db[s.DefaultDb].Beginx()

	if err != nil {
		slog.Error("begin trans failed", "error", err, "requestId", ctx.RequestId)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			err = beginx.Rollback()
			slog.Error("事务回滚", "error", err, "requestId", ctx.RequestId)
			if err != nil {
				return
			}
		} else if err != nil {
			err = beginx.Rollback()
			slog.Error("事务回滚", "error", err, "requestId", ctx.RequestId)
			if err != nil {
				return
			}
		} else {
			err = beginx.Commit()
			if err != nil {
				slog.Error("提交失败", "error", err, "requestId", ctx.RequestId)
				return
			}
		}
	}()
	queryObj(s, beginx)
	return
}
