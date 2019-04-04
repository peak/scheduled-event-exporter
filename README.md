# AWS events exporter

This application exports information about AWS scheduled events.
It only supports currently IAM roles, and only requires 1 IAM permission.

#### IAM permission required

`ec2:DescribeInstanceStatus`

#### Exported information

Exports metrics with:

| label   | value  |  
|---|---|
| instance_id of firing event  | Hours to the scheduled event  | 

####  Prometheus Alert examples

```bazaar
  - alert: AWS_Scheduled_Event
    expr: aws_events_scheduled_events_status < 96
    for: 10m
    labels:
      notify: %%TEAM%%
      severity: critical
    annotations:
      summary: "Ec2 instance {{ $labels.instance_name }} is scheduled for event (current value: {{ $value }})"
  - alert: AWS_Scheduled_Event
    expr: aws_events_scheduled_events_status < 200
    for: 10m
    labels:
      notify: %%TEAM%%
      severity: warning
    annotations:
      summary: "Ec2 instance {{ $labels.instance_name }} is scheduled for event (current value: {{ $value }})"

```

#### Docker

You can test quickly with running following command:

`docker run peakcom/scheduled-event-exporter:latest`
