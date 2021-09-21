package tool

import (
	"fmt"
	"log"
)

var (
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
)

//func init() {
// If the file doesn't exist, create it or append to the file
//file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
//if err != nil {
//	log.Fatal(err)
//}

//log.SetOutput(file)

//infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
//warningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
//errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
//}

func LogWarning(message string) {
	log.Println(fmt.Sprintf("WARNING: %s", message))
	//warningLogger.Println(message)
}

func LogInfo(message string) {
	log.Println(fmt.Sprintf("INFO: %s", message))
	//infoLogger.Println(message)
}

func LogError(message string) {
	log.Println(fmt.Sprintf("ERROR: %s", message))
	//errorLogger.Println(message)
}
