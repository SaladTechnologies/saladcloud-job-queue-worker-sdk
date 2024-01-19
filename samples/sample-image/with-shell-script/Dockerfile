FROM foobar

WORKDIR /app

# copy previously downloaded salad-queue-client
COPY ./salad-job-queue-worker /usr/local/bin/
# copy the shell script
COPY run.sh ./

CMD ["/app/run.sh"]
