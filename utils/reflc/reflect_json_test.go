package reflc

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	Name string
	Age  int
	Sex  byte `json:"gender"`
}

type Book struct {
	ISBN     string `json:"isbn"`
	Name     string
	Price    float32      `json:"price"`
	Author   *User        `json:"author"` //把指针去掉试试
	Keywords []string     `json:"kws"`
	Local    map[int]bool //TODO 暂不支持map
}

func TestMarshal(t *testing.T) {

	user := User{
		Name: "钱钟书",
		Age:  57,
		Sex:  1,
	}
	book := Book{ISBN: "4243547567",
		Name:     "围城",
		Price:    34.8,
		Author:   &user,                      //改成nil试试
		Keywords: []string{"爱情", "民国", "留学"}, //把这一行注释掉试一下，测测null
		Local:    map[int]bool{2: true, 3: false},
	}

	type args struct {
		v interface{}
	}
	arg := args{
		v: book,
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "test001", args: arg, want: []byte("{\"isbn\":\"4243547567\",\"Name\":\"围城\",\"price\":34.799999,\"author\":{\"Name\":\"钱钟书\",\"Age\":57,\"gender\":1},\"kws\":[\"爱情\",\"民国\",\"留学\"],\"Local\":{2:true,3:false}}")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.v)
			fmt.Println(string(got))
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	var test_data []byte = []byte("{\"isbn\":\"4243547567\",\"Name\":\"围城\",\"price\":34.799999,\"author\":{\"Name\":\"钱钟书\",\"Age\":57,\"gender\":1},\"kws\":[\"爱情\",\"民国\",\"留学\"],\"Local\":{2:true,3:false}}")
	var book Book
	type args struct {
		data []byte
		v    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "test001", args: args{data: test_data, v: &book}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal(tt.args.data, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(tt.args.v)
		})
	}
}
