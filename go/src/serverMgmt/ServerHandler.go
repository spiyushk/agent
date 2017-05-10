

package serverMgmt
// version No 1 dated :- 03-Apr-2017
import (
    "strings"
    "encoding/json"
    "net/http"
    "agentUtil"
    "stringUtil"
    "fileUtil"
  _ "fmt" // for unused variable issue
    "io/ioutil"
    //"fmt"    
)

/*{"fieldCount":0,"affectedRows":1,"insertId":44,
  "serverStatus":2,"warningCount":0,"message":"","protocol41":true,"changedRows":0}*/

const baseUrl = "https://ojf489mkrc.execute-api.us-west-2.amazonaws.com/dev/registerserver"
func DoServerRegnProcess() (string){
    url := baseUrl + getQueryString()
    fileUtil.WriteIntoLogFile("ServerHandler.DoServerRegnProcess(). Going to hit url = : "+url)

    res, err := http.Get(url)
    if err != nil {
        fileUtil.WriteIntoLogFile("Error at ServerHandler.DoServerRegnProcess(). LN 27. Msg = : "+err.Error())
        return "1"
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
      fileUtil.WriteIntoLogFile("Error at ServerHandler.DoServerRegnProcess(). LN. 32. Msg = : "+err.Error())
      return "1"
    }
    var data interface{} 
    err = json.Unmarshal(body, &data)
    if err != nil {
        fileUtil.WriteIntoLogFile("Error at ServerHandler.DoServerRegnProcess(). LN 38. Msg = : "+err.Error())
        return "1"
    }

   infraGuardResp, _ := data.(map[string]interface{})
    
    var affectedRows float64
    affectedRows = -1
          
    if val, ok := infraGuardResp["affectedRows"].(float64); ok {
      affectedRows = val
    } else {
      errorMsg := "ServerHandler.DoServerRegnProcess() LN 59. Unable to cast into float64"
      fileUtil.WriteIntoLogFile(errorMsg)

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
   
   userList := agentUtil.ExecComand("cat /etc/passwd | grep '/home' | cut -d: -f1", "ServerHandler.go 71")
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

 /*
    Read Command line arguments given at the time of agentInstaller.sh execution
 */
 //var sName, pId, licenseKey string
 var sName, pId, licenseKey string
 sName = "notGiven"
 pId = "notGiven"
 licenseKey = ""

 if(fileUtil.IsFileExisted("/tmp/serverInfo.txt")){
  args := stringUtil.SplitData(fileUtil.ReadFile("/tmp/serverInfo.txt", false), ":")
  if(len(args) == 3){
    sName = args[0];
    pId = args[1];
    licenseKey = args[2];
    agentUtil.ExecComand("rm -r /tmp/serverInfo.txt", "ServerHandler.go 106")
  }
 }

 
 cpuDetails := agentUtil.ExecComand("lscpu", "ServerHandler.go 105")
 cpuDetails = stringUtil.FindKey(cpuDetails)

 kernelDetails := agentUtil.ExecComand("cat /etc/*-release", "ServerHandler.go 114")
 kernelDetails = stringUtil.FindKey(kernelDetails)


 qryStr := "?serverName="+sName+"&serverIp="+serverIp+"&hostName="+hostName+"&projectId="+pId+"&userList="+users+"&licenseKey="+licenseKey
 qryStr = strings.Replace(qryStr, "\n","",-1)
 return qryStr
}






