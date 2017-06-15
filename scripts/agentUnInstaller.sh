#!/bin/bash

uninstall(){
 
getValue "/opt/infraguard/etc/agentInfo.txt" "installer"    
getValue "/opt/infraguard/etc/agentInfo.txt" "removeProcessCmd"    

echo "Concerned Installer $installerName"
echo "Removal Command $removeProcessCmd"

echo "Going to kill process $installerName"
pkill  $installerName


echo "Stopping the service..."
command="/etc/init.d/$installerName stop"
$command

echo "Stopping the process..."
pId=$(ps -ef | grep 'infraGuardMain' | grep -v 'grep' | awk '{ printf $2 }')
echo "infraGuardMain ProcessId = : $pId"
command="/bin/kill -9 $pId"
$command


if [ $? != 0 ]; then                   
   echo "Unable to kill process id $pId " 
else
   $removeProcessCmd 
   echo "Process $pId killed successfully " 
fi

echo "Deleting all concerned directories with files..."
command="rm -rf /opt/infraguard/"
$command

command="rm -rf /var/logs/infraguard/"
$command

command="rm -rf etc/init.d/$installerName"
$command

echo ""
echo "Uninstallation process completes."

} #Uninstall

getValue(){
   filename="/opt/infraguard/etc/agentInfo.txt" 
   key="$2"
   while IFS= read -r line; do

      if [[ $line == "$key" ]]; then
         echo "$line"
         installerName=${line/$key=/""}
         break;

      fi

  done < "$filename"
}


installerName=""
removeProcessCmd=""

uninstall


