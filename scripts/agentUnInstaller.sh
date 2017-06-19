#!/bin/bash

uninstall(){
 
getValue "serviceFile"    
getValue "removeProcessCmd"    

 echo "serviceFile = : $serviceFile"
 echo "removeProcessCmd = : $removeProcessCmd"

<<COMMENT1
Since serviceFile file exist in /etc/init.d/ directory. So ensure proper file to
ignore accidental deletion of entire contents in /etc/init.d folder. 
COMMENT1

if [[ $serviceFile != *"agent_controller"*  ||
         $removeProcessCmd != *"agent_controller"* ]]; then
   echo "No valid service file found. Abort process..."
   exit 1
fi


echo "Going to kill process $serviceFile"
pkill  $serviceFile


echo "Stopping the service..."
command="/etc/init.d/$serviceFile stop"
$command


pId=$(ps -ef | grep 'infraGuardMain' | grep -v 'grep' | awk '{ printf $2 }')
echo "Stopping the process i.e infraGuardMain."
command="/bin/kill -9 $pId"
$command


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


if [ `id -u` -ne 0 ] ; then
            echo "error: Agent uninstallation process requires superuser privilege. Abort process."
            exit 1
fi

#fileName="/tmp/agentInfo.txt" 
fileName="/opt/infraguard/etc/agentInfo.txt" 
if [ ! -f $fileName ]; then
    echo "Missing file $fileName. Abort uninstallation process."
    exit 1
fi


command="cp -r  /etc/init.d  /tmp/"
$command

serviceFile=""
removeProcessCmd=""
uninstall


