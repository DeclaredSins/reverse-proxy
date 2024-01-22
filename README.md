# TCP Reverse Proxy
A simple, transparent tcp proxy. This reverse proxy works for any application (without encryption). 

I designed this to be used with TCP stuff only, as other solutions seems to require some sort of config (which it does work but I just get high latency and unstable experience and I'm not okay with that.)
Since this is made with go, it does utilize go routines, supports hostname lookup. Although it currently only supports one host, the config is really simple. You only need to specify destination ip and port.

