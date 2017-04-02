

package main

import (
    
   /* 
    "io/ioutil"
    "os"
    "os/exec"
    _ "fmt" // for unused variable issue
    "net/smtp"
    "log"
    "strings"
    "encoding/json"
    "net/http"
    */
    //"agentUtil"
    "stringUtil"
    "serverMgmt"
    "fmt"
    "fileUtil"
   // "userMgmt"
    "agentUtil"
    "userMgmt"
    "github.com/jasonlvhit/gocron"  // go get github.com/robfig/cron  
    //"strconv"
)
var freqToHitApi_InSeconds uint64 = 20

/*func main() {
    nextWork := agentUtil.GetNextWork()
    if(nextWork != nil){
    fmt.Println("\n\nInfraGuard.main(). Length infraGuardResponse = : ",len(nextWork)) 
    ExecuteWork(nextWork)
    }else{
       fmt.Println("InfraGuard.main(). There is no new work") 
    }
  
}


*/


func main() {
  fmt.Println("InfraGuard.main()") 
  respStr :=serverMgmt.DoServerRegnProcess()
  
  if(respStr =="0"){
    fmt.Printf("\nServer Regn process executed successfully\n")
    fileUtil.WriteIntoLogFile("InfraGuard.main(). Server Regn process executed successfully")
    fileUtil.WriteIntoLogFile("InfraGuard.main(). Scheduling agent jon on 20 seconds")
    scheduleAgentjob()
   
  }else{
    fileUtil.WriteIntoLogFile(" >>>>>>>>>> InfraGuard.main(). Abort server regn Process. >>>>> ")
    fmt.Printf("Abort server regn Process. Chk log at /var/logs/infraguard/activityLog")
  }
}//main





func scheduleAgentjob(){
  scheduler := gocron.NewScheduler()
  scheduler.Every(freqToHitApi_InSeconds).Seconds().Do(handleUserMgmt)
    <- scheduler.Start()
}
func handleUserMgmt(){
   var nextWork []string
    nextWork = agentUtil.GetNextWork()
    if(nextWork != nil && len(nextWork) > 0){
    fmt.Println("InfraGuard.main(). Length infraGuardResponse = : ",len(nextWork)) 
    ExecuteWork(nextWork)
    }else{
       fmt.Println("InfraGuard.main(). There is no new work") 
    }
  

}//handleUserMgmt

func ExecuteWork(nextWork []string){
  //const delim = ":"
  var pubKey, userName, prefShell, privilege string
  var values []string

  for i := 0; i < len(nextWork); i++{
    values = stringUtil.SplitData(nextWork[i], agentUtil.Delimiter)
    
    if(values[1] == "addUser"){
       
        values = stringUtil.SplitData(nextWork[i+1], agentUtil.Delimiter)
        pubKey = values[1]


        values = stringUtil.SplitData(nextWork[i+2], agentUtil.Delimiter)
        //userName = values[1] + stringUtil.RandStringBytes(3)
        userName = values[1]
        
        values = stringUtil.SplitData(nextWork[i+3], agentUtil.Delimiter)
        prefShell = values[1]

        fmt.Println("userName = : ", userName)
        fmt.Println("Activity Name ----------- Add User -----------------------")
        fmt.Println("pubKey = : ", pubKey)
        fmt.Println("userName = : ", userName)
        fmt.Println("prefShell = : ", prefShell)

        status := userMgmt.AddUser(userName, prefShell, pubKey ) 
        fmt.Println("\n247. AgentUtil.ExecuteWork(). Status Add User = : ",status)

        i += 3

    }

    if(values[1] == "deleteUser"){
        values = stringUtil.SplitData(nextWork[i+1], agentUtil.Delimiter)
        userName = values[1]
        fmt.Println("----------- userName to delete = : ", userName)
        userName = "Vinyl94EC6C"
        status := userMgmt.Userdel(userName, false)
        fmt.Println("142. status deleteUser  = : ", status)
        i += 1
    }

    if(values[1] == "changePrivilege"){
      status := ""
        values = stringUtil.SplitData(nextWork[i+1], agentUtil.Delimiter)
        userName = values[1]

        values = stringUtil.SplitData(nextWork[i+2], agentUtil.Delimiter)
        privilege = values[1]
        
        msg := "privilege = : "+privilege +" >> userName = : "+userName
        fmt.Println(msg)
        if(privilege == "root"){
           status = userMgmt.GiveRootAccess(userName)
        }else{
             status = userMgmt.GiveNormalAccess(userName)
         }
       
        fmt.Println("142. status changePrivilege  = : ", status) //success
        i += 2
    }
    
   
  }

}

