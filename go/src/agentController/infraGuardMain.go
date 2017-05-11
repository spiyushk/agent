

package main
// version No 2 dated :- 10-May-2017
import (
    
  
    "stringUtil"
    "serverMgmt"
    "fmt"
    "fileUtil"
    "userMgmt"
    "agentUtil"
    "github.com/jasonlvhit/gocron"  // go get github.com/robfig/cron
  
     //"strings"
    //"strconv"
)

var freqToHitApi_InSeconds uint64 = 20


func main() {
  fmt.Println("InfraGuard.main()") 
  respStr :=serverMgmt.DoServerRegnProcess()
  //respStr="0"
  if(respStr =="0"){
    fileUtil.WriteIntoLogFile("InfraGuard.main(). Server Regn process executed successfully")
    fmt.Printf("\nServer Regn process executed successfully. Agent next job will be fire on every 20 seconds. Waiting \n")
    fileUtil.WriteIntoLogFile("---------- Agent next job will be fire on every 90 seconds. Waiting  -------------\n")
    scheduleAgentjob()
   
  }else{
    fileUtil.WriteIntoLogFile(" >>>>>>>>>> InfraGuard.main(). Abort server regn Process. >>>>> ")
    fmt.Printf("Abort server regn Process. Chk log at /var/logs/infraguard/activityLog")
    return
  }
}//main



func scheduleAgentjob(){
  scheduler := gocron.NewScheduler()
  scheduler.Every(freqToHitApi_InSeconds).Seconds().Do(seekNextWork)
  scheduler.Every(1).Hour().Do(isAlive)
  <- scheduler.Start()
}


func seekNextWork(){
  nextWork := agentUtil.GetNextWork()
  var values [] string
  var cntr int 
  for i := 0; i < len(nextWork); i++{
    values = stringUtil.SplitData(nextWork[i], agentUtil.Delimiter)
    if(values[1] == "addUser" || values[1] == "deleteUser" || 
          values[1] == "changePrivilege" || values[1] == "lockDownServer"){
      cntr = i
      cntr := userMgmt.UserAccountController(values[1], nextWork, cntr) ;
      i = cntr;
    }
  
    if(values[1] == "listEnv"){
     // agentUtil.Send_EnVData()
    }

     if(values[1] == "addEnv"){
        //agentUtil.SetEnvData("piyush","cantata_Key102", "cantata_Value102")
    }
 
  }

}

func isAlive(){
  fileUtil.WriteIntoLogFile(" -------  Infraguard agent code is still running. Next log will be after 1 Hr.... -------")
}









