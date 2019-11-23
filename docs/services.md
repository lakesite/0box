# services #

0box provides the following services, built around our initial use cases:

1. Mail reports and populate particular templates with data.
2. Mail arbitrary people.
3. Get system user's mail (as html, json, whatever)
4. Delete a system user's mail.
5. Forward a system user's mail.
6. Secure the API with an API key for all operations.

[1, 2]: POST /api/0box/v1/mail/

  > template: unique name for template under cwd/templates/<template>
  > context: Data for template, if defined.
  > from: e-mail sent from
  > to: [e-mail(s)]
  > body: if no template is defined, send plain text body content.

Returns 200 OK if successful.

[3]: GET /api/0box/v1/mail/<username from /var/spool/mail>[/mail #]

Returns all mail (or if mail # is specified, the message # from that box) from
the username specified in json format.

Returns "No such message if no message number." if no message # found.

Returns 200 OK if successful.

[4]: DELETE /api/0box/v1/mail/<username>[/mail #]

Deletes all mail (or specified mail #) for the username specified.
Returns JSON: "No such message to delete." if no message # found.

Returns 200 OK if successful.

[5]: Use case handled via 1, 2, and 3.

[6]: Use case handled via APIKeyMiddleware.
