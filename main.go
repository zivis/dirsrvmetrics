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
  "encoding/json"

	"gopkg.in/ldap.v3"
)

var conffile = flag.String("config","","path to config-file(json)")
var host = flag.String("host", "ldap://localhost:389", "Server URL")
var user = flag.String("user", "scott", "Bind User")
var password = flag.String("password", "", "User Password")
var base = flag.String("base", "cn=Monitor", "Base for metrics")

type Config struct {
  Host     string `json:"host"`
  User     string `json:"user"`
  Password string `json:"password"`
  Base     string `json:"base"`
}

func LoadConfiguration(file string) Config {
    var config Config
    configFile, err := os.Open(file)
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)
    return config
}

func main() {
	flag.Parse()

  if *conffile != "" {
    config := LoadConfiguration(*conffile)
    if config.Host != "" {
      host = &config.Host
    }
    if config.User != "" {
      user = &config.User
    }
    if config.Password != "" {
      password = &config.Password
    }
  }

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
		for _, a := range e.Attributes {
			if n, e := strconv.Atoi(a.Values[0]); e == nil {
				values[a.Name] = n
			}
		}
	}

	hostname, err := os.Hostname()


	tags := []string{
		"dirsrv",
		"server="+u.Hostname(),
		"port="+strconv.Itoa(port),
		"host="+hostname,
	}

	fmt.Print(strings.Join(tags, ",") + " metrics="+strconv.Itoa(len(values)))

	for v, n := range values {
		fmt.Print(","+v+"="+strconv.Itoa(n)+"i")
	}

	fmt.Println(" "+strconv.FormatInt(time.Now().UnixNano(), 10))

  // TODO: parse cn=monitor connection metrics
}
