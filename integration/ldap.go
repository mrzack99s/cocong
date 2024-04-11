package integration

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/mrzack99s/cocong/constants"
)

func (l *LDAPEndpointType) connect() (err error) {
	if l.TLSEnable {
		cer, e := tls.LoadX509KeyPair("./certs/ldap.crt", "./certs/ldap.key")
		if e != nil {
			cer, e = tls.LoadX509KeyPair(constants.CONFIG_DIR+"/certs/ldap.crt", constants.CONFIG_DIR+"/certs/ldap.key")
			if e != nil {
				err = e
				return
			}
		}

		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		l.instance, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", l.Hostname, l.Port), config)
		if err != nil {
			return
		}

	} else {
		l.instance, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", l.Hostname, l.Port))
		if err != nil {
			return
		}
		defer l.instance.Close()
	}

	return
}

func (l *LDAPEndpointType) Authentication(username, password string) (err error) {

	l.connect()

	splitString := strings.Split(username, "@")
	if len(splitString) <= 1 {
		err = errors.New("invalid username")
		return
	}
	if splitString[1] != l.Domain {
		err = fmt.Errorf("username %s is not authorised by the domain name", username)
		return
	}

	err = l.instance.Bind(username, password)
	if err != nil {
		err = fmt.Errorf("%s credentials are invalid", username)
		return
	}

	l.instance.Unbind()

	return
}
