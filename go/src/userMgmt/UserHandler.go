

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

        fmt.Println(msg)
        status = agentUtil.ExecComand(" chown -R "+ usrLoginName+":"+usrLoginName+ " "+ home_dir, "UserHandler.AddUser() L67")

        agentUtil.ExecComand("mkdir -p "+home_dir+"/.ssh", "UserHandler.AddUser() L71")
        fileUtil.WriteIntoFile(home_dir+"/.ssh/authorized_keys", pubKey, false, true)
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
    agentUtil.ExecComand("/bin/rm "+ filePath, "UserHandler.Userdel() L116")
  }
}

func ProcessToChangePrivilege(usrName, privType string) string{
    if(isUserExist(usrName) == false) {

      msg := "Unable to change priviliges for user "+usrName +" : This user not existed"
      fileUtil.WriteIntoLogFile(msg)
      return "1"
    }
    accessRight := usrName+"   ALL=(ALL:ALL) ALL"
    
    
    // Read /etc/sudoers file for user access right, if any
    cmd := "awk '/"+usrName+"/ {print}' /etc/sudoers"
    oldPriv := agentUtil.ExecComand(cmd, "misc. L78")
    fmt.Println("79 oldPriv : = ",oldPriv)
    if(oldPriv == "success"){
        oldPriv = ""
    }

    tmpFilePath := "/tmp/sudoers.bak"
    status := ""
  
    // Create  a back up copy of /etc/sudoers file
    status = agentUtil.ExecComand("cp /etc/sudoers "+tmpFilePath, "misc. L38")
    fmt.Println("status at 26.  = : ", status)

    status = agentUtil.ExecComand("chmod 777 "+tmpFilePath, "misc. L41")
    fmt.Println("status at 29.  = : ", status)


    if( privType == "root"){
        if(len(oldPriv) == 0){
                    
            fileUtil.WriteIntoFile(tmpFilePath, accessRight, true, false)
            fmt.Println("condition matched at 52.  = : ")
            // Replace the sudoers file with the tmpFilePath 
            status = agentUtil.ExecComand("cp "+tmpFilePath +" /etc/sudoers", "misc. L38")

            msg := " userName = "+usrName+" Requested access = : "+ privType +" done. New entry in /etc/sudoers file. Status = : "+status
            fileUtil.WriteIntoLogFile(msg)
            
        }else{
             // If old priv is commented, then remove such comment
            if(strings.Contains(oldPriv, "#")){
                status = fileUtil.ReplaceLineOrLinesIntoFile(tmpFilePath, oldPriv, accessRight)
               
                // Replace the sudoers file with the tmpFilePath 
                status = agentUtil.ExecComand("cp "+tmpFilePath +" /etc/sudoers", "misc. L38")

                msg := " userName = "+usrName+" Requested access = : "+ privType +" done. Previously access right commented in /etc/sudoers file. Status = : "+status
                fileUtil.WriteIntoLogFile(msg)
            }   
        }
    }

    // If user's old priv is root level access, then comment access right data to become a normal user
    if( privType == "user"){
        if(len(oldPriv) > 0){
            status = fileUtil.ReplaceLineOrLinesIntoFile(tmpFilePath, oldPriv, "#"+accessRight)
            // Replace the sudoers file with the tmpFilePath 
            status = agentUtil.ExecComand("cp "+tmpFilePath +" /etc/sudoers", "misc. L38")

            msg := " userName = "+usrName+" Requested access = : "+ privType +" done. Previously access right has root Access. Now it is commented in /etc/sudoers file. Status = : "+status
            fileUtil.WriteIntoLogFile(msg)
        }
    }

    msg := " userName = "+usrName+" Requested access = : "+ privType +" done. /etc/sudoers file is unaffected. Status = : "+status
    fileUtil.WriteIntoLogFile(msg)
    return status
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
// i
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

        status = ProcessToChangePrivilege(userName, privilege)
        agentUtil.SendExecutionStatus(responseUrl, status , id, userName) 
       
        fmt.Println("status changePrivilege  = : ", status) 
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
        //responseUrl := "https://h80y20gh11.execute-api.us-west-2.amazonaws.com/dev/serverlockeddown?id=5&serverIp=172.31.15.1&status=0"
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

  
  /*func getUserList()([]string){
   // b := [5]string{"tecmint","testUser","rajeevSir1","rajeevSir2", "sadfsdfsdf"}
     users := []string{"rajeevSir2"}
     return users
  }*/

  