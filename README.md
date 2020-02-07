# go-geoip2

## Overview
This project can search for IP information through Maxmind geoip2 databases.  

This project is built using Reader in [*oschwald/maxminddb-golang*](https://github.com/oschwald/maxminddb-golang).  

You can read Maxmind databases using a local file.  
Either you can read Maxmind databases using  [*Maxmind download URL*](https://dev.maxmind.com/geoip/geoipupdate/#Direct_Downloads).  
   
If you use reading databases with Maxmind download URL(only support gzip link), it is possible to update the latest databases periodically.   
It mean you will be automatically downloaded and updated to target-path in background.  
     
So you don't need to update the latest Maxmind databases manually, So very useful.

[warning] 
Maxmind download API has a daily quota of requests.  
Set to appropriate update interval.  
  
## Getting Started

## Inspiration
This projext was inspired by [*oschwald/geoip2-golang*](https://github.com/oschwald/geoip2-golang) 

## LICENSE
This project is following The MIT.
