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
### Adding Blog Post Links
Blog post links are randomly selected from the Postgres database.  To populate this database, we simply connect using the psql client and add them.  A sample of the queries is included under samples/sample_data.sql.

To get the database URI, use "cf env" and get it from the credentials under VCAP_SERVICES.
```
$ cf env guestbook
Getting env variables for app guestbook in org starkandwayne / space development as jrbudnack@starkandwayne.com...
OK

System-Provided:
{
  "VCAP_SERVICES": {
    "elephantsql": [
    {
      "credentials": {
        "max_conns": "5",
        "uri": "postgres://bcknjfbm:jkH9Torzv1W6xH-xXkAFzJXYe6fYM9Ck@babar.elephantsql.com/bcknjfbm"
        },
        "label": "elephantsql",
        "name": "guestbook-pg",
        "plan": "turtle",
        "tags": [
        "Data Stores",
        "Cloud Databases",
        "Developer Tools",
        "Data Store",
        "postgresql",
        "relational",
        "New Product"
        ]
      }
      ]
    }
  }

  No user-defined env variables have been set
```

To populate the sample data, run the following command (URI will differ):
```
psql postgres://bcknjfbm:jkH9Torzv1W6xH-xXkAFzJXYe6fYM9Ck@babar.elephantsql.com/bcknjfbm -f samples/sample_data.sql
```
