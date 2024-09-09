FROM python:3.12-slim

WORKDIR /project

RUN useradd -m -u 1000 appuser

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY ./app /project/app
COPY ./log /project/log


RUN ls -la /project  # Debugging step to list files

RUN chown -R appuser:appuser /project

USER appuser

CMD ["python", "-m", "app.main"]
