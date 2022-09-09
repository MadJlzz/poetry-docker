FROM python:3.9.13-alpine3.16

LABEL maintainer="Julien Klaer"

RUN pip install poetry==1.2.0 && \
    poetry config virtualenvs.in-project true

CMD ["poetry", "--version"]
