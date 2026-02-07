FROM python:3.11-slim

WORKDIR /app
COPY device.py .

RUN pip install flask

EXPOSE 5001
ENV PYTHONUNBUFFERED=1
CMD ["python", "device.py"]
