FROM confluentinc/cp-kafka:latest

COPY prom-jmx-agent-config.yml /etc/jmx-agent/prom-jmx-agent-config.yml
COPY jmx_prometheus_javaagent-1.0.1.jar /usr/app/jmx_prometheus_javaagent.jar
ENV KAFKA_OPTS="-javaagent:/usr/app/jmx_prometheus_javaagent.jar=9102:/etc/jmx-agent/prom-jmx-agent-config.yml"
