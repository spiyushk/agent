package main

import (
     
    "agentUtil"
    "fmt"
    //"strconv"
  
   // "io/ioutil"
    "fileUtil"

    //"stringUtil"



     
    
  //  "os"
    _ "fmt" // for unused variable issue
 
    
    "strings"
   // "bufio"


)

func main() {
    status := processToChangePrivilege("piyushsinha","root")
    fmt.Println("23. ---------------- status processToChangePrivilege = : ", status)
}

func processToChangePrivilege(usrName, privType string) string{
    accessRight := usrName+"   ALL=(ALL:ALL) ALL"
    
    
    // Read /etc/sudoers file for user right, if any
    //cmd := "sudo awk '/"+usrName+"/ {print}' /etc/sudoers"
     cmd := "sudo awk '/"+usrName+"/ {print}' /tmp/sudoers.bak"
    oldPriv := agentUtil.ExecComand(cmd, "misc. L78")
    fmt.Println("79 oldPriv : = ",oldPriv)
    if(oldPriv == "success"){
        oldPriv = ""
    }
    

    tmpFilePath := "/tmp/sudoers.bak"
    status := ""
  
    // Create  a back up copy of /etc/sudoers file
    status = agentUtil.ExecComand("sudo cp /etc/sudoers "+tmpFilePath, "misc. L38")
    fmt.Println("status at 26.  = : ", status)

    status = agentUtil.ExecComand("sudo chmod 777 "+tmpFilePath, "misc. L41")
    fmt.Println("status at 29.  = : ", status)



    if( privType == "root"){
        if(len(oldPriv) == 0){
                    
            fileUtil.WriteIntoFile(tmpFilePath, accessRight, true, false)
            fmt.Println("condition matched at 52.  = : ")
            // Replace the sudoers file with the tmpFilePath 
            status = agentUtil.ExecComand("sudo cp "+tmpFilePath +" /etc/sudoers", "misc. L38")
            
        }else{
             // If old priv is commented, then remove such comment
            if(strings.Contains(oldPriv, "#")){
                status = fileUtil.ReplaceLineOrLinesIntoFile(tmpFilePath, oldPriv, accessRight)
               
                // Replace the sudoers file with the tmpFilePath 
                status = agentUtil.ExecComand("sudo cp "+tmpFilePath +" /etc/sudoers", "misc. L38")
            
            }   
        }
    }

    // If user's old priv is root level access, then comment access right data to become a normal user
    if( privType == "user"){
        if(len(oldPriv) > 0){
            status = fileUtil.ReplaceLineOrLinesIntoFile(tmpFilePath, oldPriv, "#"+accessRight)
            // Replace the sudoers file with the tmpFilePath 
            status = agentUtil.ExecComand("sudo cp "+tmpFilePath +" /etc/sudoers", "misc. L38")
            
        }
    }
    return status
}
