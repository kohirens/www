package backend

import (
	"reflect"
	"testing"
)

type MockFixture struct {
	Name string
}

func Test_ServiceManager_Store(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		fix     *MockFixture
		wantErr bool
	}{
		{
			name:    "can-get-struct",
			key:     "any",
			fix:     &MockFixture{"test1234"},
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotKey = NewKey[*MockFixture](c.key)
			m := NewServiceManager()
			Store(m, gotKey, c.fix)

			got, err := Retrieve(m, gotKey)

			if (err != nil) != c.wantErr {
				t.Errorf("service() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if !reflect.DeepEqual(got, c.fix) {
				t.Errorf("service() got = %v, want %v", got, c.fix)
			}
		})
	}
}

func Test_service(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		typ     *MockFixture
		fix     *MockFixture
		wantErr bool
	}{
		{
			name:    "can-get-struct",
			key:     "gp",
			fix:     &MockFixture{"test1234"},
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := NewServiceManager()
			m.Add(c.key, c.fix)

			got, err := m.Get(c.key)
			if err == nil {
				c.typ = got.(*MockFixture)
			}

			if (err != nil) != c.wantErr {
				t.Errorf("service() error = %v, wantErr %v", err, c.wantErr)
				return
			}
			if c.typ == nil {
				t.Errorf("service() got = %v, want %v", c.typ, c.fix)
				return
			}
		})
	}
}
