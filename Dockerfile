FROM ghcr.io/astral-sh/uv:alpine

WORKDIR /app

COPY . /app
RUN uv sync --frozen --no-dev --no-editable
ENV PATH="/app/.venv/bin:$PATH"

ENTRYPOINT ["python", "main.py"]
