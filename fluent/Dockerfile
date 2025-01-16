FROM fluent/fluentd:v1.16.6-1.0
LABEL authors="ksCaesar"
USER root

RUN fluent-gem install \
    fluent-plugin-multi-format-parser \
    fluent-plugin-filter \
    fluent-plugin-record-modifier \
    fluent-plugin-bigquery

RUN mkdir -p /var/log/acs_go && chown fluent:fluent /var/log/acs_go
USER fluent
COPY ./keys/bq_svc_key.json /home/ec2-user/acs_go/keys/bq_svc_key.json

# docker build -t x246libra/fluentd:v1.16.6 .

# verify config
# docker run --rm -v $(pwd)/builds:/fluentd/etc/ x246libra/fluentd:v1.16.6 fluentd --dry-run

# local test config
# docker run --rm -v $(pwd)/builds:/fluentd/etc/ -v $(pwd)/client.log:/var/log/acs_go/client.log -v $(pwd)/bq_svc_key.json:/home/ec2-user/acs_go/keys/bq_svc_key.json x246libra/fluentd:v1.16.6
# docker run --rm -v $(pwd)/builds:/fluentd/etc/ -v $(pwd)/client.log:/var/log/acs_go/client.log x246libra/fluentd:v1.16.6
