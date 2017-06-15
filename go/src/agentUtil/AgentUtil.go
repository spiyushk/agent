

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
    "stringUtil"
   
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
   localQryStr = strings.TrimSpace(localQryStr)
  

  qryStr := "?serverIp="+serverIp+"&id="+id

  if(status == "success" || status == "0"){
    qryStr = qryStr + "&status=0"
  }else{
    qryStr = qryStr + "&status=1"
  }
  if(len(localQryStr) > 0){
    localQryStr = "&"+localQryStr
  }
  serverUrl = serverUrl + qryStr+localQryStr
  serverUrl = strings.Replace(serverUrl, "\n","",-1)
   
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

 
func ReadPropertyFile() map[string]string {
  var values, rows []string
  var propertyMap map[string]string
  propertyMap = make(map[string]string)
 
  data := fileUtil.ReadFile(propertyFilePath, false) 
  data = strings.Replace( data, "\"","",-1)  // Remove dbl quotes

  rows = stringUtil.SplitData(data, "\n")
  for _, row := range rows {
    row = strings.TrimSpace(row)
    
    // Ignore row which starts with letter '#'
    if(strings.HasPrefix(row, "#")){
      continue
    }
    if(strings.Contains(row, "=")){
       values = stringUtil.SplitData(row, "=")
       if(len(values) == 2){
         propertyMap[values[0]] = values[1] 
       }
       
    }
  }


  if(propertyMap != nil && len(propertyMap) > 0){
    return propertyMap  
  }
  return nil
  
}


func GetValueFromPropertyMap(propMap map[string]string, key string) string{
  if(len(key) >0){
    if(len(propMap[key]) >0 ){
      return propMap[key]
    }  
  }
  return ""
}

