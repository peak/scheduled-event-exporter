# AWS events exporter

This application exports information about AWS scheduled events.
It only supports currently IAM roles, and doesnt need any other permission.

IAM permission required:
`ec2:DescribeInstanceStatus`

Travis Build

[![Build Status](https://travis-ci.org/Kronin-Cloud/aws-events-exporter.svg?branch=master)](https://travis-ci.org/Kronin-Cloud/aws-events-exporter)

Tests

[![Coverage Status](https://coveralls.io/repos/github/Kronin-Cloud/aws-events-exporter/badge.svg?branch=master)](https://coveralls.io/github/Kronin-Cloud/aws-events-exporter?branch=master)