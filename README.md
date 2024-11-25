To run program:

Open 3 terminals.

Terminal 1:
cd server
run 'go run serverside.go'
(to simulate crash, press CTRL+C in the terminal window)

Terminal 2:
cd backupserver
run 'go run backupserver.go'

Terminal 3:
cd clients
run 'go run clientside.go'



None of the processes close by themselves, so when done with program, press CTRL+C in all of the above terminal windows.
