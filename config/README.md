# Setting Cloud Scheduler

The cloud scheduler exposes APIs via its http server to accept jobs. [Kubernetes ingress](cloudscheduler/cloudscheduler-ingress.yaml) configures an ingress for the cloud scheduler port to be exposed outside the cluster.

[Cloudscheduler deployment](cloudscheduler/cloudscheduler.yaml) deploys cloud scheduler in Kubernetes cluster.

The cloud scheduler pushes jobs (i.e., science goals) to managed node schedulers running on each node via RabbitMQ. [configure.sh](cloudscheduler/configure.sh) creates an account for the cloud scheduler to do so.

# Setting Node Scheduler

Waggle edge stack should already have [configured](https://github.com/waggle-sensor/waggle-edge-stack/blob/main/kubernetes/wes-plugin-scheduler.yaml) the cluster for the node scheduler

_Note: The files under [nodescheduler](./nodescheduler) are only for local testing and may not be up-to-date_