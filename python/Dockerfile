ARG PYTHON_VERSION=3.12.3
FROM python:${PYTHON_VERSION}-slim

ENV PYTHONDONTWRITEBYTECODE=1

WORKDIR /app

# Install dependencies for Chrome and Xvfb
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget \
    unzip \
    xvfb \
    gnupg \
    && apt-get clean

# Install Chrome
RUN --mount=type=cache,target=/root/.cache/chrome \
    wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb -P /root/.cache/chrome && \
    apt install -y /root/.cache/chrome/google-chrome-stable_current_amd64.deb && \
    rm /root/.cache/chrome/google-chrome-stable_current_amd64.deb && \
    apt-get clean

# Install ChromeDriver
RUN --mount=type=cache,target=/root/.cache/chromedriver \
    CHROMEDRIVER_VERSION=$(wget -qO- https://chromedriver.storage.googleapis.com/LATEST_RELEASE) && \
    wget -q -O /root/.cache/chromedriver/chromedriver.zip https://chromedriver.storage.googleapis.com/${CHROMEDRIVER_VERSION}/chromedriver_linux64.zip && \
    unzip /root/.cache/chromedriver/chromedriver.zip -d /usr/local/bin/ && \
    rm /root/.cache/chromedriver/chromedriver.zip

COPY requirements.txt /app/

# Install Python dependencies
RUN --mount=type=cache,target=/root/.cache/pip \
    --mount=type=bind,source=requirements.txt,target=requirements.txt \
    python -m pip install -r requirements.txt

COPY . /app

# Set the environment variable to use Xvfb
ENV DISPLAY=:99

# Start Xvfb and run the application
CMD ["sh", "-c", "Xvfb :99 -ac & python main.py"]
