

package main
// version No 1 dated :- 03-Apr-2017
import (
    
  
    "stringUtil"
    "serverMgmt"
    "fmt"
    "fileUtil"
    "userMgmt"
    "agentUtil"
    "github.com/jasonlvhit/gocron"  // go get github.com/robfig/cron 
     "strings"
    //"strconv"
)
var freqToHitApi_InSeconds uint64 = 20

/*func main() {
    nextWork := agentUtil.GetNextWork()
    if(nextWork != nil){
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
    fmt.Printf("---------- Agent next job will be fire on every 20 seconds. Waiting  -------------")
    fileUtil.WriteIntoLogFile("---------- Agent next job will be fire on every 20 seconds. Waiting  -------------")
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
     ExecuteWork(nextWork)
    }else{
       fmt.Println("InfraGuard.main(). There is no new work") 
    }
  

}//handleUserMgmt

func ExecuteWork(nextWork []string){
  //const delim = ":"
  var pubKey, userName, prefShell, privilege, id string
  var values []string

  for i := 0; i < len(nextWork); i++{
    values = stringUtil.SplitData(nextWork[i], agentUtil.Delimiter)

    fmt.Println("----------------------------------------------------------------------------------\n")
    if(values[1] == "addUser"){
        serverUrl := "https://spjuv2c0ae.execute-api.us-west-2.amazonaws.com/dev/addeduserbyagent"
        values = stringUtil.SplitData(nextWork[i+1], agentUtil.Delimiter)
        pubKey = values[1]


        values = stringUtil.SplitData(nextWork[i+2], agentUtil.Delimiter)
        //userName = values[1] + stringUtil.RandStringBytes(3)
        userName = values[1]
        
        values = stringUtil.SplitData(nextWork[i+3], agentUtil.Delimiter)
        prefShell = values[1]

        values = stringUtil.SplitData(nextWork[i+4], agentUtil.Delimiter)
        id = values[1]

      
        msg :=  "Going to add userName = : "+userName
        fmt.Println(msg)
        fileUtil.WriteIntoLogFile(msg)
        status := userMgmt.AddUser(userName, prefShell, pubKey ) 
       
        sendExecutionStatus(serverUrl, status , id, userName)
      
        i += 4

    }

    if(values[1] == "deleteUser"){
        serverUrl := "https://vglxmaiux1.execute-api.us-west-2.amazonaws.com/dev/deleteduserbyagent"
        values = stringUtil.SplitData(nextWork[i+1], agentUtil.Delimiter)
        userName = values[1]

        values = stringUtil.SplitData(nextWork[i+2], agentUtil.Delimiter)
        id = values[1]

        msg :=  "Going to delete userName = : "+userName
        fmt.Println(msg)
        fileUtil.WriteIntoLogFile(msg)

        status := userMgmt.Userdel(userName, false)
        fmt.Println("status deleteUser  = : ", status)
        sendExecutionStatus(serverUrl, status , id, userName)
        i += 2
    }

    if(values[1] == "changePrivilege"){
      serverUrl := "https://a1gpcq76u3.execute-api.us-west-2.amazonaws.com/dev/privilegechangedbyagent"
      status := ""
        values = stringUtil.SplitData(nextWork[i+1], agentUtil.Delimiter)
        userName = values[1]

        values = stringUtil.SplitData(nextWork[i+2], agentUtil.Delimiter)
        privilege = values[1]

        values = stringUtil.SplitData(nextWork[i+3], agentUtil.Delimiter)
        id = values[1]

        
        msg :=  "Going to change privilege for userName = : "+userName+ " Priv = : "+privilege
        fmt.Println(msg)
        fileUtil.WriteIntoLogFile(msg)

        status = userMgmt.ProcessToChangePrivilege(userName, privilege)


      
       /* if(privilege == "root"){
           status = userMgmt.GiveRootAccess(userName)
        }else{
             status = userMgmt.GiveNormalAccess(userName)
         }
*/
        sendExecutionStatus(serverUrl, status , id, userName) 
       
        fmt.Println("status changePrivilege  = : ", status) 
        i += 3
    }
    fmt.Println("----------------------------------------------------------------------------------\n")
   
  }

}


func sendExecutionStatus(serverUrl string, status string, id string, param ... string) string{

   serverIp := agentUtil.ExecComand("hostname --all-ip-addresses", "ServerHandler.go 74")
   serverIp = strings.TrimSpace(serverIp)
  


  qryStr := "?serverIp="+serverIp+"&id="+id
  if(len(param) >= 1){
    qryStr = qryStr + "&userName="+ param[0]
  }
  
  if(status == "success"){
    qryStr = qryStr + "&status=0"
  }else{
    qryStr = qryStr + "&status=1"
  }

  serverUrl = serverUrl + qryStr
  serverUrl = strings.Replace(serverUrl, "\n","",-1)
  status = agentUtil.HitAnyUrl(serverUrl)
  return status

}//sendExecutionStatus
