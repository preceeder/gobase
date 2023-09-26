/*
File Name:  mysqlDb_test.go.py
Description:
Author:      Chenghu
Date:       2023/9/1 16:21
Change Activity:
*/
package mysqlDb

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func TestInitMysqlWithStruct(t *testing.T) {
	type args struct {
		config map[string]MysqlConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitMysqlWithStruct(tt.args.config)
		})
	}
}

func TestInitMysqlWithViperConfig(t *testing.T) {
	type args struct {
		config viper.Viper
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitMysqlWithViperConfig(tt.args.config)
		})
	}
}

func TestMysqlPoolClose(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MysqlPoolClose()
		})
	}
}

func TestSdb_Execute(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		sqlStr string
		params map[string]any
		tx     []*sqlx.Tx
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   sql.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			if got := s.Execute(tt.args.sqlStr, tt.args.params, tt.args.tx...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSdb_Fetch(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		sqlStr string
		params map[string]any
		row    any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			s.Fetch(tt.args.sqlStr, tt.args.params, tt.args.row)
		})
	}
}

func TestSdb_Get(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		dest  interface{}
		query string
		args  []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			if err := s.Get(tt.args.dest, tt.args.query, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSdb_Insert(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		params map[string]any
		tx     []*sqlx.Tx
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   sql.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			if got := s.Insert(tt.args.params, tt.args.tx...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSdb_InsertMany(t *testing.T) {

	type args struct {
		params map[string]any
		tx     []*sqlx.Tx
	}
	tests := []struct {
		name string
		args args
		want sql.Result
	}{
		// TODO: Add test cases.
		{name: "ss", args: args{params: map[string]any{
			"tableName": "t_user_language",
			"Set": []map[string]any{
				{"userId": "1693281526098636",
					"language": "chinese"},
				{"userId": "1693281526098636",
					"language": "english"},
			},
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MysqlDb.InsertMany(tt.args.params, tt.args.tx...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSdb_InsertUpdate(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		params map[string]any
		tx     []*sqlx.Tx
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   sql.Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			if got := s.InsertUpdate(tt.args.params, tt.args.tx...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSdb_Select(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		sqlStr string
		params map[string]any
		row    any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			s.Select(tt.args.sqlStr, tt.args.params, tt.args.row)
		})
	}
}

func TestSdb_Transaction(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		queryObj func(Sdb, *sqlx.Tx)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			if err := s.Transaction(tt.args.queryObj); (err != nil) != tt.wantErr {
				t.Errorf("Transaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSdb_Update(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		params map[string]any
		tx     []*sqlx.Tx
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			if got := s.Update(tt.args.params, tt.args.tx...); got != tt.want {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSdb_defaultDb(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			if got := s.defaultDb(); got != tt.want {
				t.Errorf("defaultDb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSdb_getDb(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		params map[string]any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantDb string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			if gotDb := s.getDb(tt.args.params); gotDb != tt.wantDb {
				t.Errorf("getDb() = %v, want %v", gotDb, tt.wantDb)
			}
		})
	}
}

func TestSdb_getTableName(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		params map[string]any
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantTableName string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			if gotTableName := s.getTableName(tt.args.params); gotTableName != tt.wantTableName {
				t.Errorf("getTableName() = %v, want %v", gotTableName, tt.wantTableName)
			}
		})
	}
}

func TestSdb_sqlPares(t *testing.T) {
	type fields struct {
		Db        Mdb
		DefaultDb string
	}
	type args struct {
		osql   string
		params map[string]any
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantSql  string
		wantArgs []any
		wantDb   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Sdb{
				Db:        tt.fields.Db,
				DefaultDb: tt.fields.DefaultDb,
			}
			gotSql, gotArgs, gotDb := s.sqlPares(tt.args.osql, tt.args.params)
			if gotSql != tt.wantSql {
				t.Errorf("sqlPares() gotSql = %v, want %v", gotSql, tt.wantSql)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("sqlPares() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
			if gotDb != tt.wantDb {
				t.Errorf("sqlPares() gotDb = %v, want %v", gotDb, tt.wantDb)
			}
		})
	}
}

func Test_initMySQL(t *testing.T) {
	type args struct {
		config map[string]MysqlConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initMySQL(tt.args.config)
		})
	}
}

func Test_readMysqlConfig(t *testing.T) {
	type args struct {
		config viper.Viper
	}
	tests := []struct {
		name string
		args args
		want map[string]MysqlConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readMysqlConfig(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readMysqlConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
