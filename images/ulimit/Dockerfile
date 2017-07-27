FROM ubuntu:16.04

RUN apt-get update -y && apt-get upgrade -y
COPY ulimitabuser /usr/bin/app
COPY initscript.sh /initscript.sh
RUN chmod +x /initscript.sh
ENTRYPOINT ["/initscript.sh"]
