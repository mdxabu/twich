"""
Copyright © 2026 @mdxabu
"""

import json
import re
import sys

import httpx
import pytchat


def resolve_live_video_id(channel_input):
    """
    Given a YouTube channel identifier, find the live stream video ID.

    Accepts:
      - A raw video ID (11 chars, e.g. "JFfPyuo67E8")
      - A @handle (e.g. "@mkbhd" or "mkbhd")
      - A channel URL (e.g. "https://www.youtube.com/@mkbhd")

    Returns the video ID string.
    """
    # If it looks like a plain video ID (11 alphanumeric-ish chars), return it directly.
    if re.match(r"^[A-Za-z0-9_-]{11}$", channel_input):
        return channel_input

    # Normalize input to a @handle if it isn't already a URL.
    if channel_input.startswith("http://") or channel_input.startswith("https://"):
        # Already a URL — append /live if not present.
        base_url = channel_input.rstrip("/")
        if not base_url.endswith("/live"):
            base_url += "/live"
        live_url = base_url
    else:
        # Treat as a handle/username.
        handle = channel_input.lstrip("@")
        live_url = f"https://www.youtube.com/@{handle}/live"

    print(f"Resolving live stream from: {live_url}", file=sys.stderr)
    sys.stderr.flush()

    headers = {
        "User-Agent": (
            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
            "AppleWebKit/537.36 (KHTML, like Gecko) "
            "Chrome/120.0.0.0 Safari/537.36"
        ),
        "Accept-Language": "en-US,en;q=0.9",
    }

    try:
        with httpx.Client(follow_redirects=True, timeout=15.0) as client:
            response = client.get(live_url, headers=headers)
            response.raise_for_status()
            html = response.text
    except Exception as e:
        raise RuntimeError(f"Failed to fetch channel live page: {e}")

    # Strategy 1: Look for canonical URL with /watch?v= (YouTube sets this on live pages).
    match = re.search(
        r'<link\s+rel="canonical"\s+href="https://www\.youtube\.com/watch\?v=([^"&]+)"',
        html,
    )
    if match:
        return match.group(1)

    # Strategy 2: Look for videoId in the page's embedded JSON/JS data.
    # YouTube embeds {"videoId":"XXXXXXXXXXX"} in several places.
    match = re.search(r'"videoId"\s*:\s*"([A-Za-z0-9_-]{11})"', html)
    if match:
        return match.group(1)

    # Strategy 3: Look for /watch?v= anywhere in the HTML.
    match = re.search(r"/watch\?v=([A-Za-z0-9_-]{11})", html)
    if match:
        return match.group(1)

    # Check if the channel is actually live by looking for liveness indicators.
    if '"isLive":true' not in html and '"isLiveNow":true' not in html:
        raise RuntimeError(
            "This channel does not appear to be live right now. "
            "Please check the channel name and try again when the stream is active."
        )

    raise RuntimeError(
        "Could not extract video ID from the channel's live page. "
        "The page structure may have changed, or the channel may not be live."
    )


def main():
    if len(sys.argv) < 2:
        print(
            "Usage: yt_proxy.py <@channel_handle | video_id | channel_url>",
            file=sys.stderr,
        )
        sys.exit(1)

    channel_input = sys.argv[1]

    try:
        video_id = resolve_live_video_id(channel_input)
        print(f"Resolved video ID: {video_id}", file=sys.stderr)
        sys.stderr.flush()
    except RuntimeError as e:
        print(f"Resolution Error: {e}", file=sys.stderr)
        sys.stderr.flush()
        sys.exit(1)

    try:
        chat = pytchat.create(video_id=video_id)
        while chat.is_alive():
            for c in chat.get().sync_items():
                data = {"author": c.author.name, "message": c.message}
                print(json.dumps(data))
                sys.stdout.flush()
    except Exception as e:
        print(f"Python Error: {e}", file=sys.stderr)
        sys.stderr.flush()
        sys.exit(1)


if __name__ == "__main__":
    main()
