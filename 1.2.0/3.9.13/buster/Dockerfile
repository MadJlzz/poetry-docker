FROM python:3.9.13-buster

LABEL maintainer="Julien Klaer"

RUN pip install poetry==1.2.0 && \
    poetry config virtualenvs.in-project true

CMD ["poetry", "--version"]
