package sqlite3

import (
	"database/sql/driver"
	"reflect"
	"testing"

	_ "modernc.org/sqlite"
)

func TestSQLiteDriver_Open(t *testing.T) {
	type fields struct {
		Driver *SQLiteDriver
	}
	type args struct {
		dsn string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    driver.Conn
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				Driver: nil,
			},
			args: args{
				dsn: "file:testdb?cache=shared&_journal=WAL&_fk=1",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &SQLiteDriver{}
			got, err := d.Open(tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() got = %v, want %v", got, tt.want)
			}
		})
	}
}
