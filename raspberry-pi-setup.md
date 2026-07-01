# How to set up a Raspberry Pi

To make a very lean Raspberry Pi host, install the lite version of Rasbian and then SSH to your Raspberry Pi.

## Preparation (done as root user)

Remove some unnessary packages:

```bash
apt remove --purge bluez
apt autoremove --purge
```

Update the OS:

```bash
apt update
apt upgrade
apt install vim
```

Update the firmware:

```bash
rpi-update
```

Reboot.

## Set up shell etc

```bash
sed -i 's/%sudo     ALL=(ALL:ALL) ALL/%sudo ALL=(ALL:ALL) NOPASSWD:ALL/g' /etc/sudoers
sed -i 's/%sudo\tALL=(ALL:ALL) ALL/%sudo ALL=(ALL:ALL) NOPASSWD:ALL/g' /etc/sudoers
echo 'alias ll="ls -l"' >> /etc/bash.bashrc
echo 'alias grep="grep --color=auto"' >> /etc/bash.bashrc
echo 'export HISTSIZE=99999' >> /etc/bash.bashrc
echo 'export HISTFILESIZE=9999999' >> /etc/bash.bashrc
echo "PS1='\[\033[00;32m\]\u\[\033[00;37m\]@\[\033[00;34m\]\h\[\033[00m\]:\[\033[00;35m\]\w\[\033[00m\]\$ '" >> /etc/bash.bashrc
touch /etc/cloud/cloud-init.disabled
```

## Install go

```bash
wget https://go.dev/dl/go1.23.12.linux-armv6l.tar.gz
tar -C /usr/local -xzf go1.23.12.linux-armv6l.tar.gz
```

Add this to `/etc/profile`:

```bash
PATH=$PATH:/usr/local/go/bin:/root/go/bin
GOPATH=$HOME/golang
```

## Installation

Install the flipdot-clock binary:

```bash
go install github.com/FutureSharks/flipdot-clock@latest
```

Add service configuration:

```bash
echo -e "[Unit]\nDescription=Flipdot Clock Service\nAfter=network.target\n\n[Service]\nExecStart=/root/go/bin/flipdot-clock -clock -flip-display -clock-size 2 -clock-mode transition\nRestart=always\nRestartSec=5\nUser=root\n\n[Install]\nWantedBy=multi-user.target" > /etc/systemd/system/flipdot-clock.service
systemctl daemon-reload && systemctl start flipdot-clock && systemctl enable flipdot-clock
```
