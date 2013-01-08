package main

import (
  "fmt"
  "bufio"
  "log"
  "net"
  "net/textproto"
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
    Channel:    "#orderdeck",
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
    fmt.Printf("%s\n", line)
  }
}

