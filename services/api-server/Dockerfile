FROM python:3.12-slim

WORKDIR /app

RUN useradd -m -u 1000 appuser

COPY requirements.txt .

COPY ./app /app/app

RUN chown -R appuser:appuser /app

USER appuser
#ENV PATH="/root/.local/bin:${PATH}"

RUN pip install -r requirements.txt

CMD ["uvicorn", "app.main:app"]