## Usage

~~~
Usage PROG <options>

+ ================================================================================================ +
| Required Arguments                                                                               |
+ ============== + ===== + ======================================================================= +
| Long           | Short | Description                                                             |
+ ============== + ===== + ======================================================================= +
| --token        |  -t   | The telegram bot token. Can be provided in the following manners:       |
|                |       | - Environment variable: Argument should be env (case insensitive), and  |
|                |       |   the telegram bot token stored in an env var named TELEGRAM_BOT_TOKEN  |
|                |       |   Note: the value of token must be provided, although it is not         |
|                |       |   required for this option and its argument to appear on the command    |
|                |       |   line. Should they be absent, the default action will be to search     |
|                |       |   for the environment variable.                                         |
|                |       |                                                                         |
|                |       |                                                                         |
|                |       | - File: Argument should be the path to a file containing a single line  |
|                |       |   that is the telegram bot token. Path of the file should be prepended  |
|                |       |   with 'file:' (see Examples below)                                     |
|                |       |                                                                         |
|                |       | - Plain text argument: Argument is the plain text bot token string      |
|                |       |                                                                         |
|                |       | Examples: --token env                                                   |
|                |       |           --token xbot-token-1111-2222                                  |
|                |       |           --token file:/tmp/telegram-bot-token                          |
+ -------------- + ----- - ----------------------------------------------------------------------- +
| --channel      |  -c   | A comma-delimited list of channel IDs that will receive the message.    |
| --channel-id   |       | The channel IDs are typically an integer beginning with -100 and        |
|                |       | can be obtained by copy/pasting an element from a chennl in the webUI,  |
|                |       | or using the Json Dump Bot. This also accepts underscore variants that  |
|                |       | are common for forum/thread/topics (e.g. -10011111111_800)              |
|                |       |                                                                         |
|                |       | Example: --target -1001111111111,-100222222222                          |
+ -------------- + ----- - ----------------------------------------------------------------------- +
| --message      |  -m   | The message to send                                                     |
|                |       |                                                                         |
|                |       | Example: --message 'Test message :exclamation:'                         |
+ -------------- + ----- + ----------------------------------------------------------------------- +


+ ============================================================================================= +
| Optional Arguments                                                                            |
+ =========== + ===== + ======================================================================= +
| Long        | Short | Description                                                             |
+ =========== + ===== + ======================================================================= +
|             |       |                                                                         |
+ ----------- + ----- + ----------------------------------------------------------------------- +

+ ============================================================================================= +
| Switches                                                                                      |
+ =========== + ===== + ======================================================================= +
| Long        | Short | Description                                                             |
+ =========== + ===== + ======================================================================= +
| --help      |  -h   | Show this page                                                          |
+ ----------- + ----- + ----------------------------------------------------------------------- +


Examples:

   Send to #channel_name_1, using a token provided in a file
   $ PROG --token file:/tmp/telegram.token --channel-id -1005555555 --message 'This is a test message'

   Send a message to two different channels, using a token provided by environment variable
   $ PROG --channel-id env --target -10033333333,-10088888888 --message 'Please check alerts!!'
~~~
