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
  "io/ioutil"
)

// Set this to true to always output new lines from socket
const debug = true

// Set this to true to load bot configuration from ./server.config instead of
// prompting on startup
const use_config_file = false

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
  results["nick"] = strings.TrimSpace(nick)
  results["channel"] = strings.TrimSpace(channel)
  results["password"] = strings.TrimSpace(password)

  fmt.Printf(strings.Join([]string{"Connecting to ", results["host"], " and joining ", results["channel"]},""));
  return results
}

func prompt(r *bufio.Reader, msg string) string {
  fmt.Print(msg)
  value, _ := r.ReadString('\n')
  return strings.TrimSpace(value)
}

func configFromPrompt() []string {
  fmt.Print("Starting server from prompt.\n")
  reader := bufio.NewReader(os.Stdin)
  config := make([]string, 4)
  config[0] = prompt(reader, "Enter a host to connect to: (e.g. irc.freenode.net) ")
  config[1] = prompt(reader, "Enter a port to connect to: (e.g. 6667) ")
  config[2] = prompt(reader, "Enter a nickname: ")
  config[3] = prompt(reader, "Enter a channel to join: (e.g. #orderdeck) ")
  return config
}

func configFromFile() []string {
  fmt.Print("Starting server from config file.\n")
  bytes, _ := ioutil.ReadFile("server.config")
  return strings.Split(string(bytes), "\n")
}

//create a new bot with the given server info
//currently static but will create prompt for any server
func createBot() *IRCBot {

  var config []string

  if (use_config_file) {
    config = configFromFile()
  } else {
    config = configFromPrompt()
  }

  //serverPrompt()
  return &IRCBot {
    Server:     config[0]+":"+config[1],
    Nick:       config[2],
    Channel:    config[3],
    Pass:       "",
    Connection: nil,
    User:       "ircvote-bot",
  }
}

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

    command := strings.Split(line, " ")

      if (debug) {
        fmt.Printf("%s\n", line)
      }


    if (command[1] == "PRIVMSG") {
      privateMessageReceived(command)
    } else {
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
}

func privateMessageReceived(command []string) {
  var target, message = command[2], command[3]
  fmt.Printf("%s\n", target + ": " + message)
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

/*
Things To Implement:
:tommyvyo!~tommyvyo@unaffiliated/tommyvyo PRIVMSG #orderdeck :ay0o
:tommyvyo!~tommyvyo@unaffiliated/tommyvyo PRIVMSG #orderdeck :wsup t0myvy0oo0
:tommyvyo!~tommyvyo@unaffiliated/tommyvyo PRIVMSG #orderdeck :http://drupal.org/node/225125
:tommyvyo!~tommyvyo@unaffiliated/tommyvyo PRIVMSG t0myvy0oo0 :yo yo yo
PING :calvino.freenode.net
PONG
:tommyvyo!~tommyvyo@unaffiliated/tommyvyo PRIVMSG t0myvy0oo0 :PING 1369865522.308181
PONG
:tommyvyo!~tommyvyo@unaffiliated/tommyvyo PRIVMSG t0myvy0oo0 :TIME
:tommyvyo!~tommyvyo@unaffiliated/tommyvyo PRIVMSG t0myvy0oo0 :VERSION
:tommyvyo!~tommyvyo@unaffiliated/tommyvyo PRIVMSG t0myvy0oo0 :USERINFO
PING :calvino.freenode.net
PONG
:tommyvyo!~tommyvyo@unaffiliated/tommyvyo PRIVMSG t0myvy0oo0 :CLIENTINFO
*/

// func indexOf(value string, mySlice []string) int {
//   for index, curValue := range mySlice {
//     if curValue == value {
//       return index
//     }
//   }
//   return -1
// }
