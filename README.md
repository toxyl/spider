# Spider
Spider, the master of deception and destruction, builds webs and catches prey. But in this case, the prey are pesky bots, and the web is a trap set on specified ports.

## Spiders vs. Prey
Spider wastes the time and resources of bots. Using `spiderctl`, you deploy it to multiple hosts (one Spider per host). Each Spider responds to connection attempts to the ports it protects with banners that look like real services. After a brief pause, it starts to send chunks of random garbage (data that can contain any (non)-printable characters) until a configurable duration has been exceeded, and then, to add insult to injury, it closes the connection with a random taunt in plain text. Depending on the client's implementation, this might cause errors or even crashes, thus interrupting the bot's operation. 

## WARNING! 
Spider can interfere with legitimate software that interfaces with the specified ports, such as company vulnerability scanners. Always obtain permission from the admin before using Spider on a network. And make sure to carefully curate the whitelist to avoid trouble.

## How To Use
Spider is designed to be deployed and managed as a cluster of independent spiders distributed across multiple hosts. The main application for controlling a cluster is `spiderctl` which is used together with an inventory YAML file for each cluster you want to manage. `spiderctl` takes care of deployment (configuring, building, uploading, setting up services) and management (starting, stopping, resetting, gathering stats) of the given inventory.

In this loose cluster each Spider acts independently and has no knowledge of any other Spider. The only things they have in common are their configuration and the metrics server they use to report stats to.

To set up a spider cluster, follow these steps:
1. Build `spiderctl` if you haven't done so already.
2. Create an inventory YAML (e.g. `inventory.yaml`, can be any name you want) for the cluster.
3. Run `spiderctl inventory.yaml deploy` to deploy the cluster.
4. (optional) Use the Prometheus endpoint to scrape cluster stats, e.g. for a Grafana dashboard.

### Build `spiderctl`
```bash
./build-spiderctl
```

### Create An Inventory
Create a new yaml file (e.g. `inventory.yaml`) with the following content and adjust settings as needed.
```yaml
# This section configures the MetricNexus service used for stats.
# Set `host` and `ip` to the machine you want to install MetricNexus to and
# set `key` to whatever you would like to use as API key. 
# You also have to set the `source` setting to match the location of the
# nexus-server sources (this is in `app/nexus-server/` within the repo root).
# See https://github.com/toxyl/metric-nexus for more details on MetricNexus.
metric_nexus:
  source: /path/to/spider/app/nexus-server/
  destination: /usr/local/bin/
  host: my-host-01
  ip: 192.168.1.2
  port: 4096
  key: MyTotallyNotSafeAPIKey
  state_file: /etc/metric-nexus/state.yml
  ssl_cert_file: # leave empty to create a self-signed certificate
  ssl_key_file: # leave empty to create a self-signed certificate

# This section configures the spiders that will be installed on all hosts. 
# Set the `source` setting to match the location of the
# spider sources (this is in `app/spider/` within the repo root).
# Set `spiders` to all ports you want to protect. If there is already a service
# running on the given port, Spider will just skip it.
# Don't forget to setup the `whitelist`, especially if there are devices on the network
# that have legitimate reasons to connect to the hosts, such as vulnerability scanners. 
spider:
  source: /path/to/spider/app/spider/
  destination: /usr/local/bin/

  # Spider matches the specified ports to service names and banners, 
  # or creates random banners based on the list of taunts. 
  # If a port is already in use, Spider backs off and continues with the rest of the list. 
  # If `ufw` or `iptables` is present on the host, `spiderctl` will automatically open the ports, so your spiders won't starve.
  spiders:
    - 23 # Telnet
    - 80 # HTTP
    - 443 # HTTPS
    - 445 # SMB
    - 3389 # RDP
    - 5900 # VNC

  # Spider immediately closes the connection of any IP on the whitelist, 
  # instead of attacking them with random data. 
  # You should add the IPs of company virus scanners, VPNs, and similar devices to the whitelist to avoid trouble.
  whitelist:
    - 192.168.1.1
 
  # These taunts are the final messages that Spider can send before 
  # terminating the connection. It choses randomly which taunt to send.
  taunts:
    - "BAM! You're dead!"
    - "Welcome DEADBEEF!"
    - "Gotcha!"
    - "Your IP has been reported and you are being watched."
    - "Next time I'll retaliate!"

  # During this time (given in seconds), Spider sends random-length chunks of 
  # garbage data interleaved with random-length periods of inactivity, 
  # just like a real spider slowly devouring its prey.
  attack_length: 1200

# This section defines which credentials to use for SSH connections to the hosts.
# Note that hosts must be configured in your SSH config (`~/.ssh/config`), 
# these credentials will only overrule the corresponding parts of the config,
# others are used as configured in the SSH config. 
credentials:
  user: my_user
  group: my_user # used for file permissions
  key: ~/.ssh/my_user

# The last section defines the list of hosts to deploy Spider to.
# The names must match entries in your SSH config at `~/.ssh/config`.
hosts:
- my-host-01
- my-host-02
- my-host-03
- my-host-04
- my-host-05
```
 
### Using `spiderctl`
```
spiderctl [inventory] [playbook] <optional args>
```
Where `inventory` is a YAML file like the one created in the last step. See below for a list of playbooks.

Like Ansible `spiderctl` will remove hosts from processing if an error is encountered and continue with the hosts left. It operates in parallel, i.e. it will run each task of a play on all hosts simultaneously but it will wait for all hosts to finish before proceeding to the next task. 

| Playbook | Description | Example |
| --- | --- | --- |
| `deploy` | This will install/update the MetricNexus server and all Spiders. | ![spiderctl deploy](spiderctl_deploy.png?raw=true "spiderctl deploy") |
| `stats` | Want to know how the spiders are performing? Run this. | ![spiderctl stats](spiderctl_stats.png?raw=true "spiderctl stats") |
| `reset` | Stats of the cluster are permanent, i.e. you won't lose them if MetricNexus or a Spider restarts. However, sometimes you might want to reset all stats (prey, kills, uptime, etc.). This is what you need then. |![spiderctl reset](spiderctl_reset.png?raw=true "spiderctl reset") |
| `start` | Stops/pauses your cluster. | ![spiderctl start](spiderctl_start.png?raw=true "spiderctl start") |
| `stop` | Starts/resumes your cluster. | ![spiderctl stop](spiderctl_stop.png?raw=true "spiderctl stop")|
| `exec <args>` | Need to perform operations on your cluster that are not available as `spiderctl` command? Using this you can execute shell commands on all hosts. Unlike other playbooks the hosts are processed in a serial fashion to avoid mix-ups with interactive SSH sessions. | ![spiderctl exec](spiderctl_exec.png?raw=true "spiderctl exec") |

### Metrics
The MetricNexus server exposes a Prometheus endpoint that you can use to, for example, make Grafana dashboards for your cluster(s). Please note that this, on purpose, does not collect stats on the level of individual hosts but rather aggregates stats of all hosts. It's also important to note that the Prometheus config needs to take into account that an API key is required. Here's an example job config:
```yaml
  - job_name: 'spiders'
    scrape_interval: 10s
    metrics_path: /__metrics
    scheme: https # note that this is https, not http!
    tls_config:
      insecure_skip_verify: true # must be enabled to allow the self-signed certificate of the MetricNexus server
    authorization:
      type: "token"
      credentials: "MyTotallyNotSafeAPIKey" # must match the MetricNexus key configured in your inventory 
    static_configs:
    - targets:
      - 192.168.1.2:4096 # must match the MetricNexus IP and port configured in your inventory
```

The repo includes an example [Grafana dashboard](dashboard.json):

![Grafana Dashboard](dashboard.png?raw=true "Grafana Dashboard")

## Sidenotes For The Technically Minded
Under the hood this tool makes use of templating to generate configuration that is build into the binaries, so they can be kept small in filesize and memory footprint (on my test machines an instance of Spider takes about 16 MB of RAM, and the MetricNexus server takes about 35 MB). For that reason `app/nexus-server/config.go` and `app/spider/config.go` are excluded via the `.gitignore` as they will change with every deployment and are specific to your configuration. Furthermore, the tool uses Go to build the binaries, so you need a working Go 1.20 (or newer) installation.  
Also worth to mention: to upload files to the hosts the tool uses `scp` and for commands `ssh`. This does require that you define hosts using the names from your SSH config (`~/.ssh/config`). But, the tool will use the user name and SSH key defined in your inventory, overruling the SSH config. Other settings are not overruled. 

## Disclaimer
Please note that Spider is provided "as is" and without warranty of any kind, either express or implied, including, but not limited to, the implied warranties of merchantability and fitness for a particular purpose. The author(s) of Spider shall not be held liable for any damages arising from the use of this software, including, but not limited to, any direct, indirect, special, incidental, or consequential damages arising out of the use or inability to use this software (including, but not limited to, damages for loss of profits, business interruption, loss of data, or any other financial loss).

Please use Spider at your own risk, and be aware that there may be bugs or other issues that could cause unintended damage. While the author(s) have made every effort to ensure that the software is technically correct, they cannot guarantee its accuracy or reliability in all circumstances. If you encounter any problems while using Spider, please report them on the GitHub Issues page so that they can be addressed.