
# Interview

### Docker

1. Install [docker](https://docs.docker.com/engine/installation/). Complete [post-install](https://docs.docker.com/engine/installation/linux/linux-postinstall/)
   instructions.
2. Run ```make init & make start``` to start development server.
3. Modify `/etc/hosts` file to contain: `127.0.0.1 interview.localhost`
4. API should now be accessible through http://interview.localhost:8080/v1/hello
5. API container logs accessible through ```make logs```
