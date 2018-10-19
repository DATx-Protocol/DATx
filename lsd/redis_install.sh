#!/bin/bash

redis=redis-3.0.3.tar.gz

installdir=/app

softdir=/opt/soft

if [ ! -d $softdir ];

then mkdir /opt/soft

fi

echo "install software"

sleep 2

cd $softdir

wget http://download.redis.io/releases/redis-3.0.3.tar.gz

echo "install dependent environment"


echo "install redis"

sleep 2

tar zxvf $redis && cd `echo $redis | awk -F".tar.gz" '{print $1}'`

make && make install

mkdir -p /app/redis/bin/ && mkdir /app/redis/etc/

cp redis.conf /app/redis/etc/

cp /usr/local/bin/redis* /app/redis/bin/

ln -s /app/redis/bin/redis* /usr/bin/

sed -i 's/daemonize no/daemonize yes/' /app/redis/etc/redis.conf

sed -i 's=pidfile /var/run/redis.pid=pidfile /app/redis/redis.pid=' /app/redis/etc/redis.conf

echo "start redis"

/app/redis/bin/redis-server /app/redis/etc/redis.conf

echo "open firewall ports"

iptables -I INPUT -p tcp --dport 6379 -j ACCEPT

echo "set startup"

echo "/app/redis/bin/redis-server /app/redis/etc/redis.conf" >> /etc/rc.local

echo "redis install success，port:6379，installation directory: /app/redis/"