

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

  func SendExecutionStatus(serverUrl string, status string, id, localQryStr string) string{
   serverIp := ExecComand("hostname --all-ip-addresses", "AgentUtil.SendExecutionStatus.go 38")
   serverIp = strings.TrimSpace(serverIp)
  

  qryStr := "?serverIp="+serverIp+"&id="+id

  if(status == "success" || status == "0"){
    qryStr = qryStr + "&status=0"
  }else{
    qryStr = qryStr + "&status=1"
  }
  if(len(localQryStr) > 0){
    localQryStr = "?"+localQryStr
  }
  serverUrl = serverUrl + qryStr+localQryStr
  serverUrl = strings.Replace(serverUrl, "\n","",-1)
 /*
  Sir,
  Due to some personal work, a sum of Rs. 15000/- is needed as advance salary.
  I will repay it in two installment (10000 & 5000) starting from May salary

  I am requesting you, please sanction the same

  Warm regards
  Piyush
 */
  
   // Send execution status [success or fail] 
  
  res, err := http.Get(serverUrl)
  if err != nil {
      fileUtil.WriteIntoLogFile("AgentUtil.sendExecutionStatus() L 57. Error while process this url - serverUrl = : "+serverUrl)
      fileUtil.WriteIntoLogFile("Error at AgentUtil.sendExecutionStatus(). LN 58. Msg = : "+err.Error())
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

 
 /*func SendExecutionStatus(serverUrl string, status string, id string, param ... string) string{
   serverIp := ExecComand("hostname --all-ip-addresses", "AgentUtil.SendExecutionStatus.go 38")
   serverIp = strings.TrimSpace(serverIp)
  

  qryStr := "?serverIp="+serverIp+"&id="+id
  if(len(param) == 1){
    qryStr = qryStr + "&userName="+ param[0]
  }

  if(len(param) == 2){
    qryStr = qryStr + "&pwd="+ param[1]
  }
  
  if(status == "success"){
    qryStr = qryStr + "&status=0"
  }else{
    qryStr = qryStr + "&status=1"
  }

  serverUrl = serverUrl + qryStr
  serverUrl = strings.Replace(serverUrl, "\n","",-1)
 
  
   // Send execution status [success or fail] 
  
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
*/

