package geoip_location

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

// Init initializes the plugin
func init() {
	caddy.RegisterPlugin("geoip", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	config, err := parseConfig(c)
	if err != nil {
		return err
	}

	dbhandler, err := maxminddb.Open(config.DatabasePath)
	if err != nil {
		return c.Err("geoip: Can't open database: " + config.DatabasePath)
	}
	// Create new middleware
	newMiddleWare := func(next httpserver.Handler) httpserver.Handler {
		return &GeoIP{
			Next:      next,
			DBHandler: dbhandler,
			Config:    config,
		}
	}
	// Add middleware
	cfg := httpserver.GetConfig(c)
	cfg.AddMiddleware(newMiddleWare)

	return nil
}

func (gip GeoIP) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	gip.lookupLocation(w, r)
	return gip.Next.ServeHTTP(w, r)
}

func (gip GeoIP) lookupLocation(w http.ResponseWriter, r *http.Request) {
	record := gip.fetchGeoipData(r)

	replacer := httpserver.NewReplacer(r, nil, "")
	replacer.Set("geoip_latitude", strconv.FormatFloat(record.Location.Latitude, 'f', 6, 64))
	replacer.Set("geoip_longitude", strconv.FormatFloat(record.Location.Longitude, 'f', 6, 64))

	if rr, ok := w.(*httpserver.ResponseRecorder); ok {
		rr.Replacer = replacer
	}
}

func (gip GeoIP) fetchGeoipData(r *http.Request) GeoIPRecord {
	clientIP, _ := getClientIP(r, true)



	return record
}

func getClientIP(r *http.Request, strict bool) (net.IP, error) {
	var ip string

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, errors.New("unable to parse address")
	}

	return parsedIP, nil
}

