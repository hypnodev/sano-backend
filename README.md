# Sano

An easy app that can run in background and check if your services are online or not. With a simple dashboard included!

Built using [Go](https://go.dev/), [MongoDB](https://www.mongodb.com/) and [Vue3](https://vuejs.org/).

## Configuration (Backend)
First time you will run the app, the app will be killed and a new file called `sano.yml` will be created, open it up and configure it as following:
```yaml
app:
  port: 8080 # this will be the webserver port, if isn't available it will pick the next one available
  database:
    url: mongodb+srv://mongo@localhost/?retryWrites=true # mongodb connection url

healthCheck:
  cron: "@hourly" # default cron to check services, you can leave as it is or decrease/increase. See this for format: https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format

services:
  - name: backend # Service name
    displayName: Sano Backend # Name displayed on Frontend app
    url: https://google.com # Url to check
    cron: 0 */2 * * * # cron time. See this for format: https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format

  - name: frontend
    url: http://localhost:8080
    # You can omit cron and the displayName, if cron isn't indicated we'll pick the healthCheck.cron and name for the displayName

# If you want to add a new services, just copy/paste this block and edit the setting
# - name: microservice
#   ......
```
Now you can rerun the executable and enjoy the health check :)

## Configuration ([Frontend](https://github.com/hypnodev/sano-frontend))
```bash
$ npm install
$ VUE_APP_API_URL=https://localhost:8080 npm run serve # Change the API url if the your is different
```

## Contribute
Feel free to push any changes, improve anything and fix stuff.

## Stay in touch
- Author - [Cristian Cosenza](https://linkedin.com/in/cristiancosenza)
