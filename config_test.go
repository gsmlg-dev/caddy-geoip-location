package geoip_location

import (
	"reflect"
	"testing"

	"github.com/caddyserver/caddy"
)

func TestParseConfig(t *testing.T) {
	controller := caddy.NewTestController("http", `
		localhost:4000
		geoip path/to/maxmind/db
	`)
	actual, err := parseConfig(controller)
	if err != nil {
		t.Errorf("parseConfig return err: %v", err)
	}
	expected := Config{
		DatabasePath: "path/to/maxmind/db",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v actual %v", expected, actual)
	}
}
