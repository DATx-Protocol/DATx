#! /bin/bash

RETVAL=0
prog="supervisord"
SUPERVISORD=supervisord
PID_FILE=/var/run/supervisord.pid
CONFIG_FILE=/etc/supervisord.conf

build()
{
    mse="mse"
    cd ${mse}
    ./build.sh
    if [ $? -eq 0 ] ;then
        echo "mse build success"
    else
        echo "mse build error"
        exit
    fi

    cd ../

    lsd="lsd"

    cd ${lsd}

    ./lsd_build.sh

    if [ $? -eq 0 ] ;then
        echo "lsd build success"
    else
        echo "lsd build error"
        exit
    fi

    datxchain="datxchain"
    cd ../${datxchain}/

    echo "cd ${datxchain}"

    ./chain_build.sh
    if [ $? -eq 0 ] ;then
        echo "datxchain build success"
    else
        echo "datxchain build error"
        exit
    fi

    cd ../

   printf "\n\n\t \n"
   printf "\t_________                     ____________                   \n"
   printf "\t|        \          /\              |         \       /      \n"
   printf "\t|         \        /  \             |          \     /       \n"
   printf "\t|          \      /    \            |           \   /        \n" 
   printf "\t|          |     /      \           |            \ /         \n"
   printf "\t|          |    /________\          |             /          \n"    
   printf "\t|          /   /          \         |            / \         \n"
   printf "\t|         /   /            \        |           /   \        \n"
   printf "\t|________/   /              \       |          /     \       \n"

   printf "\\n\\tDATX has been successfully built. you can start it use datx_install.sh\n\\n"
}

install()
{
    echo -n $"install datxchain to supervisor: "
    #datxchain install
    cd ./datxchain
    ./chain_install.sh
    if [ $? -eq 0 ] ;then
        echo "datxchain install success"
    else
        echo "datxchain install failed, you must solve the problem before install it."
        exit
    fi

    cd ../

    #
    pwd=`pwd`

    #redis install
    chmod a+x ./lsd/redis_install.sh
    ./lsd/redis_install.sh

    if [ $? -eq 0 ] ;then
        echo "redis server install success"
    else
        echo "redis server install failed, you must install it first."
        exit
    fi

    #mse install
    msepath="./mse/build/mse"
    if [ ! -f "${msepath}" ];then 
        printf "the %s is not exist.\n" ${msepath}
        exit
    fi

    mkdir -p /app/mse/bin/

    cp ${msepath} /app/mse/bin/

    if [ ! -f "/app/mse/bin/mse" ];then 
        print "the /app/mse/bin/mse is not exist.\n"
        exit
    fi

    mkdir -p /app/mse/bin/node_modules
    rm -rf /app/mse/bin/node_modules/*
    cp -rf ./mse/node_modules/* /app/mse/bin/node_modules

    #lsd install
    
    lsdpath="./lsd/bin/lsd"

    if [ ! -f "${lsdpath}" ];then 
        printf "the %s is not exist.\n" ${lsdpath}
        exit
    fi

    mkdir -p /app/lsd/bin/

    cp ${lsdpath} /app/lsd/bin/

    if [ ! -f "/app/lsd/bin/lsd" ];then 
        print "the /app/lsd/bin/lsd is not exist.\n"
        exit
    fi

    #supervisor install
    pip install supervisor
    RETVAL=$?

    #setup supervisor
    echo_supervisord_conf > /etc/supervisord.conf

    #set and modify supervisor config
    mkdir -p /etc/supervisord.d/

    echo "[include]" >> /etc/supervisord.conf
    echo "files = /etc/supervisord.d/*.conf" >> /etc/supervisord.conf

    confpath="${pwd}/config.conf"

    if [ ! -f "${confpath}" ];then 
        printf "the %s is not exist.\n" ${confpath}
        exit
    fi

    #copy config.ini in local path to supervisor path
    cp ${confpath} /etc/supervisord.d/

    cd ${pwd}

    printf "\n\ninstall finished..\n\n"

    echo
    return $RETVAL
}

start()
{
        echo -n $"Starting $prog: "
        # $SUPERVISORD -c $CONFIG_FILE --pidfile $PID_FILE && success || failure
        $SUPERVISORD -c $CONFIG_FILE
        RETVAL=$?
        echo
        return $RETVAL
}

stop()
{
        echo -n $"Stopping $prog: "
        supervisorctl stop all

        PID=`ps -eaf | grep $SUPERVISORD | grep -v grep | awk '{print $2}'`
        if [[ "" !=  "$PID" ]]; then
        echo "killing $PID"
        kill -9 $PID
        fi
        
        RETVAL=$?
        echo
        return $RETVAL
}

case "$1" in
        build)
                build
                ;;
        install)
                install
                ;;
        start)
                start
                ;;
        stop)
                stop
                ;;
        restart)
                stop
                start
                ;;
        status)
                status -p $PID_FILE $SUPERVISORD
                RETVAL=$?
                ;;
        *)
                echo $"Usage: sudo $0 {build|install|start|stop|restart|status}"
                RETVAL=1
esac
exit $RETVAL