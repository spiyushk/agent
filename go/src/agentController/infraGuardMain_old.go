

package main
// version No 2 dated :- 10-May-2017
import (
    
  
   // "stringUtil"
    //"serverMgmt"
    "fmt"
    "fileUtil"
    "userMgmt"
  //  "agentUtil"
   // "github.com/jasonlvhit/gocron"  // go get github.com/robfig/cron
  
     //"strings"
    //"strconv"
)

var freqToHitApi_InSeconds uint64 = 20


func main() {
 fmt.Println(" ----------- InfraGuard.main(). ----------------- ")  
  usrLoginName := "test8"
  preferredShell := ""
  pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCXpt9zMgAnK8uLHhxRdW4H4ii2yTYw1SIEG4oR89SogncsVSdm2N+blu+9VyOVq93Fy/825EVyrwV7/leQuIKYMxO6sXOx9BDhRIKFff50dsJZZ+hGIF48N7c+EeV42rO87xBx6DOnixNLaEyaRYddM+rKo03RFRNtKZTnheYnrk+lBFoYMIP5VuO7vxzzoK88Kt1mb7LJ9Jg420bV7QFGFwdDGs3He5EfM8jxxi9XLoK5AG4X28o3uRRdUJOC0DoUMbVdKRczlv0Q7RvRM14VPnj+abvdrqt6zw6ieJpKjHclYx3kZoVg3G9Z5I90rnQmIcqcdb7YKa4DM4uLS8FD test@InfraGuard"
 

  status := userMgmt.AddUser(usrLoginName, preferredShell, pubKey);
  
  msg := "Final Status of AddUser(). For user  = : "+usrLoginName +" >> status = : "+status
  fileUtil.WriteIntoLogFile(msg)
  fmt.Println("\n\n",msg)


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


  usersToLock := []string{"test2", "test4", "test10", "test5"}
  status = userMgmt.ProcessToLockDownServer(usersToLock)

  msg = "Final Status of ProcessToLockDownServer(). >> status = : "+status
  fileUtil.WriteIntoLogFile(msg)
  fmt.Println("\n\n",msg)


  

 
}//main


