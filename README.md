# dirsrv metrics

Project wich collects 389 Directory Server metrics.

### Building

    $ go build .

### Usage:

    $ dirsrvmetrics -host ldap://localhost -user scott -password foo
    dirsrv,server=localhost,port=389,host=localhost metrics=44,currentconnections=19i,... 1556894373217369460

#### Useable as exec plugin for Telegraf

    [[inputs.exec]]
      commands = [".../dirsrvmetrics -host ..."]
      timeout = "5s"
      data_format = "influx"

### Configuration

You can use the command flags to specify all needed information, or:

The software will read a standard `ldap.conf`/`.ldaprc` file with an additional
allowed `BINDPW` key to store the password.
That way you won't have to specify the password on the commandline.

### Transport Level Encryption

The software will attempt to encrypt its communication.  This depends on the
URL specified.  Specifying `ldap` will attempt `STARTTLS`.  Using `ldaps` will
attempt to set up a TCP connection with TLS.

If using self-signed certificates use the `-ca` command flag, or in dire
situations the `-insecure` flag to skip host key verification.

