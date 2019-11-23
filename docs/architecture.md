# architecture #

a high-level overview of what might be

The big question here is, do we use a system level dependency, such as GNU
mailutils and `mail` to mbox formatted messages under /var/mail/, or do we use
someone's go library to read the files under /var/mail?  It seems ideal to
parse messages under /var/mail, and expose certain methods to our API endpoints;

  count      - how many messages does the mbox have?
  read [#]   - read message #, or if no # provided, all messages.
  delete [#] - delete message #, or if no # provided, delete all messages.

Looking at how we'd interface with `mail`, the amount of work to create system
calls, parse mail's output, etc, is undesirable.

Dependencies we could use that are maintained and reputable as of 20 Nov
include:

  1. [emersion/go-mbox](https://github.com/emersion/go-mbox)
  The last commit was on 5 Sep 2019.  MIT Licensed and seemingly well maintained.
  This is originally written by [blabber](https://github.com/blabber)

  2. [ProtonMail/go-mbox](https://github.com/ProtonMail/go-mbox)
  A big name like ProtonMail, but a fork of emersion/go-mbox, which has a
  different (older?) api.  I haven't looked to see the buffering differences
  between scanner.go and emersion's reader.go.

emersion's version seems like the best choice.  There are a variety of golang
mbox utilities on github, but most are either forks or prototypes.  Simon Ser
explains that he's using go-mbox in production.

Golang's mail.ReadMessage sets mime headers in a message struct, that is out of
order with what you might see in the /var/mail/<user> file.  Because of how I'm
currently handling message deletion, this means when we write out changes to 
the file, it will be different.  This should be fixed.
