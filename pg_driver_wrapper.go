package main

import (
    "net"
    "time"

    "database/sql"
    "database/sql/driver"
    "github.com/lib/pq"

    log "github.com/sirupsen/logrus"
)

type pgDriverWrapper struct{}

func (d pgDriverWrapper) Open(name string) (driver.Conn, error) {
    log.Debugf("Open(%s)", name)
    return pq.DialOpen(pgDriverWrapper{}, name)
}

func (d pgDriverWrapper) Dial(network, address string) (net.Conn, error) {
    log.Debugf("Dial(%s, %s)", network, address)
    conn, err := net.Dial(network, address)
    if err == nil {
        log.Infof("connected to %s", conn.RemoteAddr())
    }

    return conn, err
}

func (d pgDriverWrapper) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
    log.Debugf("DialTimeout(%s, %s, %v)", network, address, timeout)
    conn, err := net.DialTimeout(network, address, timeout)
    if err == nil {
        log.Infof("connected to %s", conn.RemoteAddr())
    }

    return conn, err
}

func init() {
    sql.Register("myPgDialer", pgDriverWrapper{})
}
