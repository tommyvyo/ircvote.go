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
)

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
  //serverPrompt()
  return &IRCBot {
    Server:     "irc.freenode.net",
    Port:       "6665",
    Nick:       "DeanBot",
    Channel:    "#orderdeck",
    Pass:       "",
    Connection: nil,
    User:       "VoteBot",
    }
}

/*On creation of a bot prompts the user for correct server information
func severPrompt() {
}*/

//connects the bot to the server & puts the connection in the bot
func (bot *IRCBot) ServerConnect() (connection net.Conn, err error) {
  connection, err = net.Dial("tcp", bot.Server + ":" + bot.Port)
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

// func indexOf(value string, mySlice []string) int {
//   for index, curValue := range mySlice {
//     if curValue == value {
//       return index
//     }
//   }
//   return -1
// }
