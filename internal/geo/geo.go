package geo

import (
	"fmt"
	"net/netip"

	"github.com/oschwald/geoip2-golang/v2"
)

var db *geoip2.Reader

func GetGeoLocation(ip string) *geoip2.Country {
	ipAddr, _ := netip.ParseAddr(ip)
	record, err := db.Country(ipAddr)
	if err != nil {
		fmt.Println(err)
	}
	return record
}

func InitGeoDB() error {
	var err error
	db, err = geoip2.Open("./data/GeoLite2-Country.mmdb")
	if err != nil {
		return err
	}
	return nil
}
