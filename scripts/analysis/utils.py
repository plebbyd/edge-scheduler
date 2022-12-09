import json
import pandas as pd
import matplotlib.pyplot as plt
from sage_data_client import query

pd.set_option('mode.chained_assignment',None)


def get_data(vsn, start='-1h', end=''):
    if end != "":
        return query(
        start=start,
        end=end,
        filter={
            "name": "sys.scheduler.status.plugin.*",
            "vsn": vsn.upper(),
            }
        )
    else:
        return query(
        start=start,
        filter={
            "name": "sys.scheduler.status.plugin.*",
            "vsn": vsn.upper(),
            }
        )


def parse_events(df):
    v = []
    for _, row in df.iterrows():
        r = json.loads(row.value)
        r["timestamp"] = row.timestamp.isoformat()
        r["node"] = row["meta.node"]
        r["vsn"] = row["meta.vsn"]
        r["event"] = row["name"]
        v.append(r)
    return pd.read_json(json.dumps(v))


def fill_completion_failure(df):
    # Filter only events related to plugin execution
    launched = df[df.event.str.contains("launched")]
    completed = df[df.event.str.contains("complete")]
    failed = df[df.event.str.contains("failed")]
    # launched.loc[launched["k3s_pod_name"] == completed["k3s_pod_name"]]
    for index, p in launched.iterrows():
        found = completed[completed.k3s_pod_name == p.k3s_pod_name]
        if len(found) > 0:
            launched.loc[index, "completed_at"] = found.iloc[0].timestamp
            launched.loc[index, "runtime"] = found.iloc[0].timestamp - p.timestamp
            launched.loc[index, "end_state"] = "completed"
        else:
            found = failed[failed.k3s_pod_name == p.k3s_pod_name]
            if len(found) > 0:
                launched.loc[index, "failed_at"] = found.iloc[0].timestamp
                launched.loc[index, "reason"] = found.iloc[0].reason
                print(launched.loc[p, "reason"])
                if "error_log" in found.iloc[0]:
                    launched.loc[index, "end_state"] = found.iloc[0]["error_log"]
                launched.loc[index, "end_state"] = "failed"
            else:
                launched.loc[index, "end_state"] = "unknown"
    return launched