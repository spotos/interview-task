FROM google/cloud-sdk:emulators

ENV PORT=8085
ENV VERBOSITY=warning

RUN mkdir -p /opt/pubsub
RUN echo "" >> /opt/pubsub/env.yaml

CMD gcloud beta emulators pubsub start --data-dir=/opt/pubsub --host-port=0.0.0.0:${PORT} --project=local --log-http --verbosity=${VERBOSITY} --user-output-enabled
