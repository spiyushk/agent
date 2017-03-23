

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
    //"stringUtil"
    "serverMgmt"
    "fmt"
    "fileUtil"
)
func main() {

  fmt.Println("InfraGuard.main()") 
  respStr :=serverMgmt.DoServerRegnProcess()
  if(respStr =="0"){
    fmt.Printf("\nServer Regn process executed successfully\n")
    fileUtil.WriteIntoLogFile("InfraGuard.main(). Server Regn process executed successfully")

  }else{
    fileUtil.WriteIntoLogFile(" >>>>>>>>>> InfraGuard.main(). Abort server regn Process. >>>>> ")
    fmt.Printf("Abort server regn Process. Chk log at /var/logs/infraguard/activityLog")
  }
  
  
}

