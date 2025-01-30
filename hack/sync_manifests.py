import os
import yaml
import requests

# Get the latest version of Kueue.
api_url = "https://api.github.com/repos/kubernetes-sigs/kueue/releases/latest"
response = requests.get(api_url)
latest_version = response.json().get("tag_name", "").lstrip("v")  # Strip 'v' from version
if not latest_version:
    raise ValueError("Failed to fetch the latest version")

output_dir = f"assets/kueue/{latest_version}"
manifests_url = f"https://github.com/kubernetes-sigs/kueue/releases/download/v{latest_version}/manifests.yaml"
input_file = f"{output_dir}/manifests.yaml"

# Check if version already exists and skip processing.
if os.path.exists(output_dir):
    print(f"Version {latest_version} already exists. Skipping YAML split.")
    exit(0)

# Ensure output directory exists.
os.makedirs(output_dir, exist_ok=True)

# Download manifests.yaml.
print(f"Downloading manifests.yaml from {manifests_url} ...")
response = requests.get(manifests_url)
if response.status_code != 200:
    raise RuntimeError(f"Failed to download manifests.yaml (HTTP {response.status_code})")

with open(input_file, "wb") as f:
    f.write(response.content)

# Read and split YAML file.
with open(input_file, "r") as f:
    docs = list(yaml.safe_load_all(f))

# Define a mapping of 'kind' to filenames.
file_map = {
    "CustomResourceDefinition": "kueue.yaml",
    "ClusterRole": "clusterrole.yaml",
    "ClusterRoleBinding": "clusterrolebinding.yaml",
    "APIService": "apiservice.yaml",
    "ValidatingWebhookConfiguration": "validatingwebhook.yaml",
    "MutatingWebhookConfiguration": "mutatingwebhook.yaml",
    "Service": "service.yaml",
    "Role": "role.yaml",
    "RoleBinding": "rolebinding.yaml",
}

separated_manifests = {}

# Sort documents into respective files
allowed_kinds = set(file_map.keys())  # Process only these kinds

for doc in docs:
    if not isinstance(doc, dict) or "kind" not in doc:
        continue  # Skip invalid or empty YAML docs

    kind = doc["kind"]
    if kind not in allowed_kinds:
        continue  # Skip unrelated kinds

    filename = file_map[kind]

    if filename not in separated_manifests:
        separated_manifests[filename] = []
    
    separated_manifests[filename].append(doc)

# Write each kind to its respective file.
for filename, content in separated_manifests.items():
    file_path = os.path.join(output_dir, filename)
    with open(file_path, "w") as f:
        yaml.dump_all(content, f, default_flow_style=False)

# Delete manifests.yaml after processing.
os.remove(input_file)
print(f"Processing complete. YAML manifests saved in {output_dir}/")
