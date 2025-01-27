package integration

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/mrzack99s/cocong/constants"
)

type LDAPConnectionPool struct {
	mu   sync.Mutex
	pool chan *ldap.Conn
}

func (p *LDAPConnectionPool) newConnection(l *LDAPEndpointType) (*ldap.Conn, error) {
	if l.TLSEnable {
		cer, e := tls.LoadX509KeyPair("./certs/ldap.crt", "./certs/ldap.key")
		if e != nil {
			cer, e = tls.LoadX509KeyPair(constants.CONFIG_DIR+"/certs/ldap.crt", constants.CONFIG_DIR+"/certs/ldap.key")
			if e != nil {
				return nil, e
			}
		}

		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		conn, e := ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", l.Hostname, l.Port), config)
		if e != nil {
			return nil, e
		}
		// Set keep-alive timeout
		conn.SetTimeout(30 * time.Second)
		return conn, nil

	} else {
		conn, e := ldap.Dial("tcp", fmt.Sprintf("%s:%d", l.Hostname, l.Port))
		if e != nil {
			return nil, e
		}
		// Set keep-alive timeout
		conn.SetTimeout(30 * time.Second)

		return conn, nil
	}

}

func (p *LDAPConnectionPool) GetConnection(l *LDAPEndpointType) (*ldap.Conn, error) {
	select {
	case conn := <-p.pool:
		if conn.IsClosing() {
			return p.newConnection(l)
		}
		return conn, nil
	default:
		p.mu.Lock()
		defer p.mu.Unlock()
		return p.newConnection(l)
	}
}

func (p *LDAPConnectionPool) ReturnConnection(conn *ldap.Conn) {
	select {
	case p.pool <- conn:
	default:
		conn.Close()
	}
}

func (p *LDAPConnectionPool) ClosePool() {
	close(p.pool)
	for conn := range p.pool {
		conn.Close()
	}
}

func (l *LDAPEndpointType) NewLDAPConnectionPool() error {

	poolSize := uint(10)
	if l.PoolSize > 0 {
		poolSize = l.PoolSize
	}

	pool := make(chan *ldap.Conn, poolSize)
	for i := 0; i < int(poolSize); i++ {
		conn, err := l.pool.newConnection(l)
		if err != nil {
			return fmt.Errorf("failed to create connection: %w", err)
		}
		pool <- conn
	}

	l.pool = &LDAPConnectionPool{
		pool: pool,
	}

	return nil
}

// func (l *LDAPEndpointType) connect() (conn, err error) {
// 	if l.TLSEnable {
// 		cer, e := tls.LoadX509KeyPair("./certs/ldap.crt", "./certs/ldap.key")
// 		if e != nil {
// 			cer, e = tls.LoadX509KeyPair(constants.CONFIG_DIR+"/certs/ldap.crt", constants.CONFIG_DIR+"/certs/ldap.key")
// 			if e != nil {
// 				err = e
// 				return
// 			}
// 		}

// 		config := &tls.Config{Certificates: []tls.Certificate{cer}}

// 		l.instance, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", l.Hostname, l.Port), config)
// 		if err != nil {
// 			return
// 		}

// 	} else {
// 		l.instance, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", l.Hostname, l.Port))
// 		if err != nil {
// 			return
// 		}
// 		defer l.instance.Close()
// 	}

// 	return
// }

func (l *LDAPEndpointType) Authentication(username, password string) (err error) {

	conn, err := l.pool.GetConnection(l)
	if err != nil {
		return err
	}
	defer l.pool.ReturnConnection(conn)

	splitString := strings.Split(username, "@")
	if len(splitString) <= 1 {
		err = errors.New("invalid username")
		return
	}
	if splitString[1] != l.Domain {
		err = fmt.Errorf("username %s is not authorised by the domain name", username)
		return
	}

	err = conn.Bind(username, password)
	if err != nil {
		err = fmt.Errorf("%s credentials are invalid", username)
		return
	}

	return
}
