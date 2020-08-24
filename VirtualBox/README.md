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

In the next step you should unpack the vagrant.tar file from GitHub 
in a new directory. This directory would be the main directory for this 
Vagrant project. You can start the process of creating the Oracle 
Linux VM from this directory with the command `vagrant up`. Usually 
this step takes some minutes. At the end you have to restart the new 
VM:
```
vagrant up  <-- Start and wait some minutes to complete.

vagrant status  <-- Get the status of this VM. Default name is ol7amon1.

vagrant halt ol7amon1  <-- Stop this VM.

vagrant up  <-- Start this VM.
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
the port 10101.
- Project documentation on the port 10102.

The command `docker-compose up -d` should be used every time to start
the containers.

## Import of AMON JSON reports. ##

The created VM already contains the tool amon2agraf. You should use 
this tool to import the data from AMON report into the MySQL database. 
Grafana uses the same MySQL account to access the data for displaying.

The file /root/agraf/mysql.json contains the MySQL account details. So 
you can make the report file (say report01.json) available to Linux VM 
and import it using the following command:  
  
`amon2agraf -config /root/agraf/mysql.json -file /vagrant/report01.json`.

## Documentation ##

You can find more documentation in the AMON Help Docker container 
using the URL http://localhost:10102.