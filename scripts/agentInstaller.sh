#!/bin/bash


create_InfraGuardDirectories(){
    echo "Creating directories in /opt ..."
    exec="mkdir -p /opt/infraguard/sbin"
    $exec

    exec="mkdir -p /opt/infraguard/etc"
    $exec

    exec="mkdir -p /var/logs/infraguard"
    $exec

    exec="touch  /var/logs/infraguard/activityLog"
    $exec

   
    exec="chmod 700 -R /opt/infraguard"
    $exec
    exec="chown root:root /opt/infraguard"
    $exec


    exec="chmod 700 -R /var/logs/infraguard"
    $exec
    exec="chown root:root /var/logs/infraguard"
    $exec

    echo "completed Directories Creation"


}


getLinuxType(){

   filename="/opt/infraguard/etc/linuxDistroInfo.txt" 
   cat /etc/*-release > $filename
   
   while IFS= read -r line; do

      if [[ $line == *"ID_LIKE"* ]]; then
         echo "$line"
         
         osType=${line/ID_LIKE=/""} # Extract string after "=" i.e ID_LIKE="fedora"
         osType="${osType%\"}" # Remove dbl quotes - suffix
         osType="${osType#\"}" # Remove dbl quotes - prefix
       
                
          if [[ $osType == "debian" ]]; then
             os="debian"
             fileAgentController="agent_controller_ubuntu.sh"
          fi


          if [[ $osType == "fedora" ]]; then
             os="fedora"
             fileAgentController="agent_controller.service"
          fi
        break;

      fi

  done < "$filename"


}



#  There are two repository on github, infraguard & spiyushk. infraguard is for prod environment & spiyushk is for
#  testing purpose. Below method will get file name to sownload from intended repository.

getFilePath(){
    repoName="$1"
    fileName="$2"
    echo "Repo Name = : $repoName"
    echo "File Name = : $fileName"
    gitFullPath=""

    if [[ $fileName == "agent_controller.sh"  ||
         $fileName == "agent_controller.service" ||
         $fileName == "agent_controller_ubuntu.sh" ]]; then
       gitFullPath="https://raw.githubusercontent.com/$repoName/agent/master/scripts/$fileName"

    fi
# https://raw.githubusercontent.com/agentinfraguard/agent/master/go/src/agentController/infraGuardMain"
    if [[ $fileName == "infraGuardMain" ]]; then
       gitFullPath="https://raw.githubusercontent.com/$repoName/agent/master/go/src/agentController/infraGuardMain"
    fi

    if [[ $fileName == "agentConstants.txt" ]]; then
       gitFullPath="https://raw.githubusercontent.com/$repoName/agent/master/go/src/agentConstants.txt"
       # gitFullPath="https://github.com/$repoName/agent/blob/master/go/src/agentController/infraGuardMain"
    fi

}

installAgent() {
# bash <(wget -qO- https://raw.githubusercontent.com/spiyushk/agent/master/scripts/agentInstaller.sh) server111 6 lKey101 
    #repoName="spiyushk"
    repoName="agentinfraguard"

    getFilePath "$repoName" "$fileAgentController"
    echo "gitFullPath = : $gitFullPath"
    echo "Downloading $fileAgentController "
    #local url="wget -O /tmp/$fileAgentController https://raw.githubusercontent.com/agentinfraguard/agent/master/scripts/$fileAgentController"
    local url="wget -O /tmp/$fileAgentController $gitFullPath"
    wget $url--progress=dot $url 2>&1 | grep --line-buffered "%" | sed -u -e "s,\.,,g" | awk '{printf("\b\b\b\b%4s", $2)}'
    command="mv /tmp/$fileAgentController  /etc/init.d"
    $command
    exec="chown root:root /etc/init.d/$fileAgentController"
    $exec
    exec="chmod 755 /etc/init.d/$fileAgentController"
    $exec
    echo "gitFullPath = : $gitFullPath"

    echo ""  
    echo "create  /tmp/serverInfo.txt with following data $serverName:$projectId:$licenseKe >> It will remove after server regn."
    echo "$serverName:$projectId:licenseKey" > /tmp/serverInfo.txt



    gitFullPath=""
    getFilePath "$repoName" "infraGuardMain"
    echo "gitFullPath = : $gitFullPath"
    echo ""
    echo "Downloading infraGuardMain executable. It will take time. Please wait...."
    url="wget -O /opt/infraguard/sbin/infraGuardMain $gitFullPath"
    
    #url="wget -O /opt/infraguard/sbin/infraGuardMain https://raw.githubusercontent.com/agentinfraguard/agent/master/go/src/test/infraGuardMain"
    wget $url--progress=dot $url 2>&1 | grep --line-buffered "%" | sed -u -e "s,\.,,g" | awk '{printf("\b\b\b\b%4s", $2)}'
    echo "infraGuardMain downloaded."
    exec="chown root:root /opt/infraguard/sbin/infraGuardMain"
    $exec
    exec="chmod 700 /opt/infraguard/sbin/infraGuardMain"
    $exec

    echo "153. gitFullPath = : $gitFullPath"


    gitFullPath=""
    getFilePath "$repoName" "agentConstants.txt"
    echo "gitFullPath = : $gitFullPath"

    echo ""
    echo "Downloading property file i.e agentConstants.txt ...."
    url="wget -O /opt/infraguard/etc/agentConstants.txt $gitFullPath"
    wget $url--progress=dot $url 2>&1 | grep --line-buffered "%" | sed -u -e "s,\.,,g" | awk '{printf("\b\b\b\b%4s", $2)}'
    echo "agentConstants.txt downloaded."

    echo "152. gitFullPath = : $gitFullPath"
    gitFullPath=""
    exec="chown root:root /opt/infraguard/etc/agentConstants.txt"
    $exec
    exec="chmod 700 /opt/infraguard/etc/agentConstants.txt"
    $exec


   

     if [[ "$os" = "debian" ]] ;then
            echo "Going to call  update-rc.d for $fileAgentController --------"
            update-rc.d $fileAgentController defaults
     else
            echo "Going to call  chkconfig for $fileAgentController --------"
             chkconfig --add /etc/init.d/$fileAgentController     
     fi
 


     export start="start"

     # Since fedore automatically added '.service' suffix in file name, so here ignore file extn
     if [[ $os == "fedora" ]]; then
         export command="/etc/init.d/agent_controller" 
     else    
         export command="/etc/init.d/$fileAgentController"
     fi
        
     sh $command ${start}

   
    } # downloadFiles_FromGitHub


    checkUserPrivileges(){
        if [ `id -u` -ne 0 ] ; then
            echo "error: requested operation requires superuser privilege"
            exit 1
        fi
    }


# Check whether agent already is running or not. If yes, then abort further process.

echo "Checking whether agent already installed/running or not."
pId=$(ps -ef | grep 'infraGuardMain' | grep -v 'grep' | awk '{ printf $2 }')
file="/opt/infraguard/sbin/infraGuardMain"

if [ -f "$file" ]
  then
    echo "Agent exe file found at $file "

    if [ -z "$pId" ] ; then
        echo "Agent is stopped."
    else
        echo "Agent is running. Process id is $pId"
    fi

  echo "Abort further process."
  exit 1

   
fi


if [ $# -ne 3 ] ; then
    echo "Insufficient arguments. Usage: $0 serverName projectId licenseKey"
    exit 1
fi


checkUserPrivileges
# Read arguments, it will saved into /tmp/serverInfo.txt & then serverMgmt/ServerHandler.go will read.
serverName=$1
projectId=$2
licenseKey=$3
gitFullPath=""

# Default value for os & fileAgentController is based on Amazon Linux AMI i.e rhel fedora
os="rhel fedora"
fileAgentController="agent_controller.sh"

create_InfraGuardDirectories
getLinuxType

echo "fileAgentController = : $fileAgentController"
echo "OS = : $os"
installAgent






