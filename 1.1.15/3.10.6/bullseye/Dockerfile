FROM python:3.10.6-bullseye

LABEL maintainer="Julien Klaer"

RUN pip install poetry==1.1.15 && \
    poetry config virtualenvs.in-project true

CMD ["poetry", "--version"]
