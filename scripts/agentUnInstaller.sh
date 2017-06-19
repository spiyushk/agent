#!/bin/bash

uninstall(){
 
getValue "serviceFile"    
getValue "removeProcessCmd"    


# value of serviceFile was saved at the time of agent installation and it may be
# agentInstaller.sh/agentInstaller.service/agentUnInstaller_ubuntu.sh
echo "serviceFile = : $serviceFile"
echo "removeProcessCmd = : $removeProcessCmd"


# Since serviceFile file exist in /etc/init.d/ directory. So ensure proper file MUST BE EXIST( to 
# ignore accidental deletion of entire contents in /etc/init.d folder)
if [[ $serviceFile != *"agent_controller"*  ||
         $removeProcessCmd != *"agent_controller"* ]]; then
   echo "No valid service file found. Abort process..."
   exit 1
fi


echo "Stopping the service..."
command="/etc/init.d/$serviceFile stop"
$command


pId=$(ps -ef | grep 'infraGuardMain' | grep -v 'grep' | awk '{ printf $2 }')
echo "Stopping the process i.e infraGuardMain."
command="/bin/kill -9 $pId"
$command


# Restrict service to restart on reboot
# On the basis of linux type either it has update-rc.d -f ... or chkconfig --del  ....
# value of removeProcessCmd was saved at the time of agent installation
$removeProcessCmd 

echo "Process $pId killed successfully " 


echo "Deleting all concerned directories ..."
command="rm -rf /opt/infraguard/"
$command

command="rm -rf /var/logs/infraguard/"
$command

command="rm -rf /etc/init.d/$serviceFile"
echo "---------------- Full Command  init.d removel --> $command"
$command

echo ""
echo "Uninstallation process completes."

} #Uninstall



# Read the given 'key' from /opt/infraguard/etc/agentInfo.txt file
# On the basis of key, shared variable will be initialized and uses in uninstall() method.
getValue(){
   key="$1"
   while IFS= read -r line; do

      if [[ $line == *"$key"* ]]; then
         val=${line/$key=/""}    
   
         if [[ $line == "serviceFile"* ]]; then
              serviceFile=$val
         fi

         if [[ $line == "removeProcessCmd"* ]]; then
              removeProcessCmd=$val
         fi

         break;
      fi

  done < "$fileName"
}

# Check whether user has root level access or not.
if [ `id -u` -ne 0 ] ; then
            echo "error: Agent uninstallation process requires superuser privilege. Abort process."
            exit 1
fi


# Check whether file existed or not.
fileName="/opt/infraguard/etc/agentInfo.txt" 
if [ ! -f $fileName ]; then
    echo "Missing file $fileName. This file should be created at the time of agent installation."
    exit 1
fi


serviceFile=""
removeProcessCmd=""
uninstall


