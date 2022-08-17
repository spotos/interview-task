FROM google/cloud-sdk:alpine

ENV PORT=8085

RUN apk add --no-cache openjdk8-jre
RUN gcloud --quiet components install beta pubsub-emulator

RUN mkdir -p /opt/pubsub
RUN echo "" >> /opt/pubsub/env.yaml

CMD gcloud beta emulators pubsub start --data-dir=/opt/pubsub --host-port=0.0.0.0:${PORT} --project=local --log-http --verbosity=debug --user-output-enabled
