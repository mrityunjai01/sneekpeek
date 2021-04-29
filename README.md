# SneekPeek
## How to Use
### Compile using go
1. Need to install go compiler
2. cd server && go build .
3. cd client && go build .
4. you have the built executables, now run server.exe on the computer you want to sneek into, make sure that you give server.exe access through the firewall if you're on a public network
5. run the server executable
6. on your machine, run client -connect <ip of the other machine on the same network>


## How does it work?
Look at the code!

there are two directories client ( containing source files for your computer)
and server (for the other computer)
