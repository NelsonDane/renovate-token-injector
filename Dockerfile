# Taken from: https://www.joshkasuboski.com/posts/distroless-python-uv/

FROM ghcr.io/astral-sh/uv:bookworm-slim AS builder
ENV UV_COMPILE_BYTECODE=1 \
    UV_LINK_MODE=copy \
    UV_PYTHON_INSTALL_DIR=/python \
    UV_PYTHON_PREFERENCE=only-managed

# Install Python before the project for caching
RUN uv python install 3.14

# uv sync with caching
WORKDIR /app
RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=uv.lock,target=uv.lock \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml \
    uv sync --frozen --no-install-project --no-dev --no-editable
COPY . /app
RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --frozen --no-dev --no-editable

FROM gcr.io/distroless/cc-debian12 AS final

# Copy the Python version
COPY --from=builder --chown=python:python /python /python

# Copy the application from the builder
WORKDIR /app
COPY --from=builder --chown=app:app /app/.venv /app/.venv
COPY --chown=app:app main.py /app/main.py

# Add venv to PATH
ENV PATH="/app/.venv/bin:$PATH"

# Run Entrypoint
CMD ["python", "main.py"]
