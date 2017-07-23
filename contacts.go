package main

import (
  "fmt"
  "log"
  "net/http"
  "time"
  "os"
  "io/ioutil"
  "bufio"
  "strconv"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

type Contact struct {
  Email       string      `json:"email"`
  Uid         string      `json:"uid"`
  Cookie      string      `json:"cookie"`

}

type Cookie struct {
  value string
}

func main() {
    fmt.Println("-------------------------- Welcome to the Capital Data Application Test! --------------------------")

  // get option from commande line
    // 1 - Call the api using the 1 - 4 option
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("------------------- This application helps you to store your contacts into a database")
    fmt.Println("--- Please select an option :")
    fmt.Println("--- 1 : Screen one level of information and few contacts")
    fmt.Println("--- 2 : Screen two or less information per contacts, few contacts");
    fmt.Println("--- 3 : Screen two or less information per contacts, several contacts");
    fmt.Println("--- 4 : Store your contacts informations and screen contacts with missing informations");

    tryattempt := 0
    mylocaloption := 0
    _ = tryattempt
    for tryattempt < 5 {
      fmt.Print("> ")
      text, _ := reader.ReadString('\n')

      myoption, err := strconv.Atoi(text[0:len(text)-1])

      if err != nil {
        log.Fatal("Error to convert into string :", err)
      }
      if myoption <= 0 || myoption > 4 {
          fmt.Println("Sorry, but please select an option between 1 and 4. Attempt left", 4 - tryattempt)
          tryattempt = tryattempt + 1
      } else {
        mylocaloption = myoption
        break
      }
    }
    if tryattempt == 5 {
      fmt.Println(" I can't do what i am not asked to. I'll stop while you think again about your options. Restart me when you've made a valid choice")
      return
    }

    if mylocaloption < 5 {
      // get data from api
      url := fmt.Sprintf("%s%d", "http://tech.kdata.fr:8080/alexandre.carle/", mylocaloption);

      // build the request
      req, err := http.NewRequest(http.MethodGet, url, nil);
      if (err != nil){
        log.Fatal("NewRequest :", err)
        return
      }

      client := &http.Client{
        Timeout: time.Second * 2,
      }

      resp, err := client.Do(req)
      if err != nil {
        log.Fatal("Do : ", err)
        return
      }

      defer resp.Body.Close()

      if mylocaloption <= 3 {
        // read data
        responseData, err := ioutil.ReadAll(resp.Body)
        if  err != nil {
          log.Fatal(err)
        } else {
          fmt.Println(string(responseData))
        }
      } else if mylocaloption == 4 {
        // connect to do
        db, err := sql.Open("sqlite3", "./contacts.db")
        checkErr(err)

        reader := bufio.NewReader(resp.Body)
        var line string
        var cookie []Cookie
        var uid string
        var email string
        var tmpParseString string
        tmpParseKey := ""
        hasUid := false
        for {
          line, err = reader.ReadString('\n')

          b := []rune(line);
          // parseState defines type of data on retrieve
          // 1 is uid
          // 2 is email
          // 3 is cookie
          parseState := 0
          for i := 0; i < len(b); i++ {
            switch b[i] {
            case '@': // Potential end of sequence
              if parseState == 3 {
                cookie = append(cookie, Cookie{tmpParseString})
              } else if parseState == 1 {
                uid = ""+tmpParseString
                hasUid = true
              }else if parseState == 2 {
                tmpParseString = tmpParseString + "@"
              }
              parseState = 2
              break
            case 'u':
              if hasUid {
                tmpParseString = tmpParseString + "u"
              } else {
                tmpParseKey = "u"
              }
              break
            case 'i':
              if tmpParseKey == "cook" {
                tmpParseKey = tmpParseKey + "i"
              } else if hasUid {
                tmpParseString = tmpParseString + "i"
              } else if tmpParseKey == "u" {
                tmpParseKey = tmpParseKey + "i"
              } else {
                tmpParseString = tmpParseString + "i"
              }

              break
            case 'd': // Potential end of sequence
              if hasUid {
                tmpParseString = tmpParseString + tmpParseKey + "d"
              } else if tmpParseKey == "ui" {
                if parseState == 3 {
                  cookie = append(cookie, Cookie{tmpParseString})
                  tmpParseString = ""
                } else if parseState == 2 {
                  email = ""+tmpParseString
                  tmpParseString = ""
                }
                parseState = 1
              }
              tmpParseKey = ""
              break
            case 'c':
              tmpParseKey = "c"
              break
            case 'o':
              if tmpParseKey == "c" || tmpParseKey == "co" {
                tmpParseKey = tmpParseKey + "o"
              } else {
                tmpParseString = tmpParseString + tmpParseKey + "c"
                tmpParseKey = ""
              }
              break
            case 'k':
              if tmpParseKey == "coo"{
                tmpParseKey = "cook"
              } else {
                tmpParseString = tmpParseString + tmpParseKey + "k"
                tmpParseKey = ""
              }
              break
            case 'e': // Potential end of sequence
              if tmpParseKey == "cooki" {
                if parseState == 1 {
                  uid = tmpParseString + ""
                  hasUid = true
                } else if parseState == 2 {
                  email = tmpParseString + ""
                } else if parseState == 3 {
                  cookie = append(cookie, Cookie{tmpParseString})
                }
                tmpParseString = ""
                tmpParseKey = ""
                parseState = 3
              } else {
                  tmpParseString = tmpParseString + "e"

              }
              tmpParseKey = ""
              break
            case '\n':
              if parseState == 1 {
                uid = tmpParseString + ""
              } else if parseState == 2 {
                email = tmpParseString + ""
              } else if parseState == 3 {
                cookie = append(cookie, Cookie{tmpParseString})
              }
              tmpParseString = ""
              hasUid = false
              break
            default:
              tmpParseString = tmpParseString + tmpParseKey + string(b[i])
              tmpParseKey = ""
            }
          }

          if err != nil {
            break
          }

          // if email exists
          if len(email) > 0 {
            // insert

            // if no cookies
            if len(cookie) == 0 {
              stmt, err := db.Prepare("INSERT INTO contacts(uid, email) values(?,?)")
              checkErr(err)

              res, err := stmt.Exec(uid, email)
              checkErr(err)

              id, err := res.LastInsertId()
              checkErr(err)

              _ = id
            } else {
              // insert for each cookie
              for _, acookie := range cookie {
                stmt, err := db.Prepare("INSERT INTO contacts(uid, email, cookie) values(?,?,?)")
                checkErr(err)

                res, err := stmt.Exec(uid, email, acookie.value)
                checkErr(err)

                id, err := res.LastInsertId()
                checkErr(err)

                _ = id
              }
            }
          }

          cookie = cookie[:0]
          email = ""
          uid = ""

        }
      }

    }


    // test data format

    // and store the data into sqlite if needed

    // screen the result of the storage and the error of storage
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
