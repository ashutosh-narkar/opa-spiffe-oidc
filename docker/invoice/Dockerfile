FROM golang:1.7
 
COPY invoice-app /usr/local/bin/invoice-app
COPY conf/start_service.sh /usr/local/bin/start_service.sh
RUN chmod u+x /usr/local/bin/start_service.sh
ENTRYPOINT /usr/local/bin/start_service.sh

EXPOSE 5000
