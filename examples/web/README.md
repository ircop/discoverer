# Very simple and dump example of web-tool for unified network iteractions.

Usage: you should form and pass url with required parameters. Tool will parse your parameters, and call required profile/method.

Some examples:

URL: x.x.x.x:8765/execute?login=script&password=pass&proto=telnet&profile=dlink&ip=100.100.100.100&port=23&enable=enablepassword&method=get_platform

where x.x.x.x - server IP;

login - switch login

password - switch password

proto - telnet (should be telnet|ssh)

port - network port for remote cli

profile - equipment profile (dlink|cisco)

ip - device IP

enable - enable password, if needed

method - method to call



Same for cisco: 

![output image](https://i.imgur.com/rAn2IGD.png)


URL: URL: x.x.x.x:8765/execute?login=script&password=pass&proto=ssh&profile=cisco&ip=100.100.100.101&port=22&enable=enablepassword&method=get_platform

Changed:

proto => ssh

IP

profile => cisco

Output:

![Output](https://i.imgur.com/NDKdCMi.png)




Or find lldp neighbors on this cisco:

Changed: method => get_lldp

Output:

![Output](https://i.imgur.com/0jvSrkL.png)


And etc. See code for all methods.

