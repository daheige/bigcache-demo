# bigcache-demo
    
    bigcache demo
    version:
    github.com/allegro/bigcache/v2 v2.2.5

# wrk test

    % wrk -d 30s -c 400 -t 8 http://localhost:1336/get-data
    Running 30s test @ http://localhost:1336/get-data
    8 threads and 400 connections
    Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    46.84ms   51.40ms 500.89ms   84.03%
    Req/Sec   710.77    737.88     4.71k    91.07%
    169716 requests in 30.10s, 26.87MB read
    Socket errors: connect 157, read 101, write 0, timeout 0
    Requests/sec:   5638.43
    Transfer/sec:      0.89MB

    pprof monitor
    http://localhost:2336/debug/pprof/

    metrics
    http://localhost:2336/metrics

    runtime logger
    2021/05/18 22:39:50 exec end map[exec_time:0.076965432]
    2021/05/18 22:39:50 request uri:  /get-data
    2021/05/18 22:39:50 request before
    2021/05/18 22:39:50 exec end map[exec_time:0.113544445]
    2021/05/18 22:39:50 request uri:  /get-data
    2021/05/18 22:39:50 request uri:  /get-data
    2021/05/18 22:39:50 exec end map[exec_time:0.054807514]
    2021/05/18 22:39:50 exec end map[exec_time:0.052259714]
    2021/05/18 22:39:50 exec end map[exec_time:0.070805861]
    2021/05/18 22:39:50 request before
    2021/05/18 22:39:50 request uri:  /get-data
    2021/05/18 22:39:50 exec end map[exec_time:0.056013418]
    2021/05/18 22:39:50 exec end map[exec_time:0.07678656]
    2021/05/18 22:39:50 exec end map[exec_time:0.081971679]

    % wrk -d 60s -c 500 -t 10 http://localhost:1336/get-data
    Running 1m test @ http://localhost:1336/get-data
    10 threads and 500 connections
    Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    47.55ms   48.45ms 455.96ms   83.97%
    Req/Sec   523.51    442.71     4.27k    94.72%
    311563 requests in 1.00m, 49.32MB read
    Socket errors: connect 259, read 103, write 0, timeout 0
    Requests/sec:   5185.10
    Transfer/sec:    840.55KB

    2021/05/18 22:44:36 exec end map[exec_time:0.065436795]
    2021/05/18 22:44:36 exec end map[exec_time:0.019547955]
    2021/05/18 22:44:36 exec end map[exec_time:0.088664581]
    2021/05/18 22:44:36 exec end map[exec_time:0.021195315]
    2021/05/18 22:44:36 exec end map[exec_time:0.039279678]
    2021/05/18 22:44:36 exec end map[exec_time:0.065451753]
    2021/05/18 22:44:36 request uri:  /get-data
    2021/05/18 22:44:36 exec end map[exec_time:0.038415149]
    2021/05/18 22:44:36 exec end map[exec_time:0.040294541]
    2021/05/18 22:44:36 exec end map[exec_time:0.041150999]
    2021/05/18 22:44:36 exec end map[exec_time:0.040331136]
    2021/05/18 22:44:36 exec end map[exec_time:0.065569611]
    2021/05/18 22:44:36 request uri:  /get-data
    2021/05/18 22:44:36 exec end map[exec_time:0.028267855]
    2021/05/18 22:44:36 exec end map[exec_time:0.056533453]
    
    bigcache http???????????????????????????????????????
    ?????????qps?????????5100qps/s,tps 0.89MB bigcache???????????????????????????
    ?????????????????????log.Println??????????????????????????????????????????????????????????????????????????????????????????
    ????????????log.Println???????????????sync.Mutex lock
``` go
// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
std.Output(2, fmt.Sprintln(v...))
}

// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
mu     sync.Mutex // ensures atomic writes; protects the following fields
prefix string     // prefix on each line to identify the logger (but see Lmsgprefix)
flag   int        // properties
out    io.Writer  // destination for output
	buf    []byte     // for accumulating text to write
	}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(calldepth int, s string) error {
now := time.Now() // get this early.
var file string
var line int
l.mu.Lock()
defer l.mu.Unlock()
// ....
    return err
}
```

    ???????????????????????????????????????????????????log.Println????????????
    ?????????part1/main.go#132
    // r.Use(AccessLog, RecoverHandler)
    // ?????????????????????
	r.Use(RecoverHandler, monitor.MonitorHandler)
    ???????????????????????????
    % wrk -d 60s -c 500 -t 10 http://localhost:1336/get-data
    Running 1m test @ http://localhost:1336/get-data
    10 threads and 500 connections
    Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.31ms    0.95ms  92.12ms   90.72%
    Req/Sec    10.38k     3.15k   40.43k    61.77%
    6202423 requests in 1.00m, 0.96GB read
    Socket errors: connect 259, read 99, write 0, timeout 0
    Requests/sec: 103204.21
    Transfer/sec:     16.34MB
    
    ?????????????????????qps??????103204req/sec,tps: 16.34MB

    ??????????????????
    % wrk -t 8 -d 60 -c 400 --latency http://localhost:1336/get-data
    Running 1m test @ http://localhost:1336/get-data
    8 threads and 400 connections
    Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.59ms    1.48ms  72.65ms   92.19%
    Req/Sec    11.80k     2.42k   30.34k    65.14%
    Latency Distribution
    50%    2.42ms
    75%    2.77ms
    90%    3.31ms
    99%    6.84ms
    5637931 requests in 1.00m, 0.87GB read
    Socket errors: connect 157, read 101, write 0, timeout 0
    Requests/sec:  93817.39
    Transfer/sec:     14.76MB

# benchmark test
    
    https://github.com/allegro/bigcache-bench

# other demo
    
    https://blog.csdn.net/youshijian99/article/details/84929438
