# Guestbook
This application is used to gather email addresses for the drawing at the Stark & Wayne booth at CF Summit.
### Deployment to Cloud Foundry
To instantiate the database, we use ElephantSQL to provision an instance of Postgres:
```
cf create-service elephantsql turtle guestbook-pg
```
To deploy the application, issue the following commands:
```
cf push guestbook -m 256M -d starkandwayne.com --no-start
cf bs guestbook guestbook-pg
cf start guestbook
```
