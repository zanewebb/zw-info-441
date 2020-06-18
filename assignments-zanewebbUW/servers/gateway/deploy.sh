./build.sh
docker push zanewebb/zanewebbuw
cat serverdeployV2.sh | sudo ssh -i 441classroom.pem ec2-user@18.221.123.38
