# How to set up a Raspberry Pi

First you need to SSH to your Raspberry Pi.

Install go:

```bash
wget https://go.dev/dl/go1.23.12.linux-armv6l.tar.gz
tar -C /usr/local -xzf go1.23.12.linux-armv6l.tar.gz
```

Add this to `/etc/profile`:

```bash
PATH=$PATH:/usr/local/go/bin:/root/go/bin
GOPATH=$HOME/golang
```

Install the flipdot-clock binary:

```bash
go install github.com/FutureSharks/flipdot-clock@latest
```

Add cron configuration:

```bash
echo "* * * * * root pgrep flipdot-clock > /dev/null || /root/go/bin/flipdot-clock" > /etc/cron.d/flipdot-clock
```

Disable some services that we don't need:

```bash
systemctl disable bluetooth.service
systemctl disable triggerhappy.service
systemctl disable dbus.service
systemctl disable avahi-daemon.service
systemctl disable ModemManager.service
systemctl disable polkit.service
```
