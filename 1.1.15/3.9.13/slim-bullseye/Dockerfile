FROM python:3.9.13-slim-bullseye

LABEL maintainer="Julien Klaer"

RUN pip install poetry==1.1.15 && \
    poetry config virtualenvs.in-project true

CMD ["poetry", "--version"]
