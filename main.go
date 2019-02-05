package main

import (
    "os"
    "os/signal"
    "fmt"
    "syscall"
    "time"
    "context"

    flags "github.com/jessevdk/go-flags"
    log "github.com/sirupsen/logrus"

    "database/sql"
    _ "github.com/lib/pq"
)

var version string = "undef"

type Options struct {
    Debug   bool   `env:"DEBUG"    long:"debug"    description:"enable debug"`
    LogFile string `env:"LOG_FILE" long:"log-file" description:"path to JSON log file"`

    Conn    string `               long:"connect" description:"pgsql connection string" required:"true"`
    Query   string `               long:"query"   description:"pgsql query string"        default:"SELECT current_time(1)"`
    Period  time.Duration `        long:"period"  description:"delay between invocations" default:"1s"`
}

func main() {
    var opts Options

    _, err := flags.Parse(&opts)
    if err != nil {
        os.Exit(1)
    }

    if opts.Debug {
        log.SetLevel(log.DebugLevel)
    }

    if opts.LogFile != "" {
        logFp, err := os.OpenFile(opts.LogFile, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0600)
        checkError(fmt.Sprintf("error opening %s", opts.LogFile), err)

        defer logFp.Close()

        // ensure panic output goes to log file
        syscall.Dup2(int(logFp.Fd()), 1)
        syscall.Dup2(int(logFp.Fd()), 2)

        // log as JSON
        log.SetFormatter(&log.JSONFormatter{})

        // send output to file
        log.SetOutput(logFp)
    }

    log.Debug("hi there! (tickertape tickertape)")
    log.Infof("version: %s", version)

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt)

    // does not actually connect to database
    db, err := sql.Open("myPgDialer", opts.Conn)
    if err != nil {
        log.Fatalf("unable to create connection to database: %v", err)
    }
    defer db.Close()

    keepGoing := true
    for keepGoing {

        select {
            case <- sigCh:
                keepGoing = false
                log.Debug("got signal")

            case <- time.After(opts.Period):
                log.Debug("tick")

                ctx, cancel := context.WithTimeout(context.Background(), opts.Period)

                if err = db.PingContext(ctx); err != nil {
                    log.Errorf("unable to ping database: %v", err)
                } else {
                    log.Debug("connection good")

                    var result interface{}
                    if err = db.QueryRowContext(ctx, opts.Query).Scan(&result); err != nil {
                        log.Errorf("unable to execute query: %v", err)
                    } else {
                        log.Debug("query result: %#v", result)
                    }
                }

                cancel()
        }

    }

    log.Info("done")
}
