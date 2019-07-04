package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ldap.v3"
)

var host = flag.String("host", "ldap://localhost:389", "Server URL")
var user = flag.String("user", "scott", "Bind User")
var password = flag.String("password", "", "User Password")
var insecure = flag.Bool("insecure", false, "Skip verify for TLS")

func main() {
	flag.Parse()

	loadDefaultConfig()

	u, err := url.Parse(*host)
	if err != nil {
		log.Fatal(err)
	}
	port, err := net.LookupPort("tcp", u.Scheme)
	if err != nil {
		port = 389
	}

	var conn *ldap.Conn
	var tlsconfig = &tls.Config{
		InsecureSkipVerify: *insecure,
		ServerName: u.Hostname(),
	}

	if u.Scheme == "ldaps" {
		conn, err := ldap.DialTLS("tcp", u.Hostname()+":"+strconv.Itoa(port), tlsconfig)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
	} else {
		conn, err := ldap.Dial("tcp", u.Hostname()+":"+strconv.Itoa(port))
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		err = conn.StartTLS(tlsconfig)
		if err != nil {
			log.Println("Could not connect via STARTTLS")
		}
	}

	err = conn.Bind(*user, *password)
	if err != nil {
		log.Fatal(err)
	}
	searchRequest := ldap.NewSearchRequest(
		"cn=Monitor",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=top)",
		[]string{},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	values := make(map[string]int)

	for _, e := range sr.Entries {
		for _, a := range e.Attributes {
			if n, e := strconv.Atoi(a.Values[0]); e == nil {
				values[a.Name] = n
			}
		}
	}

	hostname, err := os.Hostname()

	tags := []string{
		"dirsrv",
		"server=" + u.Hostname(),
		"port=" + strconv.Itoa(port),
		"host=" + hostname,
	}

	fmt.Print(strings.Join(tags, ",") + " metrics=" + strconv.Itoa(len(values)))

	for v, n := range values {
		fmt.Print("," + v + "=" + strconv.Itoa(n) + "i")
	}

	fmt.Println(" " + strconv.FormatInt(time.Now().UnixNano(), 10))

	// TODO: parse cn=monitor connection metrics
}
