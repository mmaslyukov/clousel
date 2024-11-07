
i=0
while [ 1 ]; do
  i=$((i+1));
  echo $i;
  cat ./heartbeat.json | sed  "s/\"SequenceNum\": ./\"SequenceNum\": $i/g" > /tmp/hb
  cat /tmp/hb
  mosquitto_pub -h 192.168.0.150 -p 1883 -t '/clousel/carousel/550e8400-e29b-41d4-a716-446655440000' -f /tmp/hb
  sleep 30;
done
