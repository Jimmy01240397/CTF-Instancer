# CTF-Instancer

## Run
1. Edit `docker-compose.yml`
```yaml
volumes:
- ./chal:/app/chal:ro
environment:
- PORT=8000
- SESSIONNAME=session
# Service mode web or api
- SERVICEMODE=web
# Token in ctfd api mode
- TOKEN=testtoken
# Your Instancer Title
- TITLE=
# Instance port range
- MINPORT=30000
- MAXPORT=31000
# Instance Validity
- VALIDITY=3m
# Dynamic flag in ctfd api mode
- FLAGPREFIX=TSC
- FLAGMSG=testflag
# Instance subnet prefix
- PREFIX=29
# Instance subnet pool
- SUBNETPOOL=10.200.0.0/16
# Challenge Dir
- CHALDIR=chal
- BASESCHEME=http
# Base host name. For example use aaa.com you will get <id>.aaa.com for instance host
- BASEHOST=
- CAPTCHA_SITE_KEY=
- CAPTCHA_SECRET_KEY=
# CTFD URL
- CTFDURL=
# Multiport port: port id. for example: MODE{ID}=Proxy MODE{ID}=Forward MODE{ID}=Command
- MODE0=Proxy
# Command template in command mode. for example: COMMAND{ID}=nc {{ .BaseHost }} {{ .Port }}
- COMMAND0=nc {{ .BaseHost }} {{ .Port }}
ports:
# Same as PORT environment
- 8000:8000
```

2. Move your challenge to `CHALDIR`

3. Challenge docker-compose.yml example
```yaml
version: '3'
services:
  chal:
    image: chal
    build: .
    ports:
    # Instancer will use ${PORT} to control your port
    - ${PORT0}:11111
    environment:
    - FLAG=${FLAG}
    volumes:
    - /tmp/${ID}/userid:/userid:ro
    - /tmp/${ID}/flag:/flag:ro
    networks:
      default:

networks:
  default:
    ipam:
      config:
      # Instancer will use ${SUBNET<Number>} to control your subnet
      - subnet: ${SUBNET0}
```

4. Run `docker compose up -d` and wait for 2 minute
