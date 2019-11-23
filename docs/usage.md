# usage #

0box runs as a service on a linux system and provides API endpoints interacted
with using an API key/session key.

For help on the command line:

 $ ./0box -h

## Running ##

First, make sure you have 0box configured correctly.  A basic configuration has
an apikey defined, and the mbox root path, which defaults to /var/mail.

config.toml:

```
[0box]
apikey="secretkey"
mboxroot="/var/mail"
```

  # Start the daemon with the 0box.toml configuration:
  $ ./0box -c 0box.toml

  # Get all of root's e-mail:
  $ curl -H "apikey: secretkey" http://localhost:6999/api/0box/v1/mail/root

  # Get the second email in root's mbox:
  $ curl -H "apikey: secretkey" http://localhost:6999/api/0box/v1/mail/root/2

  # Delete the second email in root's mbox:
  $ curl -H "apikey: my_key" -X DELETE http://localhost:6999/api/0box/v1/mail/root/2

  # Delete *all* email in root's mbox:
  $ curl -H "apikey: my_key" -X DELETE http://localhost:6999/api/0box/v1/mail/root/

  # Send a message to root:
  $ curl -H "api_key: secretkey" -X POST -d "from=localuser&to=root&subject=hello&body=the%20body" http://localhost:6999/api/0box/v1/mail/

* The environment variable, 0BOX_API_KEY, is checked before the apikey directive
  in config.toml.

* Currently, multiple recipients are not supported.

* When using curl, you can specify -v to get more verbose information for
  debugging.

## API ##

Current API endpoints:

/api/0box/v1/                      - Management requests for 0box itself.

/api/0box/v1/mail                  - GET, List available mailboxes.
/api/0box/v1/mail                  - POST, Send mail.
/api/0box/v1/mail/{user}           - GET, retrieve all mailbox user's messages.
/api/0box/v1/mail/{user}           - DELETE, remove all mailbox user's messages.
/api/0box/v1/mail/{user}/{message} - GET, retrieve message # from user's messages.
/api/0box/v1/mail/{user}/{message} - DELETE, delete message # from user's messages.
