# CTF-Instancer

## Run
1. Edit `docker-compose.yml`
```yaml
volumes:
- ./chal:/app/chal:ro
environment:
- PORT=8000
- SESSIONNAME=session
- DBNAME=instance.db
# Your Instancer Title
- TITLE=
# Instance port range
- MINPORT=30000
- MAXPORT=31000
# Instance Validity
- VALIDITY=3m
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
- PROXYMODE=true
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
    build: .
    ports:
    # Instancer will use ${PORT} to control your port
    - ${PORT}:11111
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
