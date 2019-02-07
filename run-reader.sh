echo "starting statusd...."
statusd -c /etc/gs-bots/_assets/ethdenver.json &

echo "waiting for statusd...."
sleep 2

echo "starting chanreader...."
chanreader
