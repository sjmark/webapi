package httputil

import (
	"net/url"
	"testing"
)

func TestSign(t *testing.T) {
	for i, tc := range []struct {
		privateKey string
		args       url.Values
		sign       string
	}{
		{"", url.Values{}, "1B2M2Y8AsgTpgAmY7PhCfg=="},
		{"privateKey", url.Values{}, "vQyOO9mYpQCwmEWEpyRcoA=="},
		{"privateKey", url.Values{"a": {"1"}}, "PqVztHIMI4ZlU0LLD+whVg=="},
		{"privateKey", url.Values{"a": {"1"}, "b": {"2"}}, "mDbyn4jDsvl0wtMfbNMsoA=="},
		{"privateKey", url.Values{"b": {"2"}, "a": {"1"}}, "mDbyn4jDsvl0wtMfbNMsoA=="},
	} {
		sign := MD5Sign(tc.privateKey, tc.args)
		if sign != tc.sign {
			t.Errorf("%dth: sign want `%s`, got `%s`", i, tc.sign, sign)
			continue
		}
		tc.args.Add("sign", sign)
		if !MD5Verify(tc.privateKey, tc.args) {
			t.Errorf("%dth: verify fail", i)
		}
	}
}
