

package agentUtil

import (
    
    "os/exec"
  _ "fmt" // for unused variable issue
    "fileUtil"
     "fmt"
    "io/ioutil"
     "encoding/json"
     "net/http"
     "stringUtil"
      "strconv"
      "strings"
   
)

const url = "https://8zq40bjoqc.execute-api.us-west-2.amazonaws.com/dev/hello"
var maxSize int = 100
var array = make([]string, maxSize) 
var cntr int = 0
const Delimiter = ":"


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

func hitApi() string{
  
    array = make([]string, maxSize)
    cntr = 0

    fileUtil.WriteIntoLogFile("UserHandler.getNextWork(). Going to hit url = : "+url)
    res, err := http.Get(url)
    if err != nil {
        fileUtil.WriteIntoLogFile("Error at UserHandler.getNextWork(). LN 119. Msg = : "+err.Error())
        return "1"
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
      fileUtil.WriteIntoLogFile("Error at UserHandler.getNextWork(). LN. 124. Msg = : "+err.Error())
      return "1"
    }
     
    var data []interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        fileUtil.WriteIntoLogFile("Error at UserHandler.getNextWork(). LN 130. Msg = : "+err.Error())
        return "1"
    }
    //fmt.Println(data)
    parseArray(data)
    TrimArrayToActualSize()
    return "0"
    
   
}//hitApi



func GetNextWork() ([]string){
  for i := 0; i < 10; i++ {
    //var infraGuardResponse []string
    respStatus := hitApi();
    if(respStatus == "0"){
      isValidData := ValidateArray()
      msg := "\n************ isValidData = : "+strconv.FormatBool(isValidData) +" ***********************\n"
      fmt.Println(msg)
      fileUtil.WriteIntoLogFile(msg)
      if(isValidData){
         if(len(array) > 0 ){
             printArray()
             var tmp = array
             array = nil
             return tmp 
         }else{
          return nil;
         }
        
      }
    }
    
  }
  return nil
}

//http://stackoverflow.com/questions/29366038/looping-iterate-over-the-second-level-nested-json-in-go-lang


func parseMap(aMap map[string]interface{}) {
    for key, val := range aMap {
        switch concreteVal := val.(type) {
        case map[string]interface{}:
           // initializeArray(key, "")
            parseMap(val.(map[string]interface{}))
        case []interface{}:
            //initializeArray(key, "")
            parseArray(val.([]interface{}))
        default:
            if val, ok := concreteVal.(string); ok {
              initializeArray(key, val)
            }

           
        }
    }
}

func parseArray(anArray []interface{}) {
    for _, val := range anArray {
        switch concreteVal := val.(type) {
        case map[string]interface{}:
            parseMap(val.(map[string]interface{}))
        case []interface{}:
            parseArray(val.([]interface{}))
        default:
            if val, ok := concreteVal.(string); ok {
               initializeArray(val, "")
            }

        }
    }
}

func initializeArray(key, val string){
  key = stringUtil.RemoveSpace(key)
  val = stringUtil.RemoveSpace(val)
  if(key == "reqData" || val == "reqData" || len(key) == 0 || len(val) == 0){
    return
  }
  array[cntr] = key + Delimiter + val
  cntr = cntr + 1
  return
}

func printArray(){
   for i := 0; i < len(array); i++ {
    val := array[i]
    fmt.Println(val)
    fileUtil.WriteIntoLogFile(val)
  }
}



func ValidateArray() bool{
  
  var values []string
   for i := 0; i < len(array); i++ {
   values = stringUtil.SplitData(array[i], Delimiter)
     
  
   if(values[0] == "activityName"){
      isValidData := checkDataSequence(values[1], i)
      if(isValidData == false){
        return false;
      }
   }

   
  }

  isValidData := checkActivityNameSequence()
  return isValidData
}



func checkDataSequence(activityName string, cnt int) bool{
    var values, sequnce []string
  
    if(activityName == "addUser"){
      sequnce = []string{"activityName","publicKey", "userName", "shell"}
    }

    if(activityName == "deleteUser"){
      sequnce = []string{"activityName","userName"}
    }
    if(activityName == "changePrivilege"){
      sequnce = []string{"activityName","username", "privilege"}
    }
    
    if((cnt + len(sequnce)) > len(array)){
      return false
    }
   

     for i := 0; i < len(sequnce); i++{
        values = stringUtil.SplitData(array[cnt], Delimiter)
        cnt++
  
        if(strings.ToLower(values[0]) != strings.ToLower(sequnce[i])){
           return false
        }
      
      }
    return true 
  }//checkDataSequence


func checkActivityNameSequence() bool{
  var values []string
  var i, cnt int
  cnt = 0
  for i = 0; i < len(array); i++ {
   values = stringUtil.SplitData(array[i], Delimiter)
   if(values[0] == "activityName"){
     cnt++;
   }

  }
  var tmpArray = make([]string, cnt)
  cnt = 0
  
  for i = 0; i < len(array); i++ {
   values = stringUtil.SplitData(array[i], Delimiter)
   if(values[0] == "activityName"){
     tmpArray[cnt] = strconv.Itoa(i) + ":"+values[1]
     cnt++;
   }

  }


  for j := 0; j < len(tmpArray)-1; j++ {
     fmt.Println("288.  tmpArray[i] = : "+tmpArray[j])
     values = stringUtil.SplitData(tmpArray[j], Delimiter)
      fmt.Println("290.  values[0] ",values[0])
      fmt.Println("291.  values[1] ",values[1])

     val1, _ := strconv.ParseInt(values[0], 10, 0)
     activityName1 := values[1]
     
     var jInt64 int64
     jInt64 = int64(j)
     if(jInt64 == 0 && jInt64 != val1){
      return false
     }

     values = stringUtil.SplitData(tmpArray[j+1], Delimiter)
     val2, _ := strconv.ParseInt(values[0], 10, 0)
    

     if(activityName1 == "addUser"){
       if((val1 + 4 )!= val2){
        return false
       }
     }

     if(activityName1 == "deleteUser"){
       if((val1 + 2 )!= val2){
        return false
       }
     }
     if(activityName1 == "changePrivilege"){
       if((val1 + 3 )!= val2){
        return false
       }
     }
  }//I
 return true



}

func TrimArrayToActualSize(){
  var cnt int = 0 
  for i := 0; i < len(array); i++ {
  val := array[i]
    if(len(val) > 0){
      cnt++
    }
  }
  var tmp = array[0:cnt]
  array = tmp
  tmp = nil
 
  
}

