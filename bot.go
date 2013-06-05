/* This is an IRC BOT it joins the specified channel below
 * by the specified nickname and sits waiting for commands.
 *
 * It can accept votes for users up or down for silly
 * internet points, the best kind of points.
 * TODO - Make a simple timer to check for 1 vote an hour/half hour/ten minutes.
 * TODO - Make a simple timer to only allow for 1 help per minute.
 * TODO - Put votes in IRCBot struct.
 * TODO - Keep track of who voted for who.
 * TODO - Stop self voting.
 */
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

<<<<<<< HEAD
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
var respReq *textproto.Reader

//Tried to add most of the commands into the bot
//in order to keep from having too many globals.
//feel free to change these back.

type IRCBot struct {
  Server      string
  Port        string
  User        string
  Nick        string
  Channel     string
  Pass        string
  Votes       map[string]int
  BotErrors   Errors
  Cmds        Commands
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

//Different commands the bot has
type Commands struct {
  var CMD_VOTE_UP string = "!voteUp"
  var CMD_VOTE_DOWN string = "!voteDown"
  var CMD_HELP = "!help"
  var CMD_VOTES = "!votes"
}

//Different errors the bot can throw
type Errors struct {
  var NO_USER_ERROR string = "!No user specified!"
  var USER_NOT_FOUND_ERROR string = "!User not in channel!"
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

//connects the bot to the server & puts the connection in the bot
func (bot *IRCBot) ServerConnect() (connection net.Conn, err error) {
  connection, err = net.Dial("tcp", bot.Server)
  if err != nil {
    log.Fatal("Unable to connect to the specified server", err)
  }
  log.Printf("Successfully connected to %s ($s)\n", bot.Server, bot.Connection.RemoteAddr())
  return bot.Connection, nil
}

//Starts the bot, initalizes votes to 0 & starts
//a constant loop reading in lines & waiting for commands.
func main() {
  votes = make(map[string]int)
  bot = createBot()
  bot.Connection, _ = bot.ServerConnect()
  sendCommand("USER", bot.Nick, "8 *", bot.Nick)
  sendCommand("NICK", bot.Nick)
  sendCommand("JOIN", bot.Channel)
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
    //Sends the PONG command in order not to get kicked from the server
    fmt.Printf("\033[93m%s\n", line)
      if strings.Contains(line, "PING") {
        sendCommand("PONG", "")
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

//Sends the specified command along with any parameters needed as well
func sendCommand(command string, parameters ...string) {
  msg := strings.Join(parameters, " ")
  cmd := fmt.Sprintf("%s %s", command, msg)
  fmt.Fprintf(bot.Connection, strings.Join([]string{cmd, "\n"}, ""));
  fmt.Printf(strings.Join([]string{cmd, "\n"}, ""));
}

//Uses the sendCommand to PRIVMSG the channel 
//as well as sends a []string message
func sendMessage(recipient string, message ...string) {
  msg := fmt.Sprintf("%s : %s", recipient, strings.Join(message, " "))
  sendCommand("PRIVMSG", msg)
}

//Votes the user up & saves it to votes.
//If the user isnt in the channel or isnt speicified it will complain
func voteUp(line string) error {
  commandLine := line[strings.Index(line, CMD_VOTE_UP):len(line)]

  if strings.Index(line, CMD_VOTE_UP) != -1 {
    commands := strings.Split(commandLine, " ")
    if len(commands) ==1 {
      sendMessage(bot.Channel, "Error: ", NO_USER_ERROR)
    } else {
      voteUser := commands[1]
      contained, err := inChannel(voteUser)
      if err != nil {
        return err
      }
      if contained {
        votes[voteUser] += 1
        sendMessage(bot.Channel, "Upvoted:", voteUser)
      } else {
        sendMessage(bot.Channel, "Error:", USER_NOT_FOUND_ERROR)
      }
    }
  }
  return nil
}

//Votes the user down & saves it to votes.
//If the user isnt in the channel or isnt speicified it will complain
func voteDown(line string) error {
  commandLine := line[strings.Index(line, CMD_VOTE_DOWN):len(line)]
  if strings.Index(line, CMD_VOTE_DOWN) != -1 {
    commands := strings.Split(commandLine, " ")
    if len(commands) ==1 {
      sendMessage(bot.Channel, "Error:", NO_USER_ERROR)
    } else {
      voteUser := commands[1]
      contained, err := inChannel(voteUser)
      if err != nil {
        return err
      }
      if contained {
        votes[voteUser] -= 1
        sendMessage(bot.Channel, "Downvoted:", voteUser)
      } else {
        sendMessage(bot.Channel, "Error:", USER_NOT_FOUND_ERROR)
      }
    }
  }
  return nil
}

//Prints a help menu to the channel
func help() {
  sendMessage(bot.Channel, "Commands are: \n")
  sendMessage(bot.Channel, CMD_HELP)
  sendMessage(bot.Channel, CMD_VOTE_UP)
  sendMessage(bot.Channel, CMD_VOTE_DOWN)
  sendMessage(bot.Channel, CMD_VOTES)
}

//Gets the current list of votes and sends them to the
//server
func getVotes() {
  sendMessage(bot.Channel, "Current votes are:")
  for key, value := range votes {
    sendMessage(bot.Channel, key, ":", strconv.Itoa(value))
  }
}

//Gets a list of the names in the channel and creates a []string from them
func getNames() ([]string, error) {
  sendCommand("NAMES", bot.Channel)
  line, err := respReq.ReadLine()
  if err != nil {
    return nil, err
  }
  nameLine := line[strings.Index(line, bot.Channel)+len(bot.Channel)+2:len(line)]
  names := strings.Split(nameLine, " ")
  for index, curName := range names {
    if strings.HasPrefix(curName, "@") {
      names[index] = curName[1:len(curName)]
    }
  }
  return names, nil
}

//checks if the user listed is in the channel
func inChannel(name string) (bool, error) {
  chanNames, err := getNames()
  if err != nil {
    return false, err
  }
  for _, curName := range chanNames {
    if curName == name {
      return true, nil
    }
  }
  return false, nil
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
