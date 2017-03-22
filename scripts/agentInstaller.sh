#!/bin/bash
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
            echo "ERROR: Cannot get the system type. Exiting."
            os="unknown"
            exit 1
    fi
}

<<"COMMENT"

start() {
(

prog=awslogs
exec="/usr/sbin/awslogsd"
lockfile=/var/lock/subsys/awslogs
pidfile=/var/run/awslogs.pid
mutexfile=/var/lock/awslogs.mutex

[ -x $exec ] || exit 5
echo -n $"Starting $prog: "
daemon $NICELEVEL --pidfile=$pidfile --check=${prog} "nohup $exec >/dev/null 2>&1 &"
retval=$?
echo [ $retval -eq 0 ] && touch $lockfile
) 9>${mutexfile}
rm -f ${mutexfile}
}



COMMENT

create_InfraGuardDirectories(){
    echo "Creating directories in /opt ..."
    exec="sudo mkdir -p /opt/infraguard/sbin"
    $exec

    exec="sudo mkdir -p /opt/infraguard/etc"
    $exec

    exec="sudo mkdir -p /var/logs/infraguard"
    $exec
    

    exec="sudo chmod 777 /opt/infraguard/sbin"
    $exec

    exec="sudo chmod 777 /opt/infraguard/etc"
    $exec

    exec="sudo chmod 777 /var/logs/infraguard"
    $exec
    
    echo "completed Directories Creation"

}



install_daemon(){
    echo 'Attempting Daemon Installation'
    cd /tmp
    if [[ "$os" = "debian"  || "$os" = "ubuntu" ]]
        then
        wget -q https://github.com/terminalcloud/terminal-tools/raw/master/daemon_0.6.4-2_amd64.deb || exit -1
        dpkg -i daemon_0.6.4-2_amd64.deb
        echo "Daemon Installation Done Successfully on debian/ununtu"
    else
        wget -q http://libslack.org/daemon/download/daemon-0.6.4-1.x86_64.rpm || exit -1
        rpm -i daemon-0.6.4-1.x86_64.rpm
        echo "Daemon Installation Done Successfully on rpm"
    fi
}

downloadFiles_FromGitHub() {
    # dpkg: error: requested operation requires superuser privilege
     echo "Server Name: $serverName"
     echo "Project Id: $projectId"
     echo "licenseKey: $licenseKey"
     echo "$serverName:$projectId:licenseKey" > /tmp/serverInfo.txt

    echo "Downloading agent_controller.sh and moving it from /tmp to /etc/init.d ..."
    echo ""
    local url="wget -O /tmp/agent_controller.sh https://raw.githubusercontent.com/agentinfraguard/agent/master/scripts/agent_controller.sh"
    wget $url--progress=dot $url 2>&1 | grep --line-buffered "%" | sed -u -e "s,\.,,g" | awk '{printf("\b\b\b\b%4s", $2)}'
    
    ######command="sudo mv /tmp/agent_controller.sh  /etc/init.d"
    command="mv /tmp/agent_controller.sh  /etc/init.d"
    $command
    #chkconfig --add /etc/init.d/agent_controller.sh

    


    #####command="sudo chmod 777 /etc/init.d/agent_controller.sh"
    command="chmod 777 /etc/init.d/agent_controller.sh"
    $command
    

    echo "Downloading infraGuardMain executable. Please wait...."
    url="wget -O /opt/infraguard/sbin/infraGuardMain https://raw.githubusercontent.com/agentinfraguard/agent/master/go/src/agentController/infraGuardMain"
    wget $url--progress=dot $url 2>&1 | grep --line-buffered "%" | sed -u -e "s,\.,,g" | awk '{printf("\b\b\b\b%4s", $2)}'
    echo "infraGuardMain downloaded."

    ###command="sudo chmod 777 /opt/infraguard/sbin/infraGuardMain"
    command="chmod 777 /opt/infraguard/sbin/infraGuardMain"
    $command
    

    export start="start"
    export command="/etc/init.d/agent_controller.sh"
        
    sh $command ${start}
    

# bash <(wget -qO- https://raw.githubusercontent.com/agentinfraguard/agent/master/scripts/agentInstaller.sh) server111 5011 lKey101
   
    }

    checkUserPrivileges(){
        if [ `id -u` -ne 0 ] ; then
            echo "error: requested operation requires superuser privilege"
            exit 1
        fi
    }



echo "$#"
if [ $# -ne 3 ] ; then
    echo "Insufficient arguments. Usage: $0 serverName projectId licenseKey"
    exit 1
fi
checkUserPrivileges
# Read arguments, it will saved into /tmp/serverInfo.txt & then serverMgmt/ServerHandler.go will read.
serverName=$1
projectId=$2
licenseKey=$3

os=""
get_osflavor
#install_daemon
echo "os found = : $os"
create_InfraGuardDirectories
downloadFiles_FromGitHub






