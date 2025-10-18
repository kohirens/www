package validation

import "testing"

func TestEmail(runner *testing.T) {
	cases := []struct {
		name  string
		email string
		want  string
		want1 bool
	}{
		{"jons-email-is-good", "jons@example.com ", "jons@example.com", true},
		{"printable-chars-are-ok", " jons{!#$%&'*+-/=?^_|~}@example.com ", "jons{!#$%&'*+-/=?^_|~}@example.com", true},
		{
			"domain-label-too-long-not-ok",
			" jons{!#$%&'*+-/=?^_|~}@subdomainwaytoolongtobeavalidemailaddressbecausethemaxlenthalabelcanbeis63characters.example.com ",
			"jons{!#$%&'*+-/=?^_|~}@subdomainwaytoolongtobeavalidemailaddressbecausethemaxlenthalabelcanbeis63characters.example.com",
			false,
		},
		{
			"missing-local-part-not-ok",
			"@example.com",
			"@example.com",
			false,
		},
	}

	for _, c := range cases {
		runner.Run(c.name, func(t *testing.T) {
			got, got1 := Email(c.email)
			if got != c.want {
				t.Errorf("Email() got = %v, want %v", got, c.want)
			}
			if got1 != c.want1 {
				t.Errorf("Email() got1 = %v, want %v", got1, c.want1)
			}
		})
	}
}

func TestMaxLen(t *testing.T) {
	cases := []struct {
		name    string
		subject string
		max     int
		want    bool
	}{
		{"less-than-max", "yes", 5, true},
		{"equal-max", "no", 2, true},
		{"exceed-max", "way-over", 5, false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxLen(tt.subject, tt.max); got != tt.want {
				t.Errorf("MaxLen() = %v, want %v", got, tt.want)
			}
		})
	}
}
