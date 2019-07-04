# dirsrv metrics

Project wich collects 389 Directory Server metrics.

Building

    $ go build .

Usage:

    $ dirsrvmetrics -host ldap://localhost -user scott -password foo
    dirsrv,server=localhost,port=389,host=localhost metrics=44,currentconnections=19i,... 1556894373217369460

Useable as exec plugin for Telegraf

    [[inputs.exec]]
      commands = [".../dirsrvmetrics -host ..."]
      timeout = "5s"
      data_format = "influx"

The software will read a standard `ldap.conf`/`.ldaprc` file with an additional
allowed `BINDPW` key to store the password.
That way you won't have to specify the password on the commandline.
