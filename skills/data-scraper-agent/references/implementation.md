# Data Scraper Agent — Reference Implementation (Python)

A complete, working stack: `requests`/`BeautifulSoup` for scraping, Gemini Flash (free tier) for enrichment, GitHub Actions cron for scheduling, and Notion as the storage backend. Read this after `SKILL.md` when you're ready to wire up a real source.

## Full Directory Structure

```
my-agent/
├── config.yaml              # User customises this (keywords, filters, preferences)
├── profile/
│   └── context.md           # User context the AI uses (resume, interests, criteria)
├── scraper/
│   ├── __init__.py
│   ├── main.py              # Orchestrator: scrape → enrich → store
│   ├── filters.py           # Rule-based pre-filter (fast, before AI)
│   └── sources/
│       ├── __init__.py
│       └── source_name.py   # One file per data source
├── ai/
│   ├── __init__.py
│   ├── client.py            # Gemini REST client with model fallback
│   ├── pipeline.py          # Batch AI analysis
│   └── memory.py            # Learn from user feedback
├── storage/
│   ├── __init__.py
│   └── notion_sync.py       # Or sheets_sync.py / supabase_sync.py
├── data/
│   └── feedback.json        # User decision history (auto-updated)
├── .env.example
├── setup.py                 # One-time DB/schema creation
├── enrich_existing.py       # Backfill AI scores on old rows
├── requirements.txt
└── .github/
    └── workflows/
        └── scraper.yml      # GitHub Actions schedule
```

## Scraper Source Template

```python
# scraper/sources/my_source.py
"""
[Source Name] — scrapes [what] from [where].
Method: [REST API / HTML scraping / RSS feed]
"""
import requests
from bs4 import BeautifulSoup
from datetime import datetime, timezone
from scraper.filters import is_relevant

HEADERS = {
    "User-Agent": "Mozilla/5.0 (compatible; research-bot/1.0)",
}


def fetch() -> list[dict]:
    """
    Returns a list of items with consistent schema.
    Each item must have at minimum: name, url, date_found.
    """
    results = []

    # ---- REST API source ----
    resp = requests.get("https://api.example.com/items", headers=HEADERS, timeout=15)
    if resp.status_code == 200:
        for item in resp.json().get("results", []):
            if not is_relevant(item.get("title", "")):
                continue
            results.append(_normalise(item))

    return results


def _normalise(raw: dict) -> dict:
    """Convert raw API/HTML data to the standard schema."""
    return {
        "name": raw.get("title", ""),
        "url": raw.get("link", ""),
        "source": "MySource",
        "date_found": datetime.now(timezone.utc).date().isoformat(),
        # add domain-specific fields here
    }
```

**HTML scraping pattern:**
```python
soup = BeautifulSoup(resp.text, "lxml")
for card in soup.select("[class*='listing']"):
    title = card.select_one("h2, h3").get_text(strip=True)
    link = card.select_one("a")["href"]
    if not link.startswith("http"):
        link = f"https://example.com{link}"
```

**RSS feed pattern:**

Use `defusedxml` instead of the stdlib `xml.etree.ElementTree` — the stdlib XML parser is vulnerable to XXE and billion-laughs attacks against untrusted feed content:

```python
from defusedxml import ElementTree as ET  # pip install defusedxml
root = ET.fromstring(resp.text)
for item in root.findall(".//item"):
    title = item.findtext("title", "")
    link = item.findtext("link", "")
```

**Paginated API pattern:**
```python
page = 1
while True:
    resp = requests.get(url, params={"page": page, "limit": 50}, timeout=15)
    data = resp.json()
    items = data.get("results", [])
    if not items:
        break
    for item in items:
        results.append(_normalise(item))
    if not data.get("has_more"):
        break
    page += 1
```

**JS-rendered pages (Playwright):**
```python
from playwright.sync_api import sync_playwright

with sync_playwright() as p:
    browser = p.chromium.launch()
    page = browser.new_page()
    page.goto(url)
    page.wait_for_selector(".listing")
    html = page.content()
    browser.close()

soup = BeautifulSoup(html, "lxml")
```

## AI Client with Model Fallback

```python
# ai/client.py
import os, json, time, requests

_last_call = 0.0

MODEL_FALLBACK = [
    "gemini-2.0-flash-lite",
    "gemini-2.0-flash",
    "gemini-2.5-flash",
    "gemini-flash-lite-latest",
]


def generate(prompt: str, model: str = "", rate_limit: float = 7.0) -> dict:
    """Call Gemini with auto-fallback on 429. Returns parsed JSON or {}."""
    global _last_call

    api_key = os.environ.get("GEMINI_API_KEY", "")
    if not api_key:
        return {}

    elapsed = time.time() - _last_call
    if elapsed < rate_limit:
        time.sleep(rate_limit - elapsed)

    models = [model] + [m for m in MODEL_FALLBACK if m != model] if model else MODEL_FALLBACK
    _last_call = time.time()

    for m in models:
        url = f"https://generativelanguage.googleapis.com/v1beta/models/{m}:generateContent?key={api_key}"
        payload = {
            "contents": [{"parts": [{"text": prompt}]}],
            "generationConfig": {
                "responseMimeType": "application/json",
                "temperature": 0.3,
                "maxOutputTokens": 2048,
            },
        }
        try:
            resp = requests.post(url, json=payload, timeout=30)
            if resp.status_code == 200:
                return _parse(resp)
            if resp.status_code in (429, 404):
                time.sleep(1)
                continue
            return {}
        except requests.RequestException:
            return {}

    return {}


def _parse(resp) -> dict:
    try:
        text = (
            resp.json()
            .get("candidates", [{}])[0]
            .get("content", {})
            .get("parts", [{}])[0]
            .get("text", "")
            .strip()
        )
        if text.startswith("```"):
            text = text.split("\n", 1)[-1].rsplit("```", 1)[0]
        return json.loads(text)
    except (json.JSONDecodeError, KeyError):
        return {}
```

## Batch Enrichment Pipeline

```python
# ai/pipeline.py
import json
import yaml
from pathlib import Path
from ai.client import generate

def analyse_batch(items: list[dict], context: str = "", preference_prompt: str = "") -> list[dict]:
    """Analyse items in batches. Returns items enriched with AI fields."""
    config = yaml.safe_load((Path(__file__).parent.parent / "config.yaml").read_text())
    model = config.get("ai", {}).get("model", "gemini-2.5-flash")
    rate_limit = config.get("ai", {}).get("rate_limit_seconds", 7.0)
    min_score = config.get("ai", {}).get("min_score", 0)
    batch_size = config.get("ai", {}).get("batch_size", 5)

    batches = [items[i:i + batch_size] for i in range(0, len(items), batch_size)]
    print(f"  [AI] {len(items)} items → {len(batches)} API calls")

    enriched = []
    for i, batch in enumerate(batches):
        print(f"  [AI] Batch {i + 1}/{len(batches)}...")
        prompt = _build_prompt(batch, context, preference_prompt, config)
        result = generate(prompt, model=model, rate_limit=rate_limit)

        analyses = result.get("analyses", [])
        for j, item in enumerate(batch):
            ai = analyses[j] if j < len(analyses) else {}
            if ai:
                score = max(0, min(100, int(ai.get("score", 0))))
                if min_score and score < min_score:
                    continue
                enriched.append({**item, "ai_score": score, "ai_summary": ai.get("summary", ""), "ai_notes": ai.get("notes", "")})
            else:
                enriched.append(item)

    return enriched


def _build_prompt(batch, context, preference_prompt, config):
    priorities = config.get("priorities", [])
    items_text = "\n\n".join(
        f"Item {i+1}: {json.dumps({k: v for k, v in item.items() if not k.startswith('_')})}"
        for i, item in enumerate(batch)
    )

    return f"""Analyse these {len(batch)} items and return a JSON object.

# Items
{items_text}

# User Context
{context[:800] if context else "Not provided"}

# User Priorities
{chr(10).join(f"- {p}" for p in priorities)}

{preference_prompt}

# Instructions
Return: {{"analyses": [{{"score": <0-100>, "summary": "<2 sentences>", "notes": "<why this matches or doesn't>"}} for each item in order]}}
Be concise. Score 90+=excellent match, 70-89=good, 50-69=ok, <50=weak."""
```

## Feedback Learning

```python
# ai/memory.py
"""Learn from user decisions to improve future scoring."""
import json
from pathlib import Path

FEEDBACK_PATH = Path(__file__).parent.parent / "data" / "feedback.json"


def load_feedback() -> dict:
    if FEEDBACK_PATH.exists():
        try:
            return json.loads(FEEDBACK_PATH.read_text())
        except (json.JSONDecodeError, OSError):
            pass
    return {"positive": [], "negative": []}


def save_feedback(fb: dict):
    FEEDBACK_PATH.parent.mkdir(parents=True, exist_ok=True)
    FEEDBACK_PATH.write_text(json.dumps(fb, indent=2))


def build_preference_prompt(feedback: dict, max_examples: int = 15) -> str:
    """Convert feedback history into a prompt bias section."""
    lines = []
    if feedback.get("positive"):
        lines.append("# Items the user LIKED (positive signal):")
        for e in feedback["positive"][-max_examples:]:
            lines.append(f"- {e}")
    if feedback.get("negative"):
        lines.append("\n# Items the user SKIPPED/REJECTED (negative signal):")
        for e in feedback["negative"][-max_examples:]:
            lines.append(f"- {e}")
    if lines:
        lines.append("\nUse these patterns to bias scoring on new items.")
    return "\n".join(lines)
```

Integration: after each run, query your storage backend for items with positive/negative status and call `save_feedback()` with the extracted patterns (implement this as a separate `feedback_sync.py`).

## Notion Storage

```python
# storage/notion_sync.py
import os
from notion_client import Client
from notion_client.errors import APIResponseError

_client = None

def get_client():
    global _client
    if _client is None:
        _client = Client(auth=os.environ["NOTION_TOKEN"])
    return _client

def get_existing_urls(db_id: str) -> set[str]:
    """Fetch all URLs already stored — used for deduplication."""
    client, seen, cursor = get_client(), set(), None
    while True:
        resp = client.databases.query(database_id=db_id, page_size=100, **{"start_cursor": cursor} if cursor else {})
        for page in resp["results"]:
            url = page["properties"].get("URL", {}).get("url", "")
            if url: seen.add(url)
        if not resp["has_more"]: break
        cursor = resp["next_cursor"]
    return seen

def push_item(db_id: str, item: dict) -> bool:
    """Push one item to Notion. Returns True on success."""
    props = {
        "Name": {"title": [{"text": {"content": item.get("name", "")[:100]}}]},
        "URL": {"url": item.get("url")},
        "Source": {"select": {"name": item.get("source", "Unknown")}},
        "Date Found": {"date": {"start": item.get("date_found")}},
        "Status": {"select": {"name": "New"}},
    }
    if item.get("ai_score") is not None:
        props["AI Score"] = {"number": item["ai_score"]}
    if item.get("ai_summary"):
        props["Summary"] = {"rich_text": [{"text": {"content": item["ai_summary"][:2000]}}]}
    if item.get("ai_notes"):
        props["Notes"] = {"rich_text": [{"text": {"content": item["ai_notes"][:2000]}}]}

    try:
        get_client().pages.create(parent={"database_id": db_id}, properties=props)
        return True
    except APIResponseError as e:
        print(f"[notion] Push failed: {e}")
        return False

def sync(db_id: str, items: list[dict]) -> tuple[int, int]:
    existing = get_existing_urls(db_id)
    added = skipped = 0
    for item in items:
        if item.get("url") in existing:
            skipped += 1; continue
        if push_item(db_id, item):
            added += 1; existing.add(item["url"])
        else:
            skipped += 1
    return added, skipped
```

## Orchestrator

```python
# scraper/main.py
import os, sys, yaml
from pathlib import Path
from dotenv import load_dotenv

load_dotenv()

from scraper.sources import my_source          # add your sources

# NOTE: This example uses Notion. If storage.provider is "sheets" or "supabase",
# replace this import with storage.sheets_sync or storage.supabase_sync and update
# the env var and sync() call accordingly.
from storage.notion_sync import sync

SOURCES = [
    ("My Source", my_source.fetch),
]

def ai_enabled():
    return bool(os.environ.get("GEMINI_API_KEY"))

def main():
    config = yaml.safe_load((Path(__file__).parent.parent / "config.yaml").read_text())
    provider = config.get("storage", {}).get("provider", "notion")

    if provider == "notion":
        db_id = os.environ.get("NOTION_DATABASE_ID")
        if not db_id:
            print("ERROR: NOTION_DATABASE_ID not set"); sys.exit(1)
    else:
        print(f"ERROR: provider '{provider}' not yet wired in main.py"); sys.exit(1)

    all_items = []

    for name, fetch_fn in SOURCES:
        try:
            items = fetch_fn()
            print(f"[{name}] {len(items)} items")
            all_items.extend(items)
        except Exception as e:
            print(f"[{name}] FAILED: {e}")

    # Deduplicate by URL
    seen, deduped = set(), []
    for item in all_items:
        if (url := item.get("url", "")) and url not in seen:
            seen.add(url); deduped.append(item)

    print(f"Unique items: {len(deduped)}")

    if ai_enabled() and deduped:
        from ai.memory import load_feedback, build_preference_prompt
        from ai.pipeline import analyse_batch

        feedback = load_feedback()
        preference = build_preference_prompt(feedback)
        context_path = Path(__file__).parent.parent / "profile" / "context.md"
        context = context_path.read_text() if context_path.exists() else ""
        deduped = analyse_batch(deduped, context=context, preference_prompt=preference)
    else:
        print("[AI] Skipped — GEMINI_API_KEY not set")

    added, skipped = sync(db_id, deduped)
    print(f"Done — {added} new, {skipped} existing")

if __name__ == "__main__":
    main()
```

## GitHub Actions Schedule

```yaml
# .github/workflows/scraper.yml
name: Data Scraper Agent

on:
  schedule:
    - cron: "0 */3 * * *"  # every 3 hours — adjust to your needs
  workflow_dispatch:        # allow manual trigger

permissions:
  contents: write   # required for the feedback-history commit step

jobs:
  scrape:
    runs-on: ubuntu-latest
    timeout-minutes: 20

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-python@v5
        with:
          python-version: "3.11"
          cache: "pip"

      - run: pip install -r requirements.txt

      # Uncomment if Playwright is enabled in requirements.txt
      # - name: Install Playwright browsers
      #   run: python -m playwright install chromium --with-deps

      - name: Run agent
        env:
          NOTION_TOKEN: ${{ secrets.NOTION_TOKEN }}
          NOTION_DATABASE_ID: ${{ secrets.NOTION_DATABASE_ID }}
          GEMINI_API_KEY: ${{ secrets.GEMINI_API_KEY }}
        run: python -m scraper.main

      - name: Commit feedback history
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git add data/feedback.json || true
          git diff --cached --quiet || git commit -m "chore: update feedback history"
          git push
```

## config.yaml Template

```yaml
# Customise this file — no code changes needed

# What to collect (pre-filter before AI)
filters:
  required_keywords: []      # item must contain at least one
  blocked_keywords: []       # item must not contain any

# Your priorities — AI uses these for scoring
priorities:
  - "example priority 1"
  - "example priority 2"

# Storage
storage:
  provider: "notion"         # notion | sheets | supabase | sqlite

# Feedback learning
feedback:
  positive_statuses: ["Saved", "Applied", "Interested"]
  negative_statuses: ["Skip", "Rejected", "Not relevant"]

# AI settings
ai:
  enabled: true
  model: "gemini-2.5-flash"
  min_score: 0               # filter out items below this score
  rate_limit_seconds: 7      # seconds between API calls
  batch_size: 5              # items per API call
```

## requirements.txt

```
requests==2.31.0
beautifulsoup4==4.12.3
lxml==5.1.0
defusedxml==0.7.1
python-dotenv==1.0.1
pyyaml==6.0.2
notion-client==2.2.1   # if using Notion
# playwright==1.40.0   # uncomment for JS-rendered sites
```

## Free Tier Limits Reference

| Service | Free Limit | Typical Usage |
|---|---|---|
| Gemini Flash Lite | 30 RPM, 1500 RPD | ~56 req/day at 3-hr intervals |
| Gemini 2.0 Flash | 15 RPM, 1500 RPD | Good fallback |
| Gemini 2.5 Flash | 10 RPM, 500 RPD | Use sparingly |
| GitHub Actions | Unlimited (public repos) | ~20 min/day |
| Notion API | Unlimited | ~200 writes/day |
| Supabase | 500MB DB, 2GB transfer | Fine for most agents |
| Google Sheets API | 300 req/min | Works for small agents |
