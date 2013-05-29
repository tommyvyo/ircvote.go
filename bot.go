package main

import (
  "fmt"
  "bufio"
  "log"
  "net"
  "net/textproto"
  "strings"
  "strconv"
  "os"
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
  results := startupPrompt()
  //serverPrompt()
  return &IRCBot {
    Server:     results["host"],
    Nick:       results["nick"],
    Channel:    results["channel"],
    Pass:       results["pass"],
    Connection: nil,
    User:       results["nick"],
  }
}

/*On creation of a bot prompts the user for correct server information
func severPrompt() {
}*/

//connects the bot to the server
func (bot *IRCBot) ServerConnect() (connection net.Conn, err error) {
  connection, err = net.Dial("tcp", bot.Server)
  if err != nil {
    log.Fatal("Unable to connect to the specified server", err)
  }
  bot.Connection = connection
  log.Printf("Successfully connected to %s ($s)\n", bot.Server, bot.Connection.RemoteAddr())
  return bot.Connection, nil
}

func startupPrompt() map[string]string {
  reader := bufio.NewReader(os.Stdin)
  var (
    host string
    port string
    nick string
    channel string
    password string
  )

  fmt.Print("Enter a hostname you'd like to connect to: (e.g: irc.freenode.net) ")
  host, _ = reader.ReadString('\n')
  fmt.Print("Enter a port: (e.g: 6667) ")
  port, _ = reader.ReadString('\n')
  fmt.Print("Enter a nickname: ")
  nick, _ = reader.ReadString('\n')
  fmt.Print("Enter a channel: (e.g: #orderdeck) ")
  channel, _ = reader.ReadString('\n')
  fmt.Print("Enter a password:")
  password, _ = reader.ReadString('\n')

  results := make(map[string]string)
  results["host"] = strings.Join([]string{strings.TrimSpace(host), strings.TrimSpace(port)}, ":");
  results["nick"] = nick
  results["channel"] = channel
  results["password"] = password

  fmt.Printf(strings.Join([]string{"Connecting to ", results["host"], " and joining ", results["channel"]},""));
  return results
}

func main() {
  votes = make(map[string]int)
  bot = createBot()
  connection, _ = bot.ServerConnect()
  sendCommand("USER", []string{bot.Nick, "8 *", bot.Nick})
  sendCommand("NICK", []string{bot.Nick})
  sendCommand("JOIN", []string{bot.Channel})
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
      sendCommand("PONG", []string{})
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

func sendCommand(command string, parameters []string) {
  msg := strings.Join(parameters, " ")
  cmd := fmt.Sprintf("%s %s", command, msg)
  fmt.Fprintf(connection, strings.Join([]string{cmd, "\n"}, ""));
  fmt.Printf(strings.Join([]string{cmd, "\n"}, ""));
}

func sendMessage(recipient string, message []string) {
  msg := strings.Join(message, " ");
  sendCommand("PRIVMSG", []string{recipient, ":", msg}) 
}

func voteUp(line string) error {
  commandLine := line[strings.Index(line, CMD_VOTE_UP):len(line)]

  if strings.Index(line, CMD_VOTE_UP) != -1 {
    commands := strings.Split(commandLine, " ")
    if len(commands) ==1 {
      sendMessage(bot.Channel, []string{"ERROR: ", NO_USER_ERROR})
    } else {
      voteUser := commands[1]
      sendCommand("NAMES", []string{bot.Channel})
      line, err := respReq.ReadLine()
      if err != nil {
        return err
      }
      if strings.Contains(line, voteUser) {
        votes[voteUser] += 1
        sendMessage(bot.Channel, []string{"Upvoted:", voteUser})
      } else {
        sendMessage(bot.Channel, []string{"Error:", USER_NOT_FOUND_ERROR}) 
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
      sendMessage(bot.Channel, []string{"ERROR:", NO_USER_ERROR})
    } else {
      voteUser := commands[1]
      sendCommand("NAMES", []string{bot.Channel})
      line, err := respReq.ReadLine()
      if err != nil {
        return err
      }
      if strings.Contains(line, voteUser) {
        votes[voteUser] -= 1
        sendMessage(bot.Channel, []string{"Downvoted:", voteUser}) 
      } else {
        sendMessage(bot.Channel, []string{"ERROR:", USER_NOT_FOUND_ERROR}) 
      }
    }
  }
  return nil
}

func help() {
  sendMessage(bot.Channel, []string{"Commands are: \n"})
  sendMessage(bot.Channel, []string{CMD_HELP})
  sendMessage(bot.Channel, []string{CMD_VOTE_UP})
  sendMessage(bot.Channel, []string{CMD_VOTE_DOWN})
  sendMessage(bot.Channel, []string{CMD_VOTES})
}

func getVotes() {
  sendMessage(bot.Channel, []string{"Current votes are:"})
  for key, value := range votes {
    sendMessage(bot.Channel, []string{key, ":", strconv.Itoa(value)})
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
