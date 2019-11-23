# 0box #

📬 drop-in mailer API	📬

## motivation ##

Often times, you want a standard API for mail operations from a system.  The
main use cases are:

  1. Mail reports and populate particular templates with data.
  2. Mail arbitrary people.
  3. Get system user's mail (as html, json, whatever)
  4. Delete a system user's mail.
  5. Forward a system user's mail.
  6. Secure the API with an API key for all operations.

* 0box aims to provide API endpoints for managing mbox formatted mail.

  For documentation on these API endpoints, see [services](docs/services.md).

* Area of focus:
  - systems mail API.
  - app integration.

Further [rationale](docs/rationale.md) provided.

To see how 0box works, see the [architecture](docs/architecture.md).
To see what services 0box provides, see the [services](docs/services.md).

## usage ##

Please see [usage](docs/usage.md).

## development ##

To run locally and develop, see [development.md](docs/development.md)

## license ##

MIT - See [LICENSE.md](LICENSE.md)

## contributing ##

Please review [standards](docs/standards.md) before submitting issues and pull
requests.  Thank you in advance for feedback, criticism, and feature requests.
