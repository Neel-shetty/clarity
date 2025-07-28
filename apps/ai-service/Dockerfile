FROM ghcr.io/astral-sh/uv:python3.13-alpine
WORKDIR /app
COPY pyproject.toml uv.lock /app/
RUN uv sync --locked --compile-bytecode
COPY . /app
RUN uv sync --locked --compile-bytecode
CMD ["uv", "run", "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
