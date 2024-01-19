# syntax=docker/dockerfile:1

FROM python:3.11-slim

EXPOSE 80/tcp

WORKDIR /app

COPY ./requirements.txt *.py ./

RUN pip install --no-cache-dir --upgrade -r /app/requirements.txt

CMD ["uvicorn", "main:app", "--host", "*", "--port", "80"]
