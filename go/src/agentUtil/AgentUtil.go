

// version No 1 dated :- 03-Apr-2017
package agentUtil
// version No 1 dated :- 03-Apr-2017
import (
    "os/exec"
  _ "fmt" // for unused variable issue
    "fileUtil"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
   
)

func ExecComand(cmd, fromFile string) string {
    
    cmdStatus,err := exec.Command("bash","-c",cmd).Output()
    execStatus := "success"
    if err != nil {
        errorMsg := "Cmd executed = : "+cmd +" : execStatus = : fail. fromFile. = :"+fromFile
        fileUtil.WriteIntoLogFile(errorMsg)
        execStatus = "fail"
        fmt.Println("34. AgentUtil.ExecComand()  errorMsg = : ", errorMsg)
    }

    if (len(string(cmdStatus)) > 0){
        execStatus =  string(cmdStatus)  
    }
    return execStatus
}
   
 func SendExecutionStatus(serverUrl string, status string, id string, param ... string) string{
   serverIp := ExecComand("hostname --all-ip-addresses", "AgentUtil.SendExecutionStatus.go 38")
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
 
  /*
    Send execution status [success or fail] 
  */
  res, err := http.Get(serverUrl)
  if err != nil {
      fileUtil.WriteIntoLogFile("Error at AgentUtil.sendExecutionStatus(). LN 61. Msg = : "+err.Error())
      status =  "1"
  }
  _, error := ioutil.ReadAll(res.Body)
  if error != nil {
    fileUtil.WriteIntoLogFile("Error at AgentUtil.sendExecutionStatus(). LN 66. Msg = : "+error.Error())
    status =  "1"
  }

  fileUtil.WriteIntoLogFile("Successfully sent execution status to this url = : "+serverUrl)
  fmt.Println("Successfully sent execution status to this url = : ",serverUrl) //success
  status =  "0"


  return status

}//sendExecutionStatus


