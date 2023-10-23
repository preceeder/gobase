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
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/jmoiron/sqlx"
	"github.com/preceeder/gobase/utils"
	"reflect"
	"testing"
)

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
			if got := s.Execute(utils.Context{}, tt.args.sqlStr, tt.args.params, tt.args.tx...); !reflect.DeepEqual(got, tt.want) {
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
			s.Fetch(utils.Context{}, tt.args.sqlStr, tt.args.params, tt.args.row)
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
			if got := s.Insert(utils.Context{}, tt.args.params, tt.args.tx...); !reflect.DeepEqual(got, tt.want) {
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
			if got := MysqlDb.InsertMany(utils.Context{}, tt.args.params, tt.args.tx...); !reflect.DeepEqual(got, tt.want) {
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
			if got := s.InsertUpdate(utils.Context{}, tt.args.params, tt.args.tx...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSdb_Select(t *testing.T) {

	InitMysqlWithStruct(map[string]MysqlConfig{
		"default": MysqlConfig{
			Host:        "127.0.0.1",
			Port:        "13306",
			Password:    "pinky@1111",
			User:        "pinky",
			Db:          "pinky",
			MaxOpenCons: 20,
			MaxIdleCons: 5,
		},
	})
	type Users struct {
		UserIds           string     `db:"userId"`
		Nicks             NullString `db:"nick"`
		LastLoginDateTime NullString `db:"lastLoginDateTime"`
	}
	ser := Users{}

	type args struct {
		sqlStr string
		params map[string]any
		row    any
	}

	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{name: "", args: args{
			sqlStr: "select nick, userId,lastLoginDateTime from t_user where id = 1018",
			params: nil,
			row:    &ser,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MysqlDb.Select(utils.Context{}, tt.args.sqlStr, tt.args.params, tt.args.row)
			fmt.Printf("%#v\n", tt.args.row)
			fmt.Println(sonic.MarshalString(tt.args.row))
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
			if err := s.Transaction(utils.Context{}, tt.args.queryObj); (err != nil) != tt.wantErr {
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
			if got := s.Update(utils.Context{}, tt.args.params, tt.args.tx...); got != tt.want {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
