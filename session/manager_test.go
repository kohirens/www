package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	tests := []struct {
		name       string
		storage    Storage
		expiration time.Duration
	}{
		{"new", &MockStorage{data: Store{}}, 5 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mngr := NewManager(tt.storage, tt.expiration)

			// use default with no error
			if e := mngr.Restore(""); e == nil {
				t.Errorf("Manager.Restore() did not error as expected")
			}

			// can load a session from storage
			if e := mngr.Restore("abcdefg"); e != nil {
				t.Errorf("Manager.Restore() = %v", e.Error())
				return
			}

			if got := mngr.Get("test2"); !reflect.DeepEqual(got, []byte("54321")) {
				t.Errorf("Manager.Restore() = %v, want %v", got, "54321")
				return
			}

			// can set and get an item from the session
			mngr.Set("test", []byte("1245"))
			if got := mngr.Get("test"); !reflect.DeepEqual(got, []byte("1245")) {
				t.Errorf("Manager.Set() = %v, want %v", got, "1245")
				return
			}

			// can remove an item from the session
			ge1 := mngr.Remove("test")
			if ge1 != nil {
				t.Errorf("Manager.Remove() = %v, want %v", ge1, "nil")
				return
			}
			if got := mngr.Get("test"); got != nil {
				t.Errorf("Manager.Remove() = %v, want %v", got, "")
				return
			}
		})
	}
}

type MockStorage struct {
	data Store
}

func (ms *MockStorage) Load(id string) (*Data, error) {
	switch id {
	case "abcdefg":
		return &Data{
			"abcdefg",
			time.Now().Add(time.Minute + 5), //exp.Format("2006-01-02T15:04:05Z07:00"),
			Store{"test2": []byte("54321")},
		}, nil
	}

	b, ok := ms.data[id]
	if !ok {
		panic("error error error")
	}

	sd := &Data{}
	if e := json.Unmarshal(b, &sd); e != nil {
		panic("error error error")
	}

	return sd, nil
}

func (ms *MockStorage) Save(data *Data) error {
	if ms.data == nil {
		ms.data = Store{}
	}

	b, _ := json.Marshal(data)

	ms.data[data.Id] = b

	return nil
}

func TestManager_SetSessionIDCookie(t *testing.T) {
	tests := []struct {
		name          string
		w             http.ResponseWriter
		r             *http.Request
		md            *MockStorage2
		cookieCount   int
		cookiePattern string
	}{
		{
			"id-set",
			&MockResponse{},
			&http.Request{},
			&MockStorage2{Store{}},
			1,
			"_sid_",
		},
		{
			"set-only-once",
			&MockResponse{},
			&http.Request{
				Header: http.Header{
					"Cookie": []string{"_sid_=10d18518-3d9b-4af8-bcd3-3823ed03ed28; Path=/; Expires=Sun, 02 Mar 2025 14:18:16 GMT; HttpOnly; Secure; SameSite=Strict"},
				},
			},
			&MockStorage2{Store{}},
			1,
			"_sid_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.md, time.Minute*3)
			m.Load(tt.w, tt.r)

			// Do NOT use w.Header().Get it will only get the first index of the header.
			cookies := tt.w.Header()["Set-Cookie"]
			gotCount := 0
			for _, cookie := range cookies {
				gotCount += strings.Count(cookie, tt.cookiePattern)
			}
			if gotCount != tt.cookieCount {
				t.Errorf("Manager.SetSessionIDCookie() = %v times, want %v", gotCount, tt.cookieCount)
				return
			}
		})
	}
}

type MockResponse struct {
	Headers http.Header
}

func (m *MockResponse) Header() http.Header {
	if m.Headers == nil {
		m.Headers = http.Header{}
	}
	return m.Headers
}

func (m *MockResponse) Write(b []byte) (int, error) {
	return len(b), nil
}

func (m *MockResponse) WriteHeader(statusCode int) {
}

type MockStorage2 struct {
	data Store
}

func (ms *MockStorage2) Load(id string) (*Data, error) {
	switch id {
	case "10d18518-3d9b-4af8-bcd3-3823ed03ed28":
		return &Data{
			"10d18518-3d9b-4af8-bcd3-3823ed03ed28",
			time.Now().Add(time.Minute + 5), //exp.Format("2006-01-02T15:04:05Z07:00"),
			ms.data,
		}, nil
	default:
		return nil, errors.New("id not found")
	}
}

func (ms *MockStorage2) Save(data *Data) error {
	if ms.data == nil {
		ms.data = Store{}
	}

	b, _ := json.Marshal(data)

	ms.data[data.Id] = b

	return nil
}
