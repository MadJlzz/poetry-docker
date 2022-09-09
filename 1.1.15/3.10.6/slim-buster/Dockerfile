# This Dockerfile was generated.
# Do not manually update this file; changes will be lost.
FROM python:3.10.6-slim-buster

LABEL maintainer="Julien Klaer"

RUN pip install poetry==1.1.15 && \
    poetry config virtualenvs.in-project true

CMD ["poetry", "--version"]
