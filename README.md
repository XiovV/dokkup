```json
{
  "groups": [
    {
      "group": "rest-api",
      "containers": ["instance-1", "instance-2", "instance-3"],
      "image_tag": "2.6.0",
      "endpoints": [
        {
          "location": "https://node1.test:8888",
          "node_name": "node1"
        },
        {
          "location": "https://node2.test:8888",
          "node_name": "node2"
        }
      ]
    }
  ]
}
```