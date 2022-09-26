# HOW IT CAN BE START?

## On the Linux file system:

Open or create the file ~/.bash_profile
Add this line: export WINDOWS_HOST=$(cat /etc/resolv.conf | grep nameserver | cut -d ' ' -f 2)
This will get the Windows Host IP Address and set it as an env variable for the distro on startup.

Restart the distro: wsl --shutdown from a Windows cmd terminal and wsl to start it back. You can type env in a WSL2 terminal to make sure the env variable is there.

## On the Windows side:

1 - In the Windows firewall, add an Inbound Rule for the TCP port 5432 (the exposed postgresQL port), follow those steps:

type wf.msc in a cmd terminal to open the firewall

select Inboud rules

select Port then next

select TCP and type 5432 in the specific port input then next

select Autorize connection and go through the 2 last steps

2 - Allow incomming connections from postgresQL config:

open the file pg_hba.conf located here by default C:\Program Files\PostgreSQL\12\data\pg_hba.conf
replace the line under IPv4 local connections with the folling:
`host    all             all             0.0.0.0/0            md5`
Now every ip address is allowed to access your pg server, this is necessary because the IP address of your WSL distro changes at every startup. Please consider the downsides that could imply.