# 🚀 Agent-Substrate GitHub Actions PR Reviewer (`antigravity-cli`)

This repository demonstrates a **Self-Hosted GitHub Actions Runner** deployed as a **stateful Agent-Substrate Actor (`ate-demo-github/gh-runner-1`)** on Google Kubernetes Engine (GKE).

## How It Works
1. **Idle State (`ACTOR_STATE_SUSPENDED`)**: When no Pull Request is being reviewed, the runner Actor is suspended in Google Cloud Storage (`gs://snapshot-substrate-demo-boredabdel-lab/ate-demo-github/`), consuming **0 CPU, 0 RAM, and 0 GKE Worker Pods**.
2. **Instant Wake-on-PR**: When a Pull Request is opened or updated, a GitHub Webhook hits the in-cluster `runner-gateway`, which routes `/wake` through `atenet-router` (`ate-target-actor: ate-demo-github/gh-runner-1`).
3. **Sub-Second Restore & Review**:
   - Agent-Substrate restores the Actor's **gVisor memory (`SNAPSHOT_CONTENT_SCOPE_FULL`)** and **durable volume (`/home/runner/state`)** onto a warm GKE `WorkerPool` pod.
   - The Actor executes the `runs-on: [self-hosted, agent-substrate]` workflow job and runs `antigravity-cli review`.
   - `antigravity-cli` analyzes the PR diff with **Vertex AI Gemini 2.5 Flash** (`boredabdel-lab`) and posts a structured review comment along with live **Substrate Actor telemetry** (showing preserved RAM wake/review counters and durable disk history across suspend/resume cycles).
4. **Auto-Suspend**: Once the job finishes, the Actor automatically suspends itself back to GCS (`ACTOR_STATE_SUSPENDED`), freeing the GKE worker pod.
