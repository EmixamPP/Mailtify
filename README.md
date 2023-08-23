# <h1 align="center">mailtify</h1>
A simple server to easily send notifications via a RESTful API and receive them as mails.

# Quick utilization
### Create a token:
`curl mailtity.domain.com/new`

-> `{"status": 200 "message":"<token>"}`

### Send a message:
`curl "mailtify.domain.com/msg?token=<token>" -F "title=my title" -F "message=my message" -F "recipients=my@mailaddr.com"`

### Delete a token:
`curl "mailtify.domain.com/del?token=<token>"`

# Which features comes next?
1. User accounts
2. Docker container
3. Retry on network error
4. SSL
