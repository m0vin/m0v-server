###


go build -v
sudo killall finisher
sudo ./runwdb.sh

~/dev/m0v/run.sh 2>&1 | ./funnel
