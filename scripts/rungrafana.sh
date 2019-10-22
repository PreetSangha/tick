docker run --rm -itd -p 3000:3000 --name=grafana -v /c/temp/influx/grafana/:/var/lib/grafana -e GF_SECURITY_ADMIN_PASSWORD=admin grafana/grafana
