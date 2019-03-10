# bank

version changes:

- v0.5 - original stable
- v0.6 - with close connection
- v0.7 - added transport
- v0.8 - modified client to be unique, added viper for config


Using netstat -n | grep -i 9596 | grep -i time_wait | wc -l   to check iddle conns

Docs read:

http://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/

https://www.reddit.com/r/golang/comments/az1gtq/how_to_calculate_resources_sizing_to_prevent_a/

https://stackoverflow.com/questions/10184975/ab-apache-bench-error-apr-poll-the-timeout-specified-has-expired-70007-on

https://golang.org/pkg/net/http/#Transport

To test:

ab -c 100 -n 50000 -s 120 -k -m GET http://192.168.1.3:9296/echo/javier