Twitch Chat TUI in Go


| Command | Description | Usage |
|---------|------------|--------|
| chat | - Fetches the Twitch live stream chat in the terminal.<brCurrently, you can only view messages. Sending messages requires authentication. | `twich chat -u "<username>"`<br> or <br>`twich chat --username "<username>"` |

## Example
### Twitch Chat
```
twitch chat -u "qoqsik"
```
### result
```
➜  ~ twich chat -u "qoqsik"
20:24:04 [Jeffthermite]: yes ask them
20:24:05 [garyhest]: what is the name of this mall?
20:24:08 [RecoveringRumpAddict]: yes ask
20:24:11 [crackpot_lady]: Go straight
20:24:13 [shadowtux]: yeah. ask
20:24:13 [suSpremacist]: go straight
```