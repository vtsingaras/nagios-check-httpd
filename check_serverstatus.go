package main

import flag "github.com/spf13/pflag"
import (
	"encoding/json"
	"github.com/olorin/nagiosplugin"
	"github.com/parnurzeal/gorequest"
	"fmt"
	"crypto/tls"
)


func main() {
	var status_url string
	var active_threshold uint
	var idle_threshold uint
	var username string
	var password string
	var insecure bool
	var hostheader string

	flag.StringVarP(&status_url, "status-url", "s","http://localhost/server-status?view=json", "mod_status handler path")
	flag.UintVarP(&active_threshold, "active-thershold", "a", 300, "Active connections threshold.")
	flag.UintVarP(&idle_threshold, "idle-thershold", "i", 30, "Idle connections thershold.")
	flag.StringVarP(&username, "username", "u","", "Basic-Auth username.")
	flag.StringVarP(&password, "password", "p","", "Basic-Auth password.")
	flag.BoolVarP(&insecure, "insecure", "k",false, "Ignore TLS errors.")
	flag.StringVarP(&hostheader, "host-header", "h", "", "Set the Host header to this string.")
	flag.Parse()

	check := nagiosplugin.NewCheck()
	defer check.Finish()
	check.AddResult(nagiosplugin.OK, "Apache httpd alive and well.")

	request := gorequest.New()
	if hostheader != "" {
		request.Set("host", hostheader)
	}
	if ( username == "" ) != ( password == "" ) {
		check.AddResult(nagiosplugin.UNKNOWN, "<username> and <password> are required together.")
		return
	}
	if username != "" {
		request.SetBasicAuth(username, password)
	}
	if insecure {
		request.TLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}

	_, body, errs := request.Get(status_url).EndBytes()
	if len(errs) > 0 {
		check.AddResult(nagiosplugin.UNKNOWN, "Unable to fetch server status!")
		return
	}
	httpd_status := new(HttpdStatus)
	err := json.Unmarshal(body, httpd_status)
	if err != nil {
		check.AddResult(nagiosplugin.UNKNOWN, "Couldn't decode server response.")
		return
	}

	if httpd_status.Connections.Active > active_threshold {
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("Active connections: %d", httpd_status.Connections.Active))
	}
	if httpd_status.Connections.Idle < idle_threshold {
		check.AddResult(nagiosplugin.CRITICAL, fmt.Sprintf("Idle connections: %d", httpd_status.Connections.Idle))
	}

	return
}