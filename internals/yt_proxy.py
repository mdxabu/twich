import pytchat
import sys
import json

video_id = sys.argv[1]

try:
    chat = pytchat.create(video_id=video_id)
    while chat.is_alive():
        for c in chat.get().sync_items():
            # Standardize output for Go
            data = {"author": c.author.name, "message": c.message}
            print(json.dumps(data))
            sys.stdout.flush() # CRITICAL: Without this, Go waits forever
except Exception as e:
    # This will be caught by our new Stderr scanner in Go
    print(f"Python Error: {e}", file=sys.stderr)
    sys.exit(1)