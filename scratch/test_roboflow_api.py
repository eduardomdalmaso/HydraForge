import urllib.request
import json

api_key = "rf_ZJl8kwtWoDRx3ZnYrEYeUox8uLm2"

projects = [
    ("leo-ueno", "people-detection-o4rdr"),
    ("human-v2", "human-dataset-v2"),
    ("lennys-workspace-gabu9", "people-kahye"),
]

for ws, proj in projects:
    url = f"https://api.roboflow.com/{ws}/{proj}?api_key={api_key}"
    print(f"\nQuerying {ws}/{proj}...")
    try:
        req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0'})
        res = urllib.request.urlopen(req)
        data = json.loads(res.read().decode('utf-8'))
        p = data.get('project', {})
        name = p.get('name')
        versions = p.get('versions', [])
        print(f"Project Name: {name} | Total Versions: {len(versions)}")
        if versions:
            latest = versions[0]
            print(f"  Latest Version: {latest.get('version')} | Images: {latest.get('images')} | Splits: {latest.get('splits')}")
    except Exception as e:
        print(f"Error querying {ws}/{proj}: {e}")
