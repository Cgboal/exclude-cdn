#### Exclude CDN

This tool wraps ProjectDiscovery's CDNCheck library to allow users to filter out CDN hosts from a list consisting of IP's, URL's, and Domains passed via stdin. 

```bash
λ cgboal [~] → crobat -s hackerone.com | httpx -silent | exclude-cdn
https://docs.hackerone.com
https://mta-sts.forwarding.hackerone.com
https://mta-sts.managed.hackerone.com
https://mta-sts.hackerone.com

λ cgboal [~] → crobat -s hackerone.com | httpx -silent              
https://mta-sts.hackerone.com
https://mta-sts.forwarding.hackerone.com
https://mta-sts.managed.hackerone.com
https://docs.hackerone.com
https://hackerone.com
https://api.hackerone.com
https://www.hackerone.com
https://api.hackerone.com
https://www.hackerone.com
https://support.hackerone.com
https://hackerone.com
```
