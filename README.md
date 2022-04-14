 
# mothDNS
#### _A simple resursive DNS server made in Go_

GoMoth is an implementation of a basic recursive DNS server without the use of any exernal libraries other than the standard libraries provided by Go. It is a passion project, implemented to learn about how DNS and DNS servers function and specially learn Golang.

## Working

GoMoth listens on a port for UPD packets. Upon recieving a valid packet it contructs a DNSPacket structure using the byte buffer received. Then we have a specific domain name and query type to lookup. Since it does not have any local files as of now to lookup from, it then sends a DNS request to a root server. 
If the query is valid the root server will respond back with a list of TLDSs (Top Level Domain servers). Then out of any of the TLDs we choose the first one whose IPaddress is available. We then send our query to one of the TLDs. It responds back with a list of name servers for the asked domain. Then we again ask the name servers for the IP corresponding to the specific domain name.

## References
- [RFC 1034](https://datatracker.ietf.org/doc/html/rfc1034)
- [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035)
- [The Go Programming Language](https://www.gopl.io/)
- [ GitHub - EmilHernvall/dnsguide ](https://github.com/EmilHernvall/dnsguide)
- [JSON format to represent DNS data](https://tools.ietf.org/id/draft-bortzmeyer-dns-json-01.html#rfc.section.3.1)
