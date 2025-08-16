# analabit (core)

The missing anal_ytics service for abiturients.

This repository contains the backend core powering analabit.ru: data ingestion, normalization, calculations, and service APIs. It’s a backend-only project that other layers build upon.

## What’s here

- Core: reusable domain logic and data processing (see `core/`, `ent/`, and supporting packages under `core/*`).
- Services: production services built on the core (see `service/aggregator`, `service/api`, `service/idmsu`, `service/producer`).
- CLI: tooling for maintenance and data operations (see `cli/`).

Looking for the UI? There’s only the web interface for using the backend, built above the analabit core.

➡️ Frontend (Web UI): https://github.com/trueegorletov/analabit-webui

## Showcase

There is some images of how does it feel when using through web UI.

[![Watch demo GIF on Imgur](https://img.shields.io/badge/Watch%20demo-Imgur-1bb76e?logo=imgur&logoColor=white)](https://imgur.com/ngIR9nI)

![Home page](https://raw.githubusercontent.com/trueegorletov/analabit-webui/refs/heads/main/.github/assets/1.png)
![Educational program dashboard](https://raw.githubusercontent.com/trueegorletov/analabit-webui/refs/heads/main/.github/assets/3.png)
![Abiturient details popup](https://raw.githubusercontent.com/trueegorletov/analabit-webui/refs/heads/main/.github/assets/2.png)

