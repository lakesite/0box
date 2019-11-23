# rationale #

So, why create 0box?

1. You have a lot of systems and you want to aggregate their mail via API to
   a reporting system, without depending on these systems to have;

   a. Google compat; e.g. forward and reverse DNS match.
   b. mail exchangers configured properly (or installed).
   c. /etc/aliases configured to forward everything to one email.

2. You want to interact with mail on remote systems interactively;

   a. List emails by user.
   b. Delete emails by user.
   c. Retrieve emails by user.

3. You want to send mail, from the system, using gmail or what have you;

   a. Send templated email.
   b. Send to a list of people.
   c. Send via third party/gmail.
