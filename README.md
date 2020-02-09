# go-geoip2
<p align="left">
<a href="https://hits.seeyoufarm.com"/><img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fgjbae1212%2Fgo-geoip2"/></a>
<a href="https://goreportcard.com/report/github.com/gjbae1212/go-geoip2"><img src="https://goreportcard.com/badge/github.com/gjbae1212/go-geoip2" alt="Go Report Card" /></a>
<a href="https://godoc.org/github.com/gjbae1212/go-geoip2"><img src="https://img.shields.io/badge/godoc-reference-5272B4"/></a>
<a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-GREEN.svg" alt="license" /></a>
</p>

## Overview
This project can search for IP information through Maxmind geoip2 databases.  

This project is built using Reader in [*oschwald/maxminddb-golang*](https://github.com/oschwald/maxminddb-golang).  

You can read Maxmind databases using a local file.  
Either you can read Maxmind databases using  [*Maxmind download URL*](https://dev.maxmind.com/geoip/geoipupdate/#Direct_Downloads).  
   
If you use reading databases with Maxmind download URL(only support gzip link), it is possible to update the latest databases periodically.   
It mean you will be automatically downloaded and updated to target-path in background.  
     
So you don't need to update the latest Maxmind databases manually, So very useful.
 
**[warning]** 
Maxmind download API has a daily quota of requests.  
Set to appropriate update interval.  
  
## Getting Started
```go
// pseudo code
package main
import (
  "net"
  "github.com/gjbae1212/go-geoip2"
  
)

func main() {
   // db, err := Open("local-file-path")
   db, err := OpenURL("maxmind license key", "GeoLite2-Country", 
      geoip2.WithUpdateInterval(6 * time.Hour), geoip2.WithRetries(2), geoip2.WithSuccessFunc(func(){}),...)
   if err != nil {
   	  panic(err)
   }
   
   ip := net.ParseIP("8.8.8.8")
   record, err := db.City(ip)
   if err != nil {
      panic(err)
   }
}
```
## Inspiration
This project was inspired by [*oschwald/geoip2-golang*](https://github.com/oschwald/geoip2-golang) 

## LICENSE
This project is following The MIT.
