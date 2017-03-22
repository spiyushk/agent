#kconfig: 35 90 12
# description: Agent Installer Test
#

# Get function from functions library

# Start the service AgentInstaller

start() {
echo $"***********  agent_controller service started ***********"
echo "Going to install Agent Code"
command="/opt/infraguard/sbin/infraGuardMain"


echo "Going to execute infraGuardMain executable."
daemon "nohup $command >/dev/null 2>&1 &"


}


stop(){
echo "Going to kill process agent_controller"
pkill  agent_controller

}

### main logic ###
case "$1" in
  start)
        start
        ;;
  stop)
        stop
        ;;

status)
        status agent_controller
        ;;

 *)
        echo $"Usage: $0 {start|stop|status}"
        exit 1
esac

exit 0



