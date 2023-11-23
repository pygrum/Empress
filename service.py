from monarch.builder import *
import os
import base64

def routine(req: BuildRequest) -> BuildResponse:
    agent_id = req.params["id"]
    cmd = "go build -o " + req.params["outfile"] + " -ldflags \"-X main.AgentID=" + agent_id + "\" ."
    result = os.popen(cmd)
    data = result.read()
    status = result.close()
    if status != None:
	    status = 1
    else:
	    status = 0
    with open(req.params["outfile"], "rb") as f:
        data = f.read()

    res = BuildResponse(status, "", base64.b64encode(data).decode('utf-8'))
    return res


if __name__ == "__main__":
    build_server = builder_service(BuildFunction(routine))
    build_server.serve_forever()
