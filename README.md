# simplesocial
Basic social networking model created in Golang, with a MySQL database.

**First installation guide**

-> Download directory to '$GOPATH/src'

-> Download all external dependencies

-> Create link of folder "$GOPATH/src/simplesocial/simplesocialtmp" in directory '$GOPATH/bin'

-> Rename link "simplesocialtmp"

**->** Create necessary databases. I will add the SQL files in upcoming commits.
-> UPDATE: .sql file has been added in /MySQLFiles.

-> *Run in terminal:*

$ go install simplesocial

**Execute created executable file**

-> *Run in terminal:*

$ cd $GOPATH/bin

$ ./simplesocial
