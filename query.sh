#!/bin/bash
#curl -s http://127.0.0.1:8080/nodemonitor/api/hpnodeinfo
curl -s -H "Content-Type:application/json" -X POST http://127.0.0.1:8080/nodemonitor/api/nodeinfo 
#curl -s -H "Content-Type:application/json" -X POST --data '{"nodetype":"hpnode"}' http://127.0.0.1:8080/monitor/api/nodeinfo 
#curl -s -H "Content-Type:application/json" -X POST -d '{"nodetype":"prenode", "status":"online"}' http://127.0.0.1:8080/monitor/api/nodeinfo 
