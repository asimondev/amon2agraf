# Installation of AMON / Grafana VirtualBox VM. #

If you want to display AMON JSON reports in Grafana, you have to 
install MySQL database, amon2agraf tool and Grafana with AMON 
dashboards. Instead of making all this steps on yourself, you should 
install VirtualBox Linux VM with Vagrant and use the provided Docker Grafana container with AMON dashboards. This README described all 
these steps.

## Installation of VirtualBox Linux VM. ##

At first you have to download and install VirtualBox and extensions. For 
automation of Linux VM installation you have to install Vagrant tool 
from HashiCorp. You can download Vagrant using this URL: https://www.vagrantup.com/downloads.html .

Using Vagrant you have to download Oracle Linux 7 Box:  

`vagrant box add ol7u7 https://yum.oracle.com/boxes/oraclelinux/ol77/ol77.box`

The provided Vagrantfile uses the name ol7u7 to reference this Vagrant box.

In the next step you should download and unpack the latest version of amon_vagrant.tar.gz from GitHub (https://github.com/asimondev/amon2agraf/tree/master/VirtualBox) in a new directory. This directory will be the main directory for this Vagrant project. The file Vagrantfile contains all details about the new Linux VM. If you like, you can change the default values for CPU (2) and memory (4096 MB) for the new VM. You can start the process of creating the Oracle Linux VM from this directory with the command `vagrant up`. Usually this step takes some minutes. The installation is completed, if you see the following output of the script:

```
    Installation finished. Please restart this Linux VM:
    vagrant halt ol7amon1
    vagrant up
```

## Checking the timezone in Linux VM. ##

Sometimes the current timezone will be set to UTC after installation.
This is wrong, because all timestamps will be converted into the
MySQL database during import into UTC. So you have to verify the
current timezone.

```
[vagrant@ol7amon1 ~]$ timedatectl
      Local time: Tue 2020-09-01 15:13:07 UTC
  Universal time: Tue 2020-09-01 15:13:07 UTC
        RTC time: Tue 2020-09-01 15:13:08
       Time zone: UTC (UTC, +0000)
     NTP enabled: yes
NTP synchronized: yes
 RTC in local TZ: no
      DST active: n/a

[vagrant@ol7amon1 ~]$ ls -l /etc/localtime
lrwxrwxrwx. 1 root root 23 Oct  3  2019 /etc/localtime -> /usr/share/zoneinfo/UTC
```

This output of *timedatectl* shows the UTC timezone. The following steps will set the timezone to Europe/Berlin.

```
[vagrant@ol7amon1 ~]$ sudo timedatectl set-timezone Europe/Berlin

[vagrant@ol7amon1 ~]$ ls -l /etc/localtime
lrwxrwxrwx 1 root root 35 Sep  1 17:14 /etc/localtime -> ../usr/share/zoneinfo/Europe/Berlin

[vagrant@ol7amon1 ~]$ timedatectl
      Local time: Tue 2020-09-01 17:15:13 CEST
  Universal time: Tue 2020-09-01 15:15:13 UTC
        RTC time: Tue 2020-09-01 15:15:14
       Time zone: Europe/Berlin (CEST, +0200)
     NTP enabled: yes
NTP synchronized: yes
 RTC in local TZ: no
      DST active: yes
 Last DST change: DST began at
                  Sun 2020-03-29 01:59:59 CET
                  Sun 2020-03-29 03:00:00 CEST
 Next DST change: DST ends (the clock jumps one hour backwards) at
                  Sun 2020-10-25 02:59:59 CEST
                  Sun 2020-10-25 02:00:00 CET
```

## Installation of Docker containers. ##

Now you can log in as root user and start the Docker Compose. All 
configurations files are in the directory /root/agraf. The MySQL 
database is already created and the account details are saved in the 
file mysql.env as environment variables. The docker-compose.yaml file 
references this file. The following Docker Compose commands will download the Docker containers from Docker Hub and start them.

```
vagrant ssh
sudo su - root
cd agraf
docker-compose up -d
```

The following Docker containers will be created:

- Grafana (default account admin/admin01) with AMON dashboards on 
the port 10201.
- Project documentation on the port 10202.

These ports are also forwarded in the Vagrantfile.  

The command `docker-compose up -d` should be used every time to start
the containers.

## Import of AMON JSON reports. ##

The created VM already contains the tool amon2agraf. You should use 
this tool to import the data from AMON reports into the MySQL database. 
Grafana uses the same MySQL account to access the data for displaying.

The file /root/agraf/mysql.json contains the MySQL account details. So 
you can make the report file (say report01.json) available to Linux VM 
and import it using the following command:  
  
`amon2agraf -config /root/agraf/mysql.json -file /vagrant/report01.json`.

## Documentation ##

You can find more documentation in you local AMON Help Docker container 
using the URL http://localhost:10202.