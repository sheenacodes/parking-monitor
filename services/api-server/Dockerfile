FROM python:3.12-slim

WORKDIR /project

RUN useradd -m -u 1000 appuser

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY . /project/

RUN chown -R appuser:appuser /project
RUN  chmod 644 /project/log/log.txt
USER appuser

CMD ["python", "-m", "app.main"]
