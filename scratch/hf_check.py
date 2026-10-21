import os
from huggingface_hub import HfApi

token = None
with open('.env') as f:
    for line in f:
        if line.startswith('HF_TOKEN='):
            token = line.strip().split('=', 1)[1].strip('"\' ')

if not token:
    print("HF_TOKEN not found in .env")
    exit(1)

api = HfApi(token=token)
user = api.whoami()
username = user.get('name')
print(f"Logged in user: {username}")

print("\n--- User Models ---")
for m in api.list_models(author=username):
    print(f" - {m.id}")

for org in ['hydra-vision', 'hydravision']:
    try:
        print(f"\n--- Org {org} Models ---")
        for m in api.list_models(author=org):
            print(f" - {m.id}")
    except Exception as e:
        print(f"Error checking {org}: {e}")
