# <h1 align="center">mailtify</h1>
A simple server to easily send notifications via a RESTful API and receive them as mails.

# Quick utilization
### Create a token:
`curl mailtity.domain.com/new`

-> `{"status": 200 "message":"<token>"}`

### Send a message:
`curl "mailtify.domain.com/msg?token=<token>" -F "title=my title" -F "message=my message" -F "recipients=my@domain.com"`

### Delete a token:
`curl "mailtify.domain.com/del?token=<token>"`

# Setup
A `config.yml` must be present in the root directory, see [config-example.yml](config-example.yml) for the structure and required fields.

# Which features come next?
1. SSL
2. User accounts
3. Docker container
4. Retry on network error
