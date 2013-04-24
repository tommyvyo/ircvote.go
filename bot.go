package main

import (
  "fmt"
  "bufio"
  "log"
  "net"
  "net/textproto"
  "strings"
)

type IRCBot struct {
  Server      string
  Port        string
  User        string
  Nick        string
  Channel     string
  Pass        string
  Connection  net.Conn
}

//create a new bot with the given server info
//currently static but will create prompt for any server
func createBot() *IRCBot {
  //serverPrompt()
  return &IRCBot {
    Server:     "irc.freenode.net",
    Port:       "6665",
    Nick:       "MachineDelVote",
    Channel:    "##botVoteTesting",
    Pass:       "",
    Connection: nil,
    User:       "VoteBot",
    }
}

/*On creation of a bot prompts the user for correct server information
func severPrompt() {
}*/

//connects the bot to the server
func (bot *IRCBot) ServerConnect() (connection net.Conn, err error) {
  connection, err = net.Dial("tcp", bot.Server + ":" + bot.Port)
  if err != nil {
    log.Fatal("Unable to connect to the specified server", err)
  }
  bot.Connection = connection
  log.Printf("Successfully connected to %s ($s)\n", bot.Server, bot.Connection.RemoteAddr())
  return bot.Connection, nil
}

func main() {
  votes := make(map[string]int)
  bot := createBot()
  connection, _ := bot.ServerConnect()
  fmt.Fprintf(connection, "USER %s 8 * :%s\n", bot.Nick, bot.Nick)
  fmt.Fprintf(connection, "NICK %s\n", bot.Nick)
  fmt.Fprintf(connection, "JOIN %s\n", bot.Channel)
  defer connection.Close()

  reader := bufio.NewReader(connection)
  respReq := textproto.NewReader(reader)
  for {
    line, err := respReq.ReadLine()
    if err != nil {
      break
    }
    words := strings.Split(line, " ")
    if strings.Contains(line, "PING") {
      fmt.Printf("\033[92mPONG\n")
      fmt.Fprintf(connection, "PONG\n")
    }
    if strings.Contains(line, bot.Channel) && strings.Contains(line, "!voteUP") {
      commandIndex := indexOf("!voteUP", words)
      if commandIndex != -1 {
        votes[words[commandIndex + 1]] += 1
        fmt.Fprintf(connection, "PRIVMSG %s Upvoted: %s\n", bot.Channel, words[commandIndex + 1])
        fmt.Printf("\033[92mPRIVMSG %s Upvoted: %s\n", bot.Channel, words[commandIndex +1])
      }
    }
    fmt.Printf("\033[93m%s\n", line)
  }
}

func indexOf(value string, mySlice []string) int {
  for index, curValue := range mySlice {
    if curValue == value {
      return index
    }
  }
  return -1
}
