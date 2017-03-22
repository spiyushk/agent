

package serverMgmt

import (
    "strings"
    "encoding/json"
    "net/http"
    "agentUtil"
    "stringUtil"
    "fileUtil"
  _ "fmt" // for unused variable issue
    "io/ioutil"
)

/*{"fieldCount":0,"affectedRows":1,"insertId":44,
  "serverStatus":2,"warningCount":0,"message":"","protocol41":true,"changedRows":0}*/

const url2 = "https://ojf489mkrc.execute-api.us-west-2.amazonaws.com/dev/registerserver?serverName=demoServer206&serverIp=12.12.12.12&hostName=linuxAMI&projectId=5&userList=aaa,bbb,ccc"

const baseUrl = "https://ojf489mkrc.execute-api.us-west-2.amazonaws.com/dev/registerserver"


func DoServerRegnProcess() (string){
   
    url := baseUrl + getQueryString()
    fileUtil.WriteIntoLogFile("ServerHandler.DoServerRegnProcess(). Going to hit url = : "+url, "")

    res, err := http.Get(url)
    if err != nil {
        fileUtil.WriteIntoLogFile("Error at ServerHandler.DoServerRegnProcess(). LN 35. Msg = : "+err.Error(), "")
        return "1"
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
      fileUtil.WriteIntoLogFile("Error at ServerHandler.DoServerRegnProcess(). LN. 40. Msg = : "+err.Error(), "")
      return "1"
    }
    var data interface{} 
    err = json.Unmarshal(body, &data)
    if err != nil {
        fileUtil.WriteIntoLogFile("Error at ServerHandler.DoServerRegnProcess(). LN 46. Msg = : "+err.Error(), "")
        return "1"
    }


    
// http://stackoverflow.com/questions/14289256/cannot-convert-data-type-interface-to-type-string-need-type-assertion
   infraGuardResp, _ := data.(map[string]interface{})
    
    var affectedRows float64
    affectedRows = -1
          
    if val, ok := infraGuardResp["affectedRows"].(float64); ok {
      affectedRows = val
    } else {
      errorMsg := "ServerHandler.DoServerRegnProcess() LN 59. Unable to cast into float64"
      fileUtil.WriteIntoLogFile(errorMsg, "")

    }
   if(affectedRows > 0){
      return "0"
    }else{
     return "1"
    }

}



func getQueryString()(string){
   serverIp := agentUtil.ExecComand("hostname --all-ip-addresses", "ServerHandler.go 74")
   hostName := agentUtil.ExecComand("hostname", "ServerHandler.go 75")
    
   serverIp = strings.TrimSpace(serverIp)
   hostName = strings.TrimSpace(hostName)
   
   userList := agentUtil.ExecComand("cat /etc/passwd | grep '/home' | cut -d: -f1", "ServerHandler.go 84")
   userList2 := strings.Split(userList,"\n")

  
   max := 5
   if(len(userList2) < max){
    max = len(userList2)
   }

  users := ""
  for i := 0; i <  max; i++ {
    if(len(users) ==0){
      users = userList2[i]
    }else{
      users = users +","+userList2[i]
    }
 
  }
 users = strings.TrimSpace(users)

 cpuDetails := agentUtil.ExecComand("lscpu", "ServerHandler.go 105")
 cpuDetails = stringUtil.FindKey(cpuDetails)

 kernelDetails := agentUtil.ExecComand("cat /etc/*-release", "ServerHandler.go 108")
 kernelDetails = stringUtil.FindKey(kernelDetails)
  
 //qryStr := "?serverName=demoServer3&serverIp="+serverIp+"&hostName="+hostName+"&projectId=5&users="+users +"&cpuDetails="+cpuDetails+"&kernelDetails="+kernelDetails;
 qryStr := "?serverName=testingServer507&serverIp="+serverIp+"&hostName="+hostName+"&projectId=5&userList="+users
 qryStr = strings.Replace(qryStr, "\n","",-1)
 return qryStr
}





