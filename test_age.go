package main
import (
    "fmt"
    "time"
)
func main() {
    now := time.Now()
    dob, _ := time.Parse("2006-01-02", "1985-12-25")
    
    age := now.Year() - dob.Year()
    fmt.Printf("Today: %s\n", now.Format("2006-01-02"))
    fmt.Printf("DOB: %s\n", dob.Format("2006-01-02"))
    fmt.Printf("Year difference: %d\n", age)
    fmt.Printf("Now month: %s, DOB month: %s\n", now.Month(), dob.Month())
    fmt.Printf("Now day: %d, DOB day: %d\n", now.Day(), dob.Day())
    
    if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
        age--
        fmt.Println("Birthday hasn't happened yet, subtracting 1")
    }
    
    fmt.Printf("Final age: %d\n", age)
}
