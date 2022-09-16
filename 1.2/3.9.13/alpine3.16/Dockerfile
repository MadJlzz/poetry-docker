# This Dockerfile was generated.
# Do not manually update this file; changes will be lost.
FROM python:3.9.13-alpine3.16

LABEL maintainer="Julien Klaer"
RUN apk add gcc libc-dev
RUN pip install poetry==1.2.1 && \
    poetry config virtualenvs.in-project true

CMD ["poetry", "--version"]
