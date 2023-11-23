from monarch.builder import *
import os


def routine(req: BuildRequest) -> BuildResponse:
    agent_id = req.params["id"]
    cmd = "go build -o " + req.params["outfile"] + " -ldflags \"-X main.AgentID=" + agent_id + "\" ."
    result = os.popen(cmd)
    data = result.read()
    status = result.close()
    if status != 0:
        status = 1
    else:
        with open(req.params["outfile"]) as f:
            data = f.read()

    res = BuildResponse(status, "", data.encode())
    return res


if __name__ == "__main__":
    build_server = builder_service(BuildFunction(routine))
    build_server.serve_forever()
