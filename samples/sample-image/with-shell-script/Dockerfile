FROM foobar

WORKDIR /app

# copy previously downloaded salad-queue-client
COPY ./salad-queue-client /usr/local/bin/
# copy the shell script
COPY run.sh ./

CMD ["/app/run.sh"]
