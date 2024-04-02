spurtCMS GraphQL APIs that exposes spurtCMS Admin side contents and how developers can harness its potential to consume efficient, self-describing APIs for their CMS website by adding or modifying api fields without impacting existing flow.

 

# Building Blocks of spurtCMS GraphQL APIs are,

### GQL Gen library: 

To simplify the process of creating GraphQL servers in Go Lang by automating the generation of schemas, types, resolvers, and more.

### Custom Go Packages: 

To enhance the functionality of GraphQL Playground, spurtCMS team has developed and integrated their custom packages like pkgcore and pkgcontent.

Playground Interface: 

To present an intitutive API demo console, spurtCMS GraphQL project has integrated React components served through a CDN for HTML content delivery.

#### Installation Steps:

 

### Step 1:

Ensure the pre-requisites of Go Environment is available ready in your system. Check our pre-requisites guideline for more details.

spurtCMS Admin application serves as the content source for spurtCMS GraphQL APIs and hence setup the admin application and execute it. Refer the spurtCMS admin setup detailed in this link.  https://dev.spurtcms.com/documentation/cms-admin

### Step 2:
To initiate the setup process, the first step involves cloning the spurtCMS Template project structure from the GIT repository. https://github.com/spurtcms/spurtcms-graphql

```
$git clone https://github.com/spurtcms/spurtcms-graphql.git
```
The cloned repository should have all the project files associated with spurtCMS GraphQL project.

### Step 3:
Now  it's time to setup the PostgreSQL database, the one which is a replica of spurtCMS Admin application as mentioned in pre-requisites (Step 1) above.

Locate the .env file of the GraphQL project folder and configure it with the details of newly imported admin database such as database name, user name, password etc

# PostgreSQL Database Configuration

```
DB_HOST=localhost 
DB_PORT=5432 
DB_NAME=your_database_name 
DB_USER=your_database_user 
DB_PASSWORD=your_database_password 
DB_SSL_MODE=disable
```
Also, configure the admin application's base url i .env which is needed for GraphQL API to output the path of Media files.
```
#DOMAIN_URL ='https://demo.spurtcms.com/' 
##Example PORT='8081'
```
### Step 4:
Now it's time to start the GraphqL application. Open the terminal from the GraphQL project folder and execute the following command.
```
$ go run server.go
```
This command initiates the Go program, installs the dependencies such as GQL Gen library, spurtCMS custom packages such as pkgcore and pkgcontent. 

### Step 5:
Your spurtCMS GraphQL API is ready to use now !! You can use the playground interface to check the sample request and response structure of each GraphQL endpoint. https://{your-domain-name}/play

### Conclusion:

As GraphQL continues to gain popularity within the developer community, mastering its concepts and tools like GQL Gen opens up a world of possibilities for building modern, data-driven applications. By following the steps outlined in this article, developers can easily integrate spurtCMS GraphQL APIs in their front-end applications to retrive content from spurtCMS admin application.
