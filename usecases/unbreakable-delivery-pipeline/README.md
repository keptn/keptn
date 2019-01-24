# Unbreakable Delivery Pipeline

The overall goal of the *Unbreakable Delivery Pipeline* is to implement a pipeline that prevents bad code changes from impacting your real end users. Therefore, it relies on three concepts known as Shift-Left, Shift-Right, and Self-Healing:

* **Shift-Left**: Ability to pull data for specific entities (processes, services, applications) through an Automation API and feed it into the tools that are used to decide on whether to stop the pipeline or keep it running.

* **Shift-Right**: Ability to push deployment information and meta data to your monitoring environment, e.g: differentiate BLUE vs GREEN deployments, push build or revision number of deployment, notify about configuration changes.

* **Self-Healing**: Ability for smart auto-remediation that addresses the root cause of a problem and not the symptom 