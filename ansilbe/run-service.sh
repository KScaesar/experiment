#!/bin/bash

set -e

BUILD_NO=$1
ENV=$2

if [ -z "$BUILD_NO" ] || [ -z "$ENV" ]; then
  echo "Usage: $0 <BUILD_NO> <ENV>"
  exit 1
fi

if ! [[ "$BUILD_NO" =~ ^[0-9]+$ ]]; then
  echo "Error: BUILD_NO must be a number."
  exit 1
fi

if [[ "$ENV" != "production" && "$ENV" != "staging" ]]; then
  echo "Error: ENV must be either 'production' or 'staging'."
  exit 1
fi

SVC=bcs
SVC_NO=$SVC-$BUILD_NO
ARCHIVE=$SVC.tar.gz

Deploy_DIR=~/deployment/
Download_DIR=~/deployment/files

if [ ! -d "$Deploy_DIR/$SVC_NO" ]; then
  cd $Download_DIR
  # 因為 bcs 起初只有 production 環境, 後續才申請到 staging 環境進行測試, 所以 jenkins 目前只有 bcs-production
  wget http://cdibuild.tw.svc.litv.tv:8080/job/$SVC-production/"$BUILD_NO"/artifact/$ARCHIVE
  tar -xzvf $ARCHIVE
  mv $SVC $Deploy_DIR/"$SVC_NO"
  rm -f $ARCHIVE
fi

cd $Deploy_DIR
if [ -d $SVC ]; then
  rm -rf $SVC
fi
ln -s "$SVC_NO" $SVC

echo "move file for .env and keys"
mkdir -p $SVC/keys
mv -f files/gcp_credential.json $SVC/keys
mv -f files/star.svc.litv.tv.crt.pem $SVC/keys
mv -f files/star.svc.litv.tv.key.pem $SVC/keys
mv -f files/.env."$ENV" $SVC/.env

echo "restart bcs service"
cd $SVC
if [ "$(docker ps -aq -f name=prefect-server)" ]; then
  echo "prefect deployment pause-schedule"
  docker exec prefect-server prefect deployment pause-schedule daily-predict-traffics/daily-predict-traffics
  docker exec prefect-server prefect deployment pause-schedule pacing-calculate-every-n-minutes/pacing-calculate-every-n-minutes
fi
docker compose -p bcs --profile prefect -f "./deploy/docker-compose.yml" up -d
Jenkins_Number=$BUILD_NO docker compose -p bcs --profile cron -f "./deploy/docker-compose.yml" up -d --build
