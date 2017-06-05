#!/bin/bash


#. /etc/init.d/functions

# chkconfig: 35 90 12
# description: Agent Installer Test
#

# Get function from functions library

# Start the service AgentInstaller

 get_osflavor(){
    if [[ -f "/etc/lsb-release" ]]
        then
            os="ubuntu"
        elif [[ -f "/etc/redhat-release" ]]
        then
            os="rpm"
        elif [[ -f "/etc/debian_version" ]]
        then
            os="debian"
        else
            os="unknown"
            
    fi
}         



start() {

echo ""
echo $"***********  agent_controller service started. Triggered from /etc/init.d/agent_controller_ubuntu.sh ***********"
command="/opt/infraguard/sbin/infraGuardMain"
#daemon "nohup $command >/dev/null 2>&1 &"
#$command &>/dev/null &

$command  > /dev/null 2>&1
exit


#disown $command &


}


stop(){
echo "Going to kill process agent_controller_ubuntu.sh"
pkill  agent_controller_ubuntu.sh

}

### main logic ###
case "$1" in
  start)
        get_osflavor
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
os=""