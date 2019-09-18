# go-playground

CRUD user management  following mvc pattern using Go http server and text template standard packages, and some all those sql operations without the use of gorm.

Rest api service for user management using Go http, gorilla router, and jinzhou Gorm for db operations, along with a configmap for the app configuration from environment that’s also to be imported into kubernetes environment. And built it with a dockerfile, to be deployed with a deployment file to local kubernetes, with a mysql instance service directly spawn up from kubernetes using an existing mysql container image from the docker hub. Finally, all steps of building deployment are packed into a Makefile to get it ready to be shipped.

Rest api service projects to play with go-kit endpoint and httptransport, and logging, proxying, instrumenting middleware for service, and pipkin, ratelimiter, circuit breaker middleware for both server and client side of service endpoints.

Implemented a web application using text template to play with JWT authentication with go crypto/bcrpt projected password that is persisted into db via gorm db operations, making use of session cookies, with alternative redis session persistence approach as well. this app would allow user to sign up, sign in, sign out, forgotten password, keeping session alive etc. The password resetting link is sent via email (gomail package), and new password is encrypted and persisted back into db.

poker play app in both cli and webserver versions, and also allows user to play cli version in a web browser as well using websocket for playing stats communication between the browser and the server.

clockface app draws an analog clock face in SVG format, ticking by second and dispatched to the browser via websocket communication

boxoffice app allows user to reserve ticket within a certain time limit to enable them to proceed to pay by card at stripe checkout. it's using a sql database for persistence, also a jwt token for keepting tracking of user activities.
