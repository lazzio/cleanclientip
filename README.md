# Clean Client IP

This Treafik plugin aims to clean client IPs from `X-Forwarded-For` header.

If some IPs inside `X-Forwarded-For` header contain a port, this plugin removes it to keep only IP addresses.

## Example of X-Forwarded-For header cleanable value 

```
"10.0.0.1:1234, 10.0.0.2, 10.0.0.3:5678, 10.0.0.4"
```

## Result

The result of this plugin is :

* `X-Forwarded-For` header with only clean IP addresses

* `RemoteAddr` set with the first IP address of `X-Forwarded-For` header

* `X-Real-Ip` header set with the first IP address of `X-Forwarded-For` header