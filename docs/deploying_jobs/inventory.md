# Inventory
The inventory contains information about our agents. The concept is very similar to Ansible's [inventory](https://docs.ansible.com/ansible/latest/inventory_guide/intro_inventory.html) file.

```yaml title="inventory.yaml"
nodes:
   # Custom name for a node. Can be set to any string. 
 - name: "lab1"

   # The IP and port where the agent is running.
   location: "192.168.0.99:8080"

   # The agent's API key
   key: "Z6wC4goD7V2EiL4XuecTuo8jVxfvwVxs"

 - name: "lab2"
   location: "192.168.0.100:8080"
   key: "EcwxaMO3kyBaKETesxInx7ga3Ti93gvI"

groups:
   # Custom name for a group. Can be set to any string.
 - name: "labs"

   # List of nodes you want to put inside the group.
   nodes: ["lab1", "lab2"]
```
