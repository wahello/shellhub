USER1=edu
USER2=edu2
USER3=edu3

TENANT=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
DEVICE=434bd6bc8c1830df1092dc518bc16b55a4cb3c6b2fcc63642d65a06f5a5de716
SESSION=b479ab90dc57be45f9e0621b474cc918e9917823c3df9bc2346bb054d7fc212c
PUBLIC_KEY=31:62:1e:27:de:00:9c:e5:de:2f:3d:17:30:0e:46:2f

for USER in $USER3; do
echo "request for user $USER"
TOKEN=`http post http://localhost/api/login username="$USER" password="211250" | jq  -r .token`
http --check-status -q patch http://localhost/api/devices/$DEVICE name=bbbbbbbbb "Authorization: Bearer $TOKEN" && echo -ne pass || echo -ne error; echo " rename device"
done
