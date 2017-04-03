

package userMgmt
// version No 1 dated :- 03-Apr-2017
import (
    
    "fmt"
    //"io/ioutil"
    // "encoding/json"
     //"net/http"
     "fileUtil"
     "agentUtil"
     //"os/exec" ---------------------------------
     //"stringUtil"

   // "os"
    "os/exec"
    _ "fmt" // for unused variable issue
  //  "net/smtp"
   // "log"
   // "strings" n
    
     //"reflect"
    // "strconv"
)

func AddUser(usrLoginName, preferredShell, pubKey string) string {
   
    removed_dir := "/home/deleted:" + usrLoginName
    home_dir := "/home/" + usrLoginName
      
     //Check whether user already existed or not. 
     status := agentUtil.ExecComand("id "+usrLoginName, "UserHandler.AddUser() L28");
     fmt.Println("33. UserHandler.AddUser()  status = : ", status)

     /* status ='fail' specify error,  'id usrLoginName' returns error due to absence of user existence
        So, below code block process to create new User Account
     */
    if(status == "fail") {
        if(fileUtil.IsFileExisted(home_dir) == false){
          if(fileUtil.IsFileExisted(removed_dir) == true){
            agentUtil.ExecComand("/bin/mv "+ removed_dir +" "+home_dir, "UserHandler.AddUser() L37")
          }
       }

       if(len(preferredShell) == 0){
          preferredShell =  "/bin/bash"
       }
       
      // Check whether group exists or not, if not, then create it
      
      
      cmd := " /usr/sbin/useradd "+ 
      "-m -d "+home_dir+                       // -d is unnecessary here, but will report error, if omit
      " -s "+preferredShell +
      //" -g "+usrLoginName+
      " "+ usrLoginName
      status = agentUtil.ExecComand(cmd, "UserHandler.AddUser() L60") 

      if(status == "success"){
        msg := "--------- user account "+usrLoginName+" successfully created ---------------"
        fileUtil.WriteIntoLogFile(msg)

        fmt.Println(msg)
        status = agentUtil.ExecComand(" chown -R "+ usrLoginName+":"+usrLoginName+ " "+ home_dir, "UserHandler.AddUser() L67")

        agentUtil.ExecComand("mkdir -p "+home_dir+"/.ssh", "UserHandler.AddUser() L71")
        fileUtil.WriteIntoFile(home_dir+"/.ssh/authorized_keys", pubKey, true)
       // status = agentUtil.ExecComand("chmod 700 "+home_dir+"/.ssh; chmod 640 "+home_dir+"/.ssh/authorized_keys", "UserHandler.AddUser() L74")
        status = agentUtil.ExecComand("chmod 777 "+home_dir+"/.ssh; chmod 777 "+home_dir+"/.ssh/authorized_keys", "UserHandler.AddUser() L74")
    

      }

    }else{
      fmt.Println("user already existed.")  
      fileUtil.WriteIntoLogFile("----- user already existed. usrLoginName = : "+usrLoginName)
    }
    return status
}

  
func Userdel(userLoginName string,  permanent bool)(string){
  permanent = false
  removed_dir := "/home/deleted:" + userLoginName
  home_dir := "/home/" + userLoginName
  userId :=  agentUtil.ExecComand("id -u "+userLoginName, "UserHandler.Userdel() L87");
  status := ""
  
  if(userId == "fail"){
    msg := "UserHandler.UserDel(). User does not exist"+userLoginName
    fileUtil.WriteIntoLogFile(msg)
    return "user does not existed"
  }
    
  if(permanent == false ){
    if(fileUtil.IsFileExisted(removed_dir)){
        agentUtil.ExecComand("/bin/rm -rf "+ removed_dir, "UserHandler.Userdel() L96")   
    }

    //Check below line in all version of linux after cross compile
    agentUtil.ExecComand("/usr/bin/pkill -u "+ userId, "UserHandler.Userdel() L100")      
    status = agentUtil.ExecComand("/usr/sbin/userdel "+ userLoginName, "UserHandler.Userdel() L101")      
    agentUtil.ExecComand("/bin/mv "+ home_dir +" "+removed_dir, "UserHandler.Userdel() L102")      
    
  }else{
    status = agentUtil.ExecComand("/usr/sbin/userdel -r "+ userLoginName, "UserHandler.Userdel() L105")      
  }
  Sudoers_del(userLoginName)
  return status
}
 

func Sudoers_del(userLoginName string){
  filePath := "/etc/sudoers.d/" + userLoginName
  if(fileUtil.IsFileExisted(filePath)){
    agentUtil.ExecComand("/bin/rm "+ filePath, "UserHandler.Userdel() L116")
  }
}


func GiveRootAccess(usrLoginName string) string{
 
  //Check whether user already existed or not. 
  status := agentUtil.ExecComand("id "+usrLoginName, "UserHandler.AddUser() L124");
  if(status == "fail"){
    fmt.Println("Unable to give root access to non existed user i.e ",usrLoginName)
    return status
  }
  

  //scriptPath := "/home/piyush/go_projects/scripts/sudoAdder.sh"
  scriptPath := "/opt/infraguard/etc/sudoAdder.sh"
  
  cmd := exec.Command("/bin/sh", "-c", scriptPath+" "+usrLoginName)
  output, err := cmd.Output()

  if err != nil {
    println(err.Error())
    msg := "UserHandler.GiveRootAccess(). Error on user "+usrLoginName+" Error Msg = : "+err.Error()
    fileUtil.WriteIntoLogFile(msg)
    return "1"
  }else{
    fmt.Println("File successfully edited...",(string(output)))
    msg := "UserHandler.GiveRootAccess(). Success on user "+usrLoginName+" Status = : "+string(output)
    fileUtil.WriteIntoLogFile(msg)
    return "0"
  }

}

func GiveNormalAccess(usrLoginName string) string{
  status := agentUtil.ExecComand("id "+usrLoginName, "UserHandler.AddUser() L151");
  if(status == "fail"){
    msg := "UserHandler.GiveNormalAccess(). user does not exist. Chk user = : "+usrLoginName
    fileUtil.WriteIntoLogFile(msg)
    fmt.Println(msg)
    return "1"
  }

  cmdStr := usrLoginName+"   ALL=(ALL:ALL) ALL"
  cmd := "sed -i '/"+cmdStr+"/s/^/#/' /etc/sudoers" 
  status = agentUtil.ExecComand(cmd, "UserHandler.GiveNormalAccess() L149");
  msg := "UserHandler.GiveNormalAccess(). Success for user = : "+usrLoginName
  fileUtil.WriteIntoLogFile(msg)
  return status

}
