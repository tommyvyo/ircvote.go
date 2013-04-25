package main

import (
  "fmt"
  "bufio"
  "log"
  "net"
  "net/textproto"
  "strings"
)

var NO_USER_ERROR string = "!No user specified!"
var USER_NOT_FOUND_ERROR string = "!User not in channel!"
var CMD_VOTE_UP string = "!voteUp"
var CMD_VOTE_DOWN string = "!voteDown"
var CMD_HELP = "!help"
var CMD_VOTES = "!votes"
var bot *IRCBot
var votes map[string]int
var connection net.Conn
var respReq *textproto.Reader


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
  votes = make(map[string]int)
  bot = createBot()
  connection, _ = bot.ServerConnect()
  fmt.Fprintf(connection, "USER %s 8 * :%s\n", bot.Nick, bot.Nick)
  fmt.Fprintf(connection, "NICK %s\n", bot.Nick)
  fmt.Fprintf(connection, "JOIN %s\n", bot.Channel)
  defer connection.Close()

  reader := bufio.NewReader(connection)
  respReq = textproto.NewReader(reader)
  for {
    line, err := respReq.ReadLine()
    if err != nil {
      break
    }
    fmt.Printf("\033[93m%s\n", line)
    if strings.Contains(line, "PING") {
      fmt.Printf("\033[92mPONG\n")
      fmt.Fprintf(connection, "PONG\n")
    }
    if strings.Contains(line, bot.Channel) && strings.Contains(line, CMD_VOTE_UP) {
      err = voteUp(line)
      if err != nil {
        break
      }
    }
    if strings.Contains(line, bot.Channel) && strings.Contains(line, CMD_VOTE_DOWN) {
      err = voteDown(line)
      if err != nil {
        break
      }
    }
    if strings.Contains(line, bot.Channel) && strings.Contains(line, CMD_HELP) {
      help()
    }
    if strings.Contains(line, bot.Channel) && strings.Contains(line, CMD_VOTES) {
      getVotes()
    }
  }
}

func voteUp(line string) error {
  commandLine := line[strings.Index(line, CMD_VOTE_UP):len(line)]
  if strings.Index(line, CMD_VOTE_UP) != -1 {
    commands := strings.Split(commandLine, " ")
    if len(commands) ==1 {
      fmt.Fprintf(connection, "PRIVMSG %s :ERROR: %s\n", bot.Channel, NO_USER_ERROR)
      fmt.Printf("\033[91mPRIVMSG %s :ERROR: %s\n", bot.Channel, NO_USER_ERROR)
    } else {
      voteUser := commands[1]
      fmt.Fprintf(connection, "NAMES %s\n", bot.Channel)
      fmt.Printf("\033[92mNAMES %s\n", bot.Channel)
      line, err := respReq.ReadLine()
      if err != nil {
        return err
      }
      if strings.Contains(line, voteUser) {
        votes[voteUser] += 1
        fmt.Fprintf(connection, "PRIVMSG %s :Upvoted: %s\n", bot.Channel, voteUser)
        fmt.Printf("\033[92mPRIVMSG %s :Upvoted: %s\n", bot.Channel, voteUser)
      } else {
        fmt.Fprintf(connection, "PRIVMSG %s :ERROR: %s\n", bot.Channel, USER_NOT_FOUND_ERROR)
        fmt.Printf("\033[91mPRIVMSG %s :ERROR: %s\n", bot.Channel, USER_NOT_FOUND_ERROR)
      }
    }
  }
  return nil
}

func voteDown(line string) error {
  commandLine := line[strings.Index(line, CMD_VOTE_DOWN):len(line)]
  if strings.Index(line, CMD_VOTE_DOWN) != -1 {
    commands := strings.Split(commandLine, " ")
    if len(commands) ==1 {
      fmt.Fprintf(connection, "PRIVMSG %s :ERROR: %s\n", bot.Channel, NO_USER_ERROR)
      fmt.Printf("\033[91mPRIVMSG %s :ERROR: %s\n", bot.Channel, NO_USER_ERROR)
    } else {
      voteUser := commands[1]
      fmt.Fprintf(connection, "NAMES %s\n", bot.Channel)
      fmt.Printf("\033[92mNAMES %s\n", bot.Channel)
      line, err := respReq.ReadLine()
      if err != nil {
        return err
      }
      if strings.Contains(line, voteUser) {
        votes[voteUser] -= 1
        fmt.Fprintf(connection, "PRIVMSG %s :Downvoted: %s\n", bot.Channel, voteUser)
        fmt.Printf("\033[92mPRIVMSG %s :Downvoted: %s\n", bot.Channel, voteUser)
      } else {
        fmt.Fprintf(connection, "PRIVMSG %s :ERROR: %s\n", bot.Channel, USER_NOT_FOUND_ERROR)
        fmt.Printf("\033[91mPRIVMSG %s :ERROR: %s\n", bot.Channel, USER_NOT_FOUND_ERROR)
      }
    }
  }
  return nil
}

func help() {
  fmt.Fprintf(connection, "PRIVMSG %s :Commands Are: \n", bot.Channel)
  fmt.Fprintf(connection, "PRIVMSG %s %s\n", bot.Channel, CMD_HELP)
  fmt.Fprintf(connection, "PRIVMSG %s %s\n", bot.Channel, CMD_VOTE_UP)
  fmt.Fprintf(connection, "PRIVMSG %s %s\n", bot.Channel, CMD_VOTE_DOWN)
  fmt.Fprintf(connection, "PRIVMSG %s %s\n", bot.Channel, CMD_VOTES)
  fmt.Printf("\033[92mPRIVMSG %s :Commands Are: \n", bot.Channel)
  fmt.Printf("\033[92mPRIVMSG %s %s\n", bot.Channel, CMD_HELP)
  fmt.Printf("\033[92mPRIVMSG %s %s\n", bot.Channel, CMD_VOTE_UP)
  fmt.Printf("\033[92mPRIVMSG %s %s\n", bot.Channel, CMD_VOTE_DOWN)
  fmt.Printf("\033[92mPRIVMSG %s %s\n", bot.Channel, CMD_VOTES)
}

func getVotes() {
  fmt.Fprintf(connection, "PRIVMSG %s :Current Votes Are: \n", bot.Channel)
  fmt.Printf("\033[92mPRIVMSG %s :Current Votes Are: \n", bot.Channel)
  for key, value := range votes {
    fmt.Fprintf(connection, "PRIVMSG %s :+%s: %d \n", bot.Channel, key, value)
    fmt.Printf("\033[92mPRIVMSG %s :+%s: %d \n", bot.Channel, key, value)
  }
}

// func indexOf(value string, mySlice []string) int {
//   for index, curValue := range mySlice {
//     if curValue == value {
//       return index
//     }
//   }
//   return -1
// }
