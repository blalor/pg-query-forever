pg-query-forever is little utility I wrote when trying to measure the downtime of a database upgrade.  

```json
{"time":"2019-02-05T16:51:22Z","msg":"hi there! (tickertape tickertape)","level":"debug"}
{"time":"2019-02-05T16:51:22Z","msg":"version: 1e32e87-dirty","level":"info"}
{"time":"2019-02-05T16:51:23Z","msg":"connected to 10.0.0.121:5432","level":"info"}
{"time":"2019-02-05T16:51:23Z","msg":"reconnected after 1.040444082s","level":"info"}
…
{"time":"2019-02-05T16:53:56Z","msg":"connected to 10.0.0.121:5432","level":"info"}
{"time":"2019-02-05T16:53:56Z","msg":"reconnected after 6.032698897s","level":"info"}
…
{"time":"2019-02-05T17:02:19Z","msg":"unable to ping database: dial tcp 10.0.0.121:5432: i/o timeout","level":"error"}
…
{"time":"2019-02-05T17:02:50Z","msg":"connected to 10.0.3.191:5432","level":"info"}
{"time":"2019-02-05T17:02:50Z","msg":"reconnected after 1m11.213686395s","level":"info"}
```

## building

    make

Generates `stage/pg-query-forever`

## using

    pg-query-forever \
        --debug \
        --log-file=pg-query.log \
        --connect='postgres://user:passwd@localhost:5432/db?connect_timeout=5'

`^C` to quit. `--help` for help.
