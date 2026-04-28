package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestManager(t *testing.T) {
	tests := []struct {
		name       string
		storage    Storage
		expiration time.Duration
	}{
		{"new", &MockStorage{data: make(Store, 10)}, 5 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mngr := NewManager(tt.storage, "", tt.expiration)

			// use default with no error
			if e := mngr.Restore(""); e == nil {
				t.Errorf("Manager.Restore() did not error as expected")
			}

			// can load a session from storage
			if e := mngr.Restore("9e934ad9-cf7a-4ab9-b8aa-9e619b30badb"); e != nil {
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
	data map[string][]byte
}

func (ms *MockStorage) Remove(key string) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStorage) Load(id string) ([]byte, error) {
	switch id {
	case "9e934ad9-cf7a-4ab9-b8aa-9e619b30badb.json":
		uid, e1 := uuid.Parse("9e934ad9-cf7a-4ab9-b8aa-9e619b30badb")
		if e1 != nil {
			panic(e1.Error())
		}
		sd := &Data{
			&uid,
			time.Now().Add(time.Minute + 5), //exp.Format("2006-01-02T15:04:05Z07:00"),
			Store{"test2": []byte("54321")},
			false,
		}
		b, e := json.Marshal(sd)
		if e != nil {
			panic("error error error")
		}

		return b, nil
	}

	b, ok := ms.data[id]
	if !ok {
		panic("error error error")
	}

	return b, nil
}

func (ms *MockStorage) Save(id string, data []byte) error {
	if ms.data == nil {
		ms.data = make(Store, 10)
	}

	b, _ := json.Marshal(data)

	ms.data[id] = b

	return nil
}

func TestManager_SetSessionIDCookie(t *testing.T) {
	mkTime01, _ := time.Parse(
		"Mon, 02 Jan 2006 15:01:05 MST",
		"Sun, 02 Mar 2025 14:18:16 GMT",
	)
	tests := []struct {
		name          string
		w             http.ResponseWriter
		r             *http.Cookie
		md            *MockStorage2
		cookieCount   int
		cookiePattern string
		wantErr       bool
	}{
		{
			"id-set",
			&MockResponse{},
			nil,
			&MockStorage2{Store{}},
			1,
			"_sid_",
			true,
		},
		{
			"set-only-once",
			&MockResponse{},
			&http.Cookie{
				Name:     "_sid_",
				Value:    "10d18518-3d9b-4af8-bcd3-3823ed03ed28",
				Quoted:   false,
				Path:     "/",
				Expires:  mkTime01,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			},
			&MockStorage2{Store{}},
			0,
			"_sid_",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.md, "", time.Minute*3)
			r := httptest.NewRequest("GET", "/", nil)
			if tt.r != nil {
				r.AddCookie(tt.r)
			}

			_ = m.LoadFromCookie(r)
			m.SetCookie(tt.w, r)

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

func TestManager_LoadFromCookie(t *testing.T) {
	mkTime01, _ := time.Parse(
		"Mon, 02 Jan 2006 15:01:05 MST",
		"Sun, 02 Mar 2025 14:18:16 GMT",
	)
	tests := []struct {
		name          string
		w             http.ResponseWriter
		r             *http.Cookie
		md            *MockStorage2
		cookieCount   int
		cookiePattern string
		wantErr       bool
	}{
		{
			"id-set",
			&MockResponse{},
			nil,
			&MockStorage2{Store{}},
			1,
			"_sid_",
			true,
		},
		{
			"set-only-once",
			&MockResponse{},
			&http.Cookie{
				Name:     "_sid_",
				Value:    "10d18518-3d9b-4af8-bcd3-3823ed03ed28",
				Quoted:   false,
				Path:     "/",
				Expires:  mkTime01,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			},
			&MockStorage2{Store{}},
			1,
			"_sid_",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.md, "", time.Minute*3)
			r := httptest.NewRequest("GET", "/", nil)
			if tt.r != nil {
				r.AddCookie(tt.r)
			}
			e := m.LoadFromCookie(r)

			if e != nil != tt.wantErr {
				t.Errorf("Manager.LoadFromCookie() = %v, wantErr %v", e, tt.wantErr)
				return
			}
		})
	}
}

// TestManager_HasExpired test that session times out when it's supposed to.
func TestManager_HasExpired(t *testing.T) {
	cases := []struct {
		name     string
		duration time.Duration
		want     bool
	}{
		{
			"not-expired",
			time.Minute, // 1 minute in the future
			false,
		},
		{
			"expired",
			-1 * time.Minute, // 1 minute in the past.
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m := NewManager(&MockStorage{}, "", c.duration)
			fmt.Printf("current time %v\n", time.Now().UTC().Format(time.RFC3339))
			fmt.Printf("session time %v\n", m.Expiration().UTC().Format(time.RFC3339))
			if m.HasExpired() != c.want {
				t.Errorf("Manager.HasExpired() = %v, want %v", m.HasExpired(), true)
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

func (ms *MockStorage2) Remove(key string) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStorage2) Load(id string) ([]byte, error) {
	switch id {
	case "10d18518-3d9b-4af8-bcd3-3823ed03ed28.json":
		sid, e1 := uuid.Parse("10d18518-3d9b-4af8-bcd3-3823ed03ed28")
		if e1 != nil {
			panic(e1)
		}
		sd := &Data{
			&sid,
			time.Now().Add(time.Minute + 5), //exp.Format("2006-01-02T15:04:05Z07:00"),
			ms.data,
			false,
		}
		b, e := json.Marshal(sd)
		if e != nil {
			panic("error error error")
		}
		return b, nil
	default:
		return nil, errors.New("id not found")
	}
}

func (ms *MockStorage2) Save(id string, data []byte) error {
	if ms.data == nil {
		ms.data = make(Store, 10)
	}

	b, _ := json.Marshal(data)

	ms.data[id] = b

	return nil
}
