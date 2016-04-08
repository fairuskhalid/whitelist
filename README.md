# whitelist

This is a Image Whitelist Docker plugin implementation based on Authorization. Using this plugin will limit the images that can be run on the host. The plugin will look for the allowed images from the whitelist before an image can be run.

![wlplugin.png]({{site.baseurl}}/wlplugin.png)


## Fastway to try (container)
1. download plugin and server container from docker hub
  - docker pull fairus/wlserver:v1
  - docker pull fairus/wlplugin:v1

2. run the container
  - docker run -d --restart=always -p 8080:8080 fairus/wlserver:v1
  - docker run -d --restart=always -v /var/run:/var/run -v /run/docker/plugins/:/run/docker/plugins -v /etc/group:/etc/group fairus/wlplugin:v1 /wlplugin -wlhost http://192.168.56.101:8080/getlist

3. add plugin option for docker daemon (below is using systemd)
  - sudo systemctl edit --full docker.service
  - ExecStart=/usr/bin/docker daemon ..... --authorization-plugin=whitelist-plugin
  - sudo service docker restart
  
1. to update the whitelist
  - copy out the file: docker cp [wlserver]:whitelist.dat whitelist.dat
  - update whitelist.dat with image id
  - copy inn the file: docker cp whitelist.dat [wlserver]:whitelist.dat
