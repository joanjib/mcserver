# mcserver
This project it's in early state of develpment

The purpose of this project is to serve as a general server for the development of mission critical Server Apps.

Characteristics of the project:

* SQL as the main development programming language: SQL files (actions, functionalities) in the server are exposed to clients (mobile Apps, web apps, IoT,...) using a highly concurrent system through Websockets, allowing bidirectional communication.

* All SQL actions are SERIALIZED. This means the system is less prone to execution conflicts between SQL statements. Also because of this, to improve performance of the SQL serialized statements, the server allows fine tuning on when transaction have to be committed to disk, allowing to do most of the computation in RAM.

* The system is not oriented to server apps with huge amount of requests per second. It's focused in solving problems, in "closed environments", like hospitals, nursing homes, stores, etc.

* The server sends an ACK to the clients on each operation. There are two ways this ACK can be sent:
  - First sending the resulting data: in this way the client can pre-allocate resources, for example. And when data is actually saved to disk in the server, it's sent the ACK
  - Or can be sent the resulting data and ACK at the same time (for example when nothing can be done meanwhile the computation and the time data is saved to disk)

* Why SQL? The answer is taken form the web site (sqlite.org): "SQL is a very high-level language. A few lines of SQL can replace hundreds or thousands of lines of procedural code. SQL thus reduces the amount of work needed to develop and maintain the application, and thereby helps to reduce the number of bugs in the application."

* The server allows most the time to develop using only SQL, but the server is being programmed using Golang.

* The server automatically compiles the SQL statements to gain some performance, when usually this process is done like this (source sqlite.org): "The best way to understand how SQL database engines work is to think of SQL as a programming language, not as a "query language". Each SQL statement is a separate program. Applications construct SQL program source files and send them to the database engine. The database engine compiles the SQL source code into executable form, runs that executable, then sends the result back to the application."

In the server I present, most the previous steps are not necessary (and performance is gained)
