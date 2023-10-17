//   File Name:  method_test.go.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/10/17 17:08
//    Change Activity:

package reflc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDUnmarshal(t *testing.T) {
	type args struct {
		p     reflect.Type
		value string
	}

	var def []string = []string{}
	p := reflect.ValueOf(def)
	tests := []struct {
		name    string
		args    args
		want    *reflect.Value
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "", args: args{p: p.Type(), value: "[34,34,45]"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DUnmarshal(tt.args.p, tt.args.value)
			fmt.Println(got.Interface())
			if (err != nil) != tt.wantErr {
				t.Errorf("DUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DUnmarshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
