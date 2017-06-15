

package main
// version No 2 dated :- 10-May-2017
import (
    
  
    "stringUtil"
    //"serverMgmt"
    "fmt"
    "fileUtil"
    "userMgmt"
    "agentUtil"
   // "github.com/jasonlvhit/gocron"  // go get github.com/robfig/cron
  
    // "strings"
    //"strconv"
)

var freqToHitApi_InSeconds uint64 = 20


func main() {
 fmt.Println(" ----------- InfraGuard.main(). ----------------- ")  
 fileUtil.WriteIntoLogFile(" ----------- InfraGuard.main(). ----------------- ")


 
  status := ""
  msg := ""

   if(fileUtil.IsFileExisted("/tmp/serverInfo.txt")){
  args := stringUtil.SplitData(fileUtil.ReadFile("/tmp/serverInfo.txt", false), ":")
  if(len(args) == 3){
    msg = "\n--------------- sName = : "+args[0]+ " >> pId = : "+args[1] + " >> LKey = : "+args[2]
    fileUtil.WriteIntoLogFile(msg)
    agentUtil.ExecComand("rm -r /tmp/serverInfo.txt", "ServerHandler.go 106")
  }else{
    fileUtil.WriteIntoLogFile("\n-------------- unable to read /tmp/serverInfo.txt")
  }
 }


  usrLoginName := "test55"
  preferredShell := ""
  pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDbbt5OuNlKkd3GWkO78eetw+fqzntLhPXt0QUWCY1CqqerZKCG1AKNxh1+SxYk/iUcsiqqQ1GXHjpVqTi0dx/EaTgs2H3WS7lrsr3hLPEJDqaEfxs2m+DFzqYV/PLExaH7Bc6PNTjz39g/MNZyScMw6XAhNoivnaGmlsWmxT2rezXrmAV1dDqR56W4jg7X9PchrMvAo80N7Bo8P6DrpKHEi2eQSH6keZ4IQS/A8730FDdn4WmmOkHSq375y57w5H1jM5zC7ZL2h4YY5BHxcE3/SBgiIMv/6NyqJDPZ9Vf8Ff8cxM7O402fe4yChvzTijlm+yT+VvBkjz7jwPn6p+PH test50@piyush-ubuntu"
  

  status = userMgmt.AddUser(usrLoginName, preferredShell, pubKey);
  msg = "Final Status of AddUser(). For user  = : "+usrLoginName +" >> status = : "+status
  fileUtil.WriteIntoLogFile(msg)
  fmt.Println("\n\n",msg)

  

 /*
  status = userMgmt.ProcessToChangePrivilege(usrLoginName, "root")
  fmt.Println("Final Status of  ProcessToChangePrivilege = : ", status)

  msg = "Final Status of ProcessToChangePrivilege(). For user  = : "+usrLoginName+" & priv = root. >> status = : "+status
  fileUtil.WriteIntoLogFile(msg)
  fmt.Println(msg)


  status = userMgmt.ProcessToChangePrivilege("test2", "normal")
  fmt.Println("Final Status of  ProcessToChangePrivilege = : ", status)

  msg = "Final Status of ProcessToChangePrivilege(). For user  test2  & priv = normal. >> status = : "+status
  fileUtil.WriteIntoLogFile(msg)
  fmt.Println("\n\n",msg)

*/


  /*usersToLock := []string{"test2", "test51", "test52", "test855"}
  status = userMgmt.ProcessToLockDownServer(usersToLock)

  msg = "\n\nFinal Status of ProcessToLockDownServer(). users are test51 & test52, test855 >> status = : "+status
  fileUtil.WriteIntoLogFile(msg)
  fmt.Println("\n\n",msg)



  usersToUnlock := []string{"test52"}
  status = userMgmt.ProcessToUnlockServer(usersToUnlock)

  msg = "Final Status of ProcessTo Unlock Server() for user test52. >> status = : "+status
  fileUtil.WriteIntoLogFile(msg)
  fmt.Println("\n\n",msg)
*/




  

 
}//main


