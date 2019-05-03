package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/url"
	"strconv"

	"gopkg.in/ldap.v3"
)

var host = flag.String("host", "ldap://localhost:389", "Server URL")
var user = flag.String("user", "scott", "Bind User")
var password = flag.String("password", "", "User Password")

func main() {
	flag.Parse()

	u, err := url.Parse(*host)
	if err != nil {
		log.Fatal(err)
	}
	port, err := net.LookupPort("tcp",u.Scheme)
	if err != nil {
		port = 389
	}

	l, err := ldap.Dial("tcp", u.Hostname()+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.StartTLS(&tls.Config{ServerName:u.Hostname()})
	if err != nil {
		log.Println("Could not connect via STARTTLS")
	}

	err = l.Bind(*user, *password)
	if err != nil {
		log.Fatal(err)
	}

  // parse cn=monitor
  // parse cn=snmp,cn=monitor
  // parse cn=counters,cn=monitor
  // parse cn=monitor connection metrics
  // publish to telegraf
}
