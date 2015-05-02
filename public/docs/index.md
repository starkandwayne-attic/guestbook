# Try Cloud Foundry Docs

## Community Documentation

While the tutorials below contain all the information you need to get up and running, having this link in hand will be extremely useful:

[http://docs.cloudfoundry.org/](Cloud Foundry documentation)


## Install the Cloud Foundry command line interface (CLI)

If you've used the Cloud Foundry CLI in the pasy, your best bet is to remove it before continuing:

	gem uninstall cf

### Windows

Choose one of the links below to download:

- [Windows 64 Bit](https://s3.amazonaws.com/go-cli/releases/v6.1.2/installer-windows-amd64.zip)
- [Windows 32 bit](https://s3.amazonaws.com/go-cli/releases/v6.1.2/installer-windows-386.zip)

Once the download has completed, unzip the file and double click the `cf` program. Follow the instructions in the installation program.

### Mac OSX

Download this installation package:

- [Mac OS X 64 bit](https://s3.amazonaws.com/go-cli/releases/v6.1.2/installer-osx-amd64.pkg)

Open the package file to install. Follow the instructions in the installation program.

### Ubuntu/Debian-based Linux

Choose one of the links below to download:

- [Debian 32 bit](https://s3.amazonaws.com/go-cli/releases/v6.1.2/cf-cli_i386.deb)
- [Debian 64 bit](https://s3.amazonaws.com/go-cli/releases/v6.1.2/cf-cli_amd64.deb)

Choose the appropriate command to run based on the file you downloaded:

	sudo dpkg -i cf-cli_i386.deb

or

	sudo dpkg -i cf-cli_amd64.deb


### Red Hat/Centos Linux

Choose one of the links below to download:

- [Red Hat 32 bit](https://s3.amazonaws.com/go-cli/releases/v6.1.2/cf-cli_i386.rpm)
- [Red Hat 64 bit](https://s3.amazonaws.com/go-cli/releases/v6.1.2/cf-cli_amd64.rpm)


Choose the appropriate command to run based on the file you downloaded:

	sudo rpm –ivh cf-cli_i386.rpm

or

	sudo rpm –ivh cf-cli_amd64.rpm

## Login to Cloud Foundry

To begin using your new Cloud Foundry instance, you'll need to configure the `cf` tool. Login details should will be sent to you when the instance was created, which can take up to 1 hour. If you don't receive an email with login details, please contact <trycf@starkandwayne.com>.

You'll need to enter the following, replacing the AWS IP addresses with those that were provided in email:

	cf api https://api.<your_aws_ip>.xip.io  --skip-ssl-validation
	cf login

An `Email>` prompt will appear, enter `admin`. Next, enter `admin` at the `Password>` prompt.



## Build an example project

### Build a Java project

#### Overview

The Java example is a simple "Hello world" application. While it's not very exciting in and of itself, it will show you the basic process for pushing applications to Cloud Foundry.

This example is built with [Spring Boot](http://projects.spring.io/spring-boot/), however, several other application frameworks are available as part of the [Cloud Foundry Java Buildpack](http://docs.gopivotal.com/pivotalcf/buildpacks/java/):

- [Grails](https://grails.org/)
- [Groovy](http://groovy.codehaus.org/)
- Java main() apps
- [Play Framework](http://www.playframework.com/)
- [Servlets](http://www.oracle.com/technetwork/java/index-jsp-135475.html)
- [Spring Boot](http://projects.spring.io/spring-boot/)


#### Building

Install [Maven](http://maven.apache.org/) if not installed already. You can see if your installation works by running this command:

	mvn

Next, you'll need to download and build the Try CF Java example application. The example application is available to download from [here]().



Unzip this file, and change directory to ./java (TODO: change the project dir)

You can build the example app with:

	mvn package

#### Deploying

Next, push the application to Cloud Foundry:

	cf push trycf_java_example -p ./target/cf_java_example-1.0-SNAPSHOT.jar -m 256m

The `-m 256m` argument starts the application with 128MB of memory. You can scale this up at any later time via the [cli application](http://docs.run.pivotal.io/devguide/installcf/whats-new-v6.html).

When the deployment completes, you should see something like the following in the command output:

```
...
requested state: started
instances: 1/1
usage: 256M x 1 instances
urls: trycf_java_example.54.235.200.94.xip.io
```

Note the value in the last row, `urls`, which is the public location where you app is available. Point your browser at *http://trycf-java-example.<your_aws_ip>.xip.io* to see your app.

*Note that the IP address will not be the same as in the example above.*


### Build a Ruby project


#### Overview

The Ruby example is a simple "Hello world" application. While it's not very exciting in and of itself, it will show you the basic process for pushing applications to Cloud Foundry.

This example is built with [Sinatra](http://www.sinatrarb.com/), however, several other application frameworks are available as part of the [Cloud Foundry Ruby Buildpack](http://docs.gopivotal.com/pivotalcf/buildpacks/ruby/)

- [Ruby](https://www.ruby-lang.org/en/)
- [Rack](http://rack.github.io/)
- [Rails](http://rubyonrails.org/)
- [Sinatra](http://www.sinatrarb.com/)

#### Building

You'll need Ruby, Rubygems and Sinatra installed. If you already have Ruby and Rubygems, just install the Sinatra gem:

	gem install sinatra

There's nothing to build! You can run the example app locally with:

	rackup -p 4567

You should be able to connect to [http://localhost:3000](http://localhost:3000) and see the message "Hello world".

#### Deploying

	cf push trycf_ruby_example -m 256m

When the deployment completes, you should see something like the following in the command output:

```
...
requested state: started
instances: 1/1
usage: 256M x 1 instances
urls: trycf-ruby-example.54.235.200.94.xip.io
```

Note the value in the last row, `urls`, which is the public location where you app is available. Point your browser at *http://trycf-ruby-example.<your_aws_ip>.xip.io* to see your app.

*Note that the IP address will not be the same as in the example above.*



### Build a Node project

#### Overview

The Node example is a simple "Hello world" application. While it's not very exciting in and of itself, it will show you the basic process for pushing applications to Cloud Foundry.

This example is built with [Express](http://expressjs.com/), however, several other application frameworks are available as part of the [Cloud Foundry Node Buildpack](http://docs.gopivotal.com/pivotalcf/buildpacks/node/)


#### Building

You'll need [Node](http://nodejs.org/) and [npm](https://www.npmjs.org/) installed.

Then, download the dependencies with the following:

```
npm install
```

Run the app with:

```
node app.js
```

You should be able to connect to [http://localhost:5000](http://localhost:3000) and see the message "Hello world".

#### Deploying

	cf push trycf_node_example -m 256m


When the deployment completes, you should see something like the following in the command output:

```
...
requested state: started
instances: 1/1
usage: 256M x 1 instances
urls: trycf_node_example.54.235.200.94.xip.io
```

Note the value in the last row, `urls`, which is the public location where you app is available. Point your browser at *http://trycf-node-example.<your_aws_ip>.xip.io* to see your app.

*Note that the IP address will not be the same as in the example above.*

## Frequently Asked Questions

- - - - - - - - -

  * __Q__: How long does it take?
  * __A__: 20-30 minutes, maybe more if the tubes of the internet are clogged.
  
- - - - - - - - -

  * __Q__: Does Try Cloud Foundry store my AWS credentials? 
  * __A__:  We only use your AWS keys while creating your Cloud Foundry instance. We don't store your AWS keys, and they are deleted upon completion. 
  
- - - - - - - - -

  * __Q__: How much does it cost to run Cloud Foundry on AWS? 
  * __A__:  While Try Cloud Foundry does not charge anything to setup a Cloud Foundry instance, Amazon will charge by the hour. We estimate that running Cloud Foundry on AWS will cost about 0.28 cents/hour, however, you can calculate a cost using [this](http://aws.amazon.com/ec2/pricing/) document. We are not liable for any damages or costs you incur while using Try Cloud Foundry. 
  
- - - - - - - - -

  * __Q__: I'm having problems with AWS Elastic IPs 
  * __A__: If you received an email from Try Cloud Foundry explaining that an AWS Elasic IP could not be allocated or added, you'll need to free up an AWS Elasic IP. See the [AWS documentation](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/elastic-ip-addresses-eip.html) for more info.
  
- - - - - - - - -

  * __Q__: I'm having a problem with Try Cloud Foundry, how do I get support? 
  * __A__:  While we don't officially offer support for Try Cloud Foundry, you can send an email to [and we'll see what we can do](mailto:trycf@starkandwayne.com). 
  
- - - - - - - - -

## Getting help

If for some reason you are having problems getting Cloud Foundry up and running, you can email us at <trycf@starkandwayne.com>.

You can also ask questions on the [Cloud Foundry Developers](https://groups.google.com/a/cloudfoundry.org/group/vcap-dev/topics) mailing list.

The Cloud Foundry site has some technical overview videos in the [Learn](http://www.cloudfoundry.org/learn/index.html) section.


---

***Thanks for trying Cloud Foundry!***

The Stark and Wayne team
