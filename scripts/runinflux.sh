#docker run  -itd --name influx -v /c/temp/influx/influx:/var/lib/influxdb -p 8083:8083 -p 8086:8086 influxdb
#-v influxdb-volume:/var/lib/influxdb \ influxdb /init-influxdb.sh 


docker run -it -p 8086:8086 -p 8083:8083     -e INFLUXDB_HTTP_ADMIN_ENABLED=true     --rm --network net --name influxdb influxdb

docker run \
--rm \
--name influxdb \
-e INFLUXDB_DB=telegraf -e INFLUXDB_ADMIN_ENABLED=true \
-e INFLUXDB_ADMIN_USER=admin -e INFLUXDB_ADMIN_PASSWORD=admin \
-e INFLUXDB_USER=telegraf -e INFLUXDB_USER_PASSWORD=telegraf \ 
-p 8083:8083 -p 8086:8086 \
influxdb 
