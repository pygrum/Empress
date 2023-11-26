import random
import subprocess

from monarch.builder import *
from subprocess import STDOUT, check_output
import os
import base64


def routine(req: BuildRequest) -> BuildResponse:
    out = "empress"

    agent_id = req.params.get("id")  # definitely exists
    host = req.params.get("host")
    port = req.params.get("port")
    ci = req.params.get("callback_interval")
    cs = req.params.get("callback_salt")
    interval = req.params.get("interval")
    salt = req.params.get("salt")
    # if true, then debug is true. otherwise false
    debug = str(req.params.get("debug")).lower() == "true"
    mode = req.params.get("mode")

    if mode == "session":
        mode = 0
    else:
        mode = 1

    with open("config/config.json", "w") as f:
        j = {
            "id": agent_id,
            "host": host,
            "port": port,
            "callback_interval": int(ci),
            "callback_salt": int(cs),
            "interval": int(interval),
            "salt": int(salt),
            "mode": mode,
        }
        j_string = json.dumps(j)
        f.write(j_string)

    if os.path.exists(out):
        try:
            os.remove(out)
        # if dir has the same name or whatever
        except OSError:
            out = random.randbytes(10).decode('utf-8')

    cmd = "go build -o %s " % out
    if debug:
        cmd += "-tags debug "
    cmd += "."

    data = b""
    my_env = os.environ.copy()
    my_env["GOOS"] = req.params.get("os")
    my_env["GOARCH"] = req.params.get("arch")

    try:
        error = check_output(cmd.split(), stderr=STDOUT, env=my_env)
    except subprocess.CalledProcessError as e:
        error = e.output.decode('utf-8')

    # path doesn't exist if build failed
    if not os.path.exists(out):
        status = 1
    else:
        status = 0
        error = ""
        with open(out, "rb") as f:
            data = f.read()

    res = BuildResponse(status, error, base64.b64encode(data).decode('utf-8'))
    return res


if __name__ == "__main__":
    build_server = builder_service(BuildFunction(routine))
    build_server.serve_forever()
