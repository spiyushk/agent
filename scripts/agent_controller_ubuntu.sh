#!/bin/bash


#. /etc/init.d/functions

# chkconfig: 35 90 12
# description: Agent Installer Test

# Get function from functions library
# Start the service AgentInstaller

 

start() {

echo ""
echo $"***********  agent_controller service started. Triggered from /etc/init.d/agent_controller_ubuntu.sh ***********"
command="/opt/infraguard/sbin/infraGuardMain"
$command > /dev/null 2>&1 &

}


stop(){
echo "Going to kill process agent_controller_ubuntu.sh"
pkill  agent_controller_ubuntu.sh

}


case "$1" in
  start)
        start
        ;;
  stop)
        stop
        ;;

status)
        status agent_controller_ubuntu.sh
        ;;

 *)
        echo $"Usage: $0 {start|stop|status}"
        exit 1
esac

exit 0
