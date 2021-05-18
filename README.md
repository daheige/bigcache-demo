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
    
    bigcache http服务，基本上qps单机在5100qps/s,tps 0.89MB bigcache效率还是相当不错的

# benchmark test
    
    https://github.com/allegro/bigcache-bench

# other demo
    
    https://blog.csdn.net/youshijian99/article/details/84929438
