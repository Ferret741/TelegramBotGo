// ===========================================================================
// Package
package main




// ===========================================================================
// Import
import (
    "fmt"
    "net/http"
    "net/url"
    "log"
    "os"
    "strings"
    "io"
    "encoding/json"
)




// ===========================================================================
// Variables
var TOKEN_ENVVAR_NAME string = "TELEGRAM_BOT_TOKEN"
var ERR_ARGUMENT_MISSING string = "Option %s is expecting an argument. Please view help file"
var ERR_MESSAGE_NOT_PROVIDED string = "No message provided. Please provide with -m/--message"
var ERR_CHANNEL_NOT_PROVIDED string = "No chat id provided. Please provide with -c/--channel/--channel-id"
var ERR_TOKEN_ENVVAR_NOT_FOUND string = "No environment variable for %s was found"
var ERR_TOKEN_NOT_PROVIDED string = "No token provided. Please refer to help file (-h/--help) to understand token declaration methods"
var ERR_TOKEN_FILE_NOT_EXIST string = "Token file %s does not exist! Please verify file location"
var ERR_TOKEN_FILE_READ string = "Error encountered while reading from %s: %s"
var HELP_PAGE string = "" +
    "\x1B[1;31m\n" +
    "Usage %s <options>\n" +
    "\x1B[0;31m\n" +
    "+ ================================================================================================ +\n" +
    "| \x1B[1;31mRequired Arguments                                                                               \x1B[0;31m|\n" +
    "+ ============== + ===== + ======================================================================= +\n" +
    "| Long           | Short | Description                                                             |\n" +
    "+ ============== + ===== + ======================================================================= +\n" +
    "| --token        |  -t   | The telegram bot token. Can be provided in the following manners:       |\n" +
    "|                |       | - Environment variable: Argument should be env (case insensitive), and  |\n" +
    "|                |       |   the telegram bot token stored in an env var named TELEGRAM_BOT_TOKEN  |\n" +
    "|                |       |   Note: the value of token must be provided, although it is not         |\n" +
    "|                |       |   required for this option and its argument to appear on the command    |\n" +
    "|                |       |   line. Should they be absent, the default action will be to search     |\n" +
    "|                |       |   for the environment variable.                                         |\n" +
    "|                |       |                                                                         |\n" +
    "|                |       |                                                                         |\n" +
    "|                |       | - File: Argument should be the path to a file containing a single line  |\n" +
    "|                |       |   that is the telegram bot token. Path of the file should be prepended  |\n" +
    "|                |       |   with 'file:' (see Examples below)                                     |\n" +
    "|                |       |                                                                         |\n" +
    "|                |       | - Plain text argument: Argument is the plain text bot token string      |\n" +
    "|                |       |                                                                         |\n" +
    "|                |       | \x1B[2;31mExamples: --token env                                                   \x1B[0;31m|\n" +
    "|                |       | \x1B[2;31m          --token xbot-token-1111-2222                                  \x1B[0;31m|\n" +
    "|                |       | \x1B[2;31m          --token file:/tmp/telegram-bot-token                          \x1B[0;31m|\n" +
    "+ -------------- + ----- - ----------------------------------------------------------------------- +\n" +
    "| --channel      |  -c   | A comma-delimited list of channel IDs that will receive the message.    |\n" +
    "| --channel-id   |       | The channel IDs are typically an integer beginning with -100 and        |\n" +
    "|                |       | can be obtained by copy/pasting an element from a chennl in the webUI,  |\n" +
    "|                |       | or using the Json Dump Bot. This also accepts underscore variants that  |\n" +
    "|                |       | are common for forum/thread/topics (e.g. -10011111111_800)              |\n" +
    "|                |       |                                                                         |\n" +
    "|                |       | \x1B[2;31mExample: --target -1001111111111,-100222222222                          \x1B[0;31m|\n" +
    "+ -------------- + ----- - ----------------------------------------------------------------------- +\n" +
    "| --message      |  -m   | The message to send. This can be given multiple times to send more than |\n" +
    "|                |       | one message in a single command.                                        |\n" +
    "|                |       |                                                                         |\n" +
    "|                |       | \x1B[2;31mExample: --message 'Test message :exclamation:'                         \x1B[0;31m|\n" +
    "+ -------------- + ----- + ----------------------------------------------------------------------- +\n" +
    "\n" +
    "\n" +
    "+ ============================================================================================= +\n" +
    "| \x1B[1;31mOptional Arguments                                                                            \x1B[0;31m|\n" +
    "+ =========== + ===== + ======================================================================= +\n" +
    "| Long        | Short | Description                                                             |\n" +
    "+ =========== + ===== + ======================================================================= +\n" +
    "|             |       |                                                                         |\n" +
    "+ ----------- + ----- + ----------------------------------------------------------------------- +\n" +
    "\n" +
    "+ ============================================================================================= +\n" +
    "| \x1B[1;31mSwitches                                                                                      \x1B[0;31m|\n" +
    "+ =========== + ===== + ======================================================================= +\n" +
    "| Long        | Short | Description                                                             |\n" +
    "+ =========== + ===== + ======================================================================= +\n" +
    "| --help      |  -h   | Show this page                                                          |\n" +
    "+ ----------- + ----- + ----------------------------------------------------------------------- +\n" +
    "\n" +
    "\n" +
    "Examples:\n" +
    "\n" +
    "   \x1B[2;31mSend to #channel_name_1, using a token provided in a file \n" +
    "   \x1B[0;31m$ %s --token file:/tmp/telegram.token --channel-id '-1005555555' --message 'This is a test message'\n" +
    "\n" +
    "   \x1B[2;31mSend a message to two different channels, using a token provided by environment variable\n" +
    "   \x1B[0;31m$ %s --token env --target '-10033333333,-10088888888' --message 'Please check alerts!!'\n" +
    "\x1B[0m\n" +
    "\n"




// ===========================================================================
// Types
type Telegram struct {

    // This is the token used for bot authorisation. It can be
    // obtained from the bot registration with BotFather
    Token string

    // The base URL from which all other API events are based
    BaseURL string

    // The message string that will be sent
    Messages []string

    // The message ID to which we plan on responding (if there
    // is one)
    MessageThreadID string

    // The channel ID. Typically for channels this will be an
    // integer beginning in -100. Copy/pasting an element from
    // a channel in the webUI will provide you with the channel
    // ID, or you can use Json Dump Bot
    ChannelID string
}




// ===========================================================================
// Functions

// Take the CLI arguments and
func populateTelegramConfig(telegram *Telegram) {

    // Verify that the array length is of sufficient
    // length for any switch requiring an argument
    // e.g. len(os.Args) must be at least index+2
    verifyArgsLength := func(position int, argument string) {
        if position + 1 >= len(os.Args) {
            msg := fmt.Sprintf(ERR_ARGUMENT_MISSING, argument)
            log.Fatal(msg)
        }
    }

    // Iterate over the CLI argument list
    for index, value := range os.Args {

        switch value {
            case "--message", "-m":
                verifyArgsLength(index, value)
                telegram.Messages = append(telegram.Messages, os.Args[index+1])
                // telegram.Message = os.Args[index+1]

            case "--token", "-t":
                verifyArgsLength(index, value)

                // If the argument is "ENV" then check to
                // see if there is a TELEGRAM_BOT_TOKEN entry
                if strings.ToUpper(os.Args[index+1]) == "ENV" {

                    // Entry was found: assign telegram instance token
                    if token := telegram.getTokenFromEnvironment(); token != "" {
                        telegram.Token = token

                    // Entry was not found: Exit after informing the user
                    } else {
                        msg := fmt.Sprintf(ERR_TOKEN_ENVVAR_NOT_FOUND, TOKEN_ENVVAR_NAME)
                        log.Fatal(msg)
                    }

                // Assume that the token was provided on the
                // command line argument list
                } else if strings.HasPrefix(os.Args[index+1], "file:") {
                    token_file := strings.Split(os.Args[index+1],":")[1]
                    telegram.read_token_file(token_file)

                // Assume that the token was provided on the
                // command line argument list
                } else {
                    telegram.Token = os.Args[index+1]
                }


            case "--channel-id","--channel","-c":
                verifyArgsLength(index, value)
                telegram.ChannelID = os.Args[index+1]


            case "--help","-h":
                fmt.Printf(HELP_PAGE, os.Args[0], os.Args[0], os.Args[0])
                os.Exit(0)
        }
    }
}


// --------------------------------------------------
// Perform the actual sending of an HTTP request
func sendHttpPostFormRequest(send_url string, send_values url.Values) (*http.Response, error) {

    // Perform the request and obtain the results
    return http.PostForm(send_url, send_values)
}


// --------------------------------------------------
// Simple way to output things. We could use this
// to define debugging levels later on
func output(message, level string){

    // Default colour
    var colour string = "\x1B[1;34m"

    // Boolean flag for printing the message
    var print_me bool = true

    // Conditionalise based off the given level
    switch strings.ToUpper(level) {

        case "DEBUG":
            colour = "\x1B[1;30m"

        case "FAILURE":
            colour = "\x1B[1;31m"

        case "SUCCESS":
            colour = "\x1B[1;32m"

        case "WARNING":
            colour = "\x1B[1;33m"

        case "INFO":
            colour = "\x1B[1;34m"
    }

    // Format the message
    message = fmt.Sprintf("%s%s\x1B[0m\n",colour,message)

    // If the print_me flag is true, then print
    // the message
    if print_me {
        fmt.Printf(message)
    }
}




// ===========================================================================
// Methods

// --------------------------------------------------
// Set defaults for the Telegram object
func (t *Telegram) setDefaults(){
    t.BaseURL           = "https://api.telegram.org/"
    t.Token             = t.getTokenFromEnvironment()
    t.MessageThreadID   = ""
    t.ChannelID         = ""
}


// --------------------------------------------------
// Verify that all the required fields for the
// Telegram object are present
func (t *Telegram) verifyFields(){

    // These are the fields that must have a value, if any
    // of them do not, then we should fail out
    switch {
        case len(t.Messages) == 0:
            log.Fatal(ERR_MESSAGE_NOT_PROVIDED)

        case t.Token == "":
            log.Fatal(ERR_TOKEN_NOT_PROVIDED)

        case t.ChannelID == "":
            log.Fatal(ERR_CHANNEL_NOT_PROVIDED)
    }
}



// --------------------------------------------------
// Send a message
func (t *Telegram) sendMessages() {

    // Create the full URL
    full_url := fmt.Sprintf("%sbot%s/sendMessage", t.BaseURL, t.Token)

    // Create a mapping for the params
    params := url.Values{}
    params.Set("parse_mode", "MarkdownV2")

    // Iterate over the list of targets
    for _, channel_id := range strings.Split(t.ChannelID, ",") {




        // Unset the message_thread_id header value, if it is set
        if params.Has("message_thread_id") {
            params.Del("message_thread_id")
            t.MessageThreadID = ""
        }


        // Assume that if the channel_id has an underscore, there is
        // reference to a message_thread_id.
        if strings.Contains(channel_id,"_") {


            // Split the channel into two portions, the channel_id and
            //the message thread id
            channel_parts    := strings.Split(channel_id,"_")
            channel_id        = channel_parts[0]
            t.MessageThreadID = channel_parts[1]

            // Add the thread message id to the telegram object. We
            // probably don't need to do that here and remove the
            // object attribute unless we plan on using it some other
            // place in the program. We could simply refer to it
            // as channel_parts[1], and then not worry about
            // unsetting it.
            params.Add("message_thread_id", t.MessageThreadID)
        }

        // Set the chat_id as current channel ID
        params.Set("chat_id", channel_id)

        // Iterate over all the message in self.Messages slice
        for _, messageValue := range t.Messages {

            // Set the message text from the messages slice
            params.Set("text", messageValue)

            // Send the request for dispatching
            resp, err := sendHttpPostFormRequest(full_url, params)
            if err != nil {
                fmt.Printf("Error encountered during HTTP POST request: %s", err)
            }

            // Consume the body here and close the Body
            content, _ := io.ReadAll(resp.Body)
            resp.Body.Close()

            // Perform tasks based on return status code
            switch resp.StatusCode {

                // Don't need to do anything here!
                case 200:

                // Anything aside from a 200 status code is
                // probably not good!!
                default:

                    // JSON type to represented comes from testing
                    // Message: {"ok":false,"error_code":400,"description":"Bad Request: chat not found"}
                    // Create a struct type that will hold the
                    // expected values from the body
                    type TelegramErrorMessage struct {
                        Ok  bool
                        Error_code int
                        Description string
                        Hello string
                    }

                    var message TelegramErrorMessage

                    // Perform the unmarshaling here and store the
                    // JSON values in the message struct
                    err_json := json.Unmarshal(content, &message)
                    if err_json != nil {
                        fmt.Printf("Error encountered during JSON unmarshaling: %s", err_json)
                    }

                    // Print out the error
                    fmt.Printf("<channel_id: %s>: [%d] %s\n", channel_id, message.Error_code, message.Description)

            }
        }
    }
}


// --------------------------------------------------
// Read a token from a file
func (t *Telegram) read_token_file(token_file string) {

    // Defalt value for the token string
    var token_string string = ""

    // Read contents from the file
    read_token_string, read_token_error := os.ReadFile(token_file)

    // Handle errors here
    // The file does not exist
    if os.IsNotExist(read_token_error) {
        err_msg := fmt.Sprintf(ERR_TOKEN_FILE_NOT_EXIST, token_file)
        log.Fatal(err_msg)

    // Some other error exists
    } else if read_token_error != nil {
        err_msg := fmt.Sprintf(ERR_TOKEN_FILE_READ, token_file, read_token_error)
        log.Fatal(err_msg)
    }

    // Convert byte array to string and trim
    token_string = strings.TrimSpace(string(read_token_string))

    // Set the token string
    t.Token = token_string
}


// --------------------------------------------------
// Obtain token from environment variables
func (t *Telegram) getTokenFromEnvironment() string {

    // Set a default
    var token_string = ""

    // Set the string if it is found
    if envVar, isFound := os.LookupEnv(TOKEN_ENVVAR_NAME); isFound {
       token_string = envVar
    }

    // Return whatever it is that we have
    return token_string
}




// ===========================================================================
// Main Body
func main(){

    // Create a new configuration instance
    TelegramInstance := &Telegram{}

    // Set the values for the configuration instance
    TelegramInstance.setDefaults()

    // Obtain additional telegram stuff from provided CLI argments
    populateTelegramConfig(TelegramInstance)

    // Ensure that the proper fields have been set
    TelegramInstance.verifyFields()

    // Send the message
    TelegramInstance.sendMessages()

}



// ===========================================================================
// Foooter
// vim:syntax=go:ft=go:sts=4:ts=4:sw=4:et:ai:nu:
