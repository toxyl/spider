# Spider
Spider, the master of deception and destruction, builds webs and catches prey. But in this case, the prey are pesky bots and the web is a trap set on specified ports.

## What It Does
In short: Spider wastes the time and resources of bots. How? By simulating a fake server on specified ports, called a "spider," that responds to connection attempts with a banner that looks like a real service. After a brief pause, Spider starts to send chunks of random data until a configurable duration has been exceeded, and then it closes the connection with a taunt in plain text.

## Spiders vs. Prey
In this game of deception and destruction, the spiders (fake servers) attack the prey (bots) by sending garbage data that can contain any (non)-printable characters. Depending on the client's implementation, this might cause errors or even crashes, thus interrupting the bot's operation. And to add insult to injury, Spider terminates the connection with a taunt from a list of customizable messages.

## WARNING! 
Spider can interfere with legitimate software that interfaces with the specified ports, such as company vulnerability scanners. Always obtain permission from the admin before using Spider on a network. And make sure to carefully curate the whitelist to avoid trouble.

## Configuration
### Host
The host to protect with spiders.
```yaml
host: 0.0.0.0
```
Usually, `0.0.0.0` is fine, but you can change it if you need Spider to protect a specific IP.

### Spiders
The list of ports to protect with spiders.
```yaml
spiders:
  - 21 # FTP
  - 23 # Telnet
  - 25 # SMTP
  - 80 # HTTP
  - 110 # POP3
  - 161 # SNMP
  - 389 # LDAP
  - 443 # HTTPS
  - 445 # SMB
  - 3306 # MySQL
  - 3389 # RDP
  - 5901 # VNC
```
The application matches the specified ports to service names and banners, or creates random banners based on the list of taunts. If a port is already in use, Spider backs off and continues with the rest of the list. Don't forget to open the ports, else your spiders will starve.

### Whitelist
The list of IPs that can connect without getting attacked.
```yaml
whitelist:
  - 192.168.0.1 # router
  - 192.168.0.2 # device 1
  - 192.168.0.3 # device 2
```
Spider immediately closes the connection of any IP on the whitelist, instead of attacking them with random data. You should add the IPs of company virus scanners, VPNs, and similar devices to the whitelist to avoid trouble.

### Taunts
The list of taunts available to the spiders.
```yaml
taunts:
  - "Spider's gotcha!"
  - "Whatcha doin here?"
```
These taunts are the final messages that Spider sends before terminating the connection. The messages are randomly chosen from this list, adding insult to injury.

### Attack Length
Defines how long Spider keeps attacking before killing its prey.
```yaml
attack_length: 10
```
During this time (given in seconds), Spider sends random-length chunks of garbage data interleaved with random-length periods of inactivity, just like a real spider slowly devouring its prey.

## Disclaimer
Please note that Spider is provided "as is" and without warranty of any kind, either express or implied, including, but not limited to, the implied warranties of merchantability and fitness for a particular purpose. The author(s) of Spider shall not be held liable for any damages arising from the use of this software, including, but not limited to, any direct, indirect, special, incidental, or consequential damages arising out of the use or inability to use this software (including, but not limited to, damages for loss of profits, business interruption, loss of data, or any other financial loss).

Please use Spider at your own risk, and be aware that there may be bugs or other issues that could cause unintended damage. While the author(s) have made every effort to ensure that the software is technically correct, they cannot guarantee its accuracy or reliability in all circumstances. If you encounter any problems while using Spider, please report them on the GitHub Issues page so that they can be addressed.