package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"

	"gopkg.in/ldap.v3"
)

var host = flag.String("host", "ldap://localhost:389", "Server URL")
var user = flag.String("user", "scott", "Bind User")
var password = flag.String("password", "", "User Password")
var base = flag.String("base", "cn=Monitor", "Base for metrics")

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
	searchRequest := ldap.NewSearchRequest(
		*base,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=top)",
		[]string{},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	values := make(map[string]int)

	for _, e := range sr.Entries  {
		fmt.Println(e.DN)
		for _, a := range e.Attributes {
			if n, e := strconv.Atoi(a.Values[0]); e == nil {
				values[a.Name] = n
			}
		}
	}

	for v, n := range values {
		fmt.Println("  "+v+strconv.Itoa(n))
	}

  // TODO: parse cn=monitor connection metrics
  // publish to telegraf
}
