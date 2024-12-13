To run program:

Open 3 terminals.

Terminal 1:
cd server
run 'go run serverside.go'
(crashes after 10 seconds - this is a feature, read report for details)

Terminal 2:
cd backupserver
run 'go run backupserver.go'

Terminal 3:
cd clients
run 'go run clientside.go'



When done with program (and client has received 'Auction is over message'), press CTRL+C in backupserver and clientside  terminal windows to terminate the program for good.

Inspect logs with 'cat /tmp/logs.txt' and 'cat /tmp/logstime.txt'.
