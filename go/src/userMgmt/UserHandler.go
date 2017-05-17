

package userMgmt
// version No 1 dated :- 03-Apr-2017
import (
    
    "fmt"
    //"io/ioutil"
    // "encoding/json"
     //"net/http"

     //  "io/ioutil"
     //  "net/http"
    
     //"reflect"
    // "strconv"
    
     "fileUtil"
     "agentUtil"
     //"os/exec" ---------------------------------
     //"stringUtil"

   // "os"
    //"os/exec"
    _ "fmt" // for unused variable issue
  //  "net/smtp"
   // "log"
      "strings" 
      "stringUtil"

)

func AddUser(usrLoginName, preferredShell, pubKey string) string {
   
    removed_dir := "/home/deleted:" + usrLoginName
    home_dir := "/home/" + usrLoginName
    status := ""
    
    fmt.Println("usrLoginName = : ", usrLoginName)
    fmt.Println("preferredShell = : ", preferredShell)
    fmt.Println("pubKey = : ", pubKey)

    fileUtil.WriteIntoLogFile("Going to create new user account for user "+usrLoginName)

    //return ""
    if(isUserExist(usrLoginName) == false) {
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
     

        agentUtil.ExecComand("mkdir -p "+home_dir+"/.ssh", "UserHandler.AddUser() L71")

       
       
        rsaFileName := "authorized_keys"
        fileUtil.WriteIntoFile(home_dir+"/.ssh/"+rsaFileName, pubKey, false, true)
        status = agentUtil.ExecComand("chmod 700 "+home_dir+"/.ssh; chmod 600 "+home_dir+"/.ssh/"+rsaFileName, "UserHandler.AddUser() L74")
        status = agentUtil.ExecComand(" chown -R "+ usrLoginName+":"+usrLoginName+ " "+ home_dir, "UserHandler.AddUser() L67")
        fmt.Println(msg)

      }

    }else{
      fmt.Println("user already existed.")  
      fileUtil.WriteIntoLogFile("----- user already existed. usrLoginName = : "+usrLoginName)
    }
    return status
}

  
func Userdel(userLoginName string,  permanent bool)(string){
 
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
    status = agentUtil.ExecComand("sudo /usr/sbin/userdel -r "+ userLoginName, "UserHandler.Userdel() L105")      
  }
  Sudoers_del(userLoginName)
  return status
}
 

func Sudoers_del(userLoginName string){
  filePath := "/etc/sudoers.d/" + userLoginName
  if(fileUtil.IsFileExisted(filePath)){
    agentUtil.ExecComand("/bin/rm "+ filePath, "UserHandler.Userdel() L126")
  }
}

/*
  Uses :- To give sudo permission, If user is not a sudo user or a sudo user converted to normal user
*/
func ProcessToChangePrivilege(usrName, privType string) string{
    
    rootPrivGrpName := GetSudo_GrpName() // On ubuntu, it is 'sudo' & on fedora, it is 'wheel'
    if(len(rootPrivGrpName) == 0){
         msg := "UserHandler.ProcessToChangePrivilege(). Unable to locate sudo/wheel. "+
                    "May be it is commented. Abort further process."
         fileUtil.WriteIntoLogFile(msg)
         fmt.Println("166. msg ", msg)
         return "1"
    }
    msg := ""
    cmd := ""
    grpOfUser := getUser_AllGrp(usrName) // grpOfUser stores, all group name in which user is a member
    isUsrHasRootPriv := strings.Contains(grpOfUser, rootPrivGrpName)

    if(privType == "root" && isUsrHasRootPriv){
        msg = "User has already sudo permission. "+usrName +" From ProcessToChangePrivilege(). L148. "
        fileUtil.WriteIntoLogFile(msg)
        return "0"
    }

     if(privType != "root" && isUsrHasRootPriv == false){
        msg = "User has not sudo permission. "+usrName +" Nothing to do . From ProcessToChangePrivilege(). L154. "
        fileUtil.WriteIntoLogFile(msg)
        return "0"
    }
   
    status := ""
    if(privType == "root"){
      userPwd :=  ChangePwd(usrName)
      if(userPwd == "1"){   // 1 indicate pwd does not created, so abort.
          msg = "UserHandler.ProcessToChangePrivilege(). Unable to create new pwd. Abort further process. L 162"
          fileUtil.WriteIntoLogFile(msg)
          fmt.Println("166. msg ", msg)
          return "1"
      }
       cmd = "usermod -aG "+rootPrivGrpName+" "+usrName
       fmt.Println("154. Going to run command = : ", cmd)
       status = agentUtil.ExecComand(cmd, "UserHandler.ProcessToChangePrivilege() L170")
       fmt.Println("\n\n 156. ****************  Status ", status)
       msg = cmd + " >> Status = : "+status
       fileUtil.WriteIntoLogFile(msg)
      
       if(status == "success"){
          return userPwd
       }
    
    }

    if(privType != "root"){
      cmd = ""
      if(rootPrivGrpName == "wheel"){ // For fedora
        cmd = "gpasswd -d  "+usrName +" "+rootPrivGrpName
        status = agentUtil.ExecComand(cmd, "UserHandler.ProcessToChangePrivilege() L184")
      }

      if(rootPrivGrpName == "sudo"){    // For ubuntu
         cmd = "deluser "+usrName +" "+rootPrivGrpName
         status = agentUtil.ExecComand(cmd, "UserHandler.ProcessToChangePrivilege() L189")
      }
     
      msg = cmd + " >> Status = : "+status
      fileUtil.WriteIntoLogFile(msg)
      fmt.Println("\n\n 194. ****************  msg ", msg)
      if(status == "success"){
        return "0"
      }
    }

   return "1"
}


/*
  Since different distro have different sudo group. e.g in ubuntu it is 'sudo'
  wheaeas  in fedora , sudo group is 'wheel' group
  It is assumed that sudoers group either sudo/wheel are uncommented wherever applicable.
*/
func GetSudo_GrpName()string{
  rootPrivGrpName := ""
     status := agentUtil.ExecComand("getent group sudo", "UserHandler.GetSudo_GrpName() L211")
        if(status == "fail"){
            status = agentUtil.ExecComand("getent group wheel", "UserHandler.GetSudo_GrpName() L213")
            if(status != "fail"){
              rootPrivGrpName = "wheel"
            }
        }else{
          rootPrivGrpName = "sudo"
        }
   return rootPrivGrpName    
}

// get all those group name in which user is a member
func getUser_AllGrp(usrName string) string{  
  groupNames := agentUtil.ExecComand("id -nG "+usrName, "UserHandler.getUser_AllGrp() L225")
  return groupNames
}


// To become  a root user, user must have own password
func ChangePwd(usrName string) string{
   userPwd := usrName + stringUtil.GetRandomString(4)
   fmt.Println("randomStr  on 4 = : ", userPwd)
 
   cmd := "usermod --password $(echo "+userPwd+" | openssl passwd -1 -stdin) "+usrName
   fileUtil.WriteIntoLogFile("Going to execute this command = : "+cmd)
   status := agentUtil.ExecComand(cmd, "UserHandler.ChangePwd() L237")
   fmt.Println("239. Status ChangePwd = : ", status)
   fileUtil.WriteIntoLogFile(status)
   msg := " New Pwd for usrName = : "+usrName + " Is >> "+userPwd
   fileUtil.WriteIntoLogFile("\n"+msg)
   fileUtil.WriteIntoLogFile("\n")


   if(status == "success"){
      return userPwd
   }
   return "1"

}

func isUserExist(usrName string) bool{
  status := agentUtil.ExecComand("id "+usrName, "UserHandler.isUserExist() L193");
  fmt.Println("33. UserHandler.AddUser()  status = : ", status)

     /* status ='fail' specify error,  'id usrLoginName' returns error due to absence of user existence
        So, below code block process to create new User Account
     */
    if(status == "fail") {
      return false;
    }
    return true;
}

// i
func UserAccountController(activityName string, nextWork []string, callerLoopCntr int) (int){
    var pubKey, userName, prefShell, privilege, id string
    var values []string
    
    if(activityName == "addUser"){
     
        responseUrl := "https://spjuv2c0ae.execute-api.us-west-2.amazonaws.com/dev/addeduserbyagent"
        values = stringUtil.SplitData(nextWork[callerLoopCntr+1], agentUtil.Delimiter)
        pubKey = values[1]


        values = stringUtil.SplitData(nextWork[callerLoopCntr+2], agentUtil.Delimiter)
        //userName = values[1] + stringUtil.RandStringBytes(3)
        userName = values[1]
        
        values = stringUtil.SplitData(nextWork[callerLoopCntr+3], agentUtil.Delimiter)
        prefShell = values[1]

        values = stringUtil.SplitData(nextWork[callerLoopCntr+4], agentUtil.Delimiter)
        id = values[1]

      
        msg :=  "Going to add userName = : "+userName
        fmt.Println(msg)
        fileUtil.WriteIntoLogFile(msg)
        status := AddUser(userName, prefShell, pubKey ) 
       
        agentUtil.SendExecutionStatus(responseUrl, status , id, userName)
      
        callerLoopCntr += 4
        return callerLoopCntr

    }

    if(activityName == "deleteUser"){
        responseUrl := "https://vglxmaiux1.execute-api.us-west-2.amazonaws.com/dev/deleteduserbyagent"
        values = stringUtil.SplitData(nextWork[callerLoopCntr+1], agentUtil.Delimiter)
        userName = values[1]

        values = stringUtil.SplitData(nextWork[callerLoopCntr+2], agentUtil.Delimiter)
        id = values[1]

        msg :=  "Going to delete userName = : "+userName
        fmt.Println(msg)
        fileUtil.WriteIntoLogFile(msg)

        status := Userdel(userName, false)
        fmt.Println("status deleteUser  = : ", status)
        agentUtil.SendExecutionStatus(responseUrl, status , id, userName)
        callerLoopCntr += 2
        return callerLoopCntr
    }

    if(activityName == "changePrivilege"){
      responseUrl := "https://a1gpcq76u3.execute-api.us-west-2.amazonaws.com/dev/privilegechangedbyagent"

     
      status := ""
        values = stringUtil.SplitData(nextWork[callerLoopCntr+1], agentUtil.Delimiter)
        userName = values[1]

        values = stringUtil.SplitData(nextWork[callerLoopCntr+2], agentUtil.Delimiter)
        privilege = values[1]

        values = stringUtil.SplitData(nextWork[callerLoopCntr+3], agentUtil.Delimiter)
        id = values[1]

        
        msg :=  "Going to change privilege for userName = : "+userName+ " Priv = : "+privilege
        fmt.Println(msg)
        fileUtil.WriteIntoLogFile(msg)

        if(isUserExist(userName) == false) {
          msg = "Request to change priviliges for non existed user "+userName +" --> Abort rest process."
          fileUtil.WriteIntoLogFile(msg)
          status = "1"
          
        }else{
           msg := "UserHandler.go L338. ProcessToChangePrivilege. usrName = : "+userName+" >> Requested privilege. Type = : "+privilege
           fmt.Println(msg)
           fileUtil.WriteIntoLogFile(msg)
           status = ProcessToChangePrivilege(userName, privilege)
        }
      
      // if status length is > 4, it means status stores user's new pwd

         msg = "\nUserHandler.go L346. Final status of ProcessToChangePrivilege. status = : "+status
         fmt.Println(msg)
         fileUtil.WriteIntoLogFile(msg)


        if(len(status) > 4){
            usrNewPwd := status
            agentUtil.SendExecutionStatus(responseUrl, "0" , id, userName, usrNewPwd) 
        }else{
            agentUtil.SendExecutionStatus(responseUrl, status , id, userName) 
        }
      
        callerLoopCntr += 3
        return callerLoopCntr
    }

    /*
      ----------------------------------  Lock down server -------------------------------------
      Pssible data format is given below
      activityName:lockDownServer requiredData:{"userList":"ec2-user,pratyush,sampath,piyush,prashant.gyan"} id:5]
    */

     if(activityName == "lockDownServer"){
      
        responseUrl := "https://h80y20gh11.execute-api.us-west-2.amazonaws.com/dev/serverlockeddown"
        status := ""
        var userList []string 
        values = stringUtil.SplitData(nextWork[callerLoopCntr+1], agentUtil.Delimiter)
        if(len(values) == 2){
          userList = stringUtil.SplitData(values[1], ",") 
          fmt.Println("userList from api = : ", userList)

          //userList = getUserList() 
          values = stringUtil.SplitData(nextWork[callerLoopCntr+2], agentUtil.Delimiter)
          id = values[1]
        }
      
       
        fmt.Println("Going to lock Down Server. Deletable users are = : ",userList)
        fileUtil.WriteIntoLogFile("Going to lock Down Server. Following users going to lock = : "+strings.Join(userList,","))
       
        // Below callerLoopCntr is used to control the loop iteration in caller function.
        callerLoopCntr += 2

        for j := 0; j < len(userList); j++{
            userName = userList[j]

            // To stop accidental lock down from local host 
            if(strings.Contains(userName, "piyush")){
              continue;
            }

           /*
             disallow userName from logging in --> sudo usermod --expiredate 1 userName
             set expiration date of userName to Never :- sudo usermod --expiredate "" userName
           */
            status = agentUtil.ExecComand("usermod --expiredate 1 "+ userName, "UserHandler.lockDownServer() L326")
            fmt.Println("status to lock user = : ",status)

            msg :=  "Locking status of user =: "+userName +" is "+status
            fmt.Println(msg)
            fileUtil.WriteIntoLogFile(msg)
         }

        agentUtil.SendExecutionStatus(responseUrl, status , id) 
        fmt.Println("334. UserAccountController status lockDownServer  = : ", status) 
        return callerLoopCntr
    }
    return callerLoopCntr

  }//UserAccountController

  


  