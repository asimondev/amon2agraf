# AMON2AGRAF
Load AMON reports into MySQL database.

This *README* describes the usage of AMON2AGRAF tool to import AMON JSON reports 
into the MySQL database. The AMON JSON reports are produced by AMON with 
**-t baseline**, **--output-format json**, **-output-file Filename** options.
Usually you would additionally use the following AMON options **-i interval**, 
**--start-time YYYY-MM-DD HH:MM[:SS]**, **--duration Minutes**. The JSON 
report should be copied to the server with MySQL database and Grafana. On 
this server the amon2agraf is used to load the database into the MySQL table. 
Grafana uses this MySQL database to display the data using AMON dashboards.


# How to use AMON2AGRAF? #

1. Create AMON JSON report using AMON.
1. Transfer the report to the server running your MySQL server.
1. Import the data into the MySQL server using amon2agraf.
1. Start Grafana and use the AMON dashboards.


## Data export examples. ##

Use -h option to see the help text for AMON.

Create a baseline JSON report for the next 60 minutes using 1 minute interval.

    amon -t baseline -i 60 --output-format json --output-file test_baseline.json 
    -- duration 60
    

Create a baseline JSON report for a RAC database from 10am till 4pm.
    amon -tr baseline -i 60 --output-format json --output-file test_baseline.json 
    --start-time "2020-08-16 10:00" -- duration 360

Sometimes you have canceled the AMON report with Ctrl-C or just don't want to wait 
for the whole duration time. In this case you can copy the existing JSON report 
file and add **]}** an the end. This would make this JSON object valid. After that 
the report can be reported as usual.


## Data import examples. ##

Use -h option to see the help text for amon2agraf.

Import data for into the MySQL database grafana using the specified database
account. The location of the AMON report is provided by the "-file" option.

    amon2agraf -database grafana -user grafana -password GrafanaPassword -file amon_baseline.json.

Usually you would use the same database account for all AMON imports. In this 
case you should create a JSON config file with your database account and only 
specify the report file name from the command line. For example, you could 
create the JSON file agraf.env:
   {
      "database": "grafana",
      "user": "grafana",
      "password": "MyGrafana1",
   }

You can now import the amon_baseline.json file using the following options:

    amon2agraf -config agraf.env -file amon_baseline.json

