# EC2 Security Group Tool

This command-line tool uses the [AWS Go SDK](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials) to perform some basic operations.

Currently, this queries existing security groups, listing the total number of Incoming IP Permissions, Outgoing IP Permissions and the EC2 Instances using them. Additionally, it outputs an AWS CLI to delete unused security groups.

* `awstool instances` - lists instances and associated security groups
* `awstool security-groups` - lists security groups and associated instances, sorted by security groups with instances
* `awstool security-groups with-delete` - as above, but with AWS CLI commands to delete unused security groups (with `--dry-run` flag)

This can be run with the [shared AWS credentials file](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials) (more info at [configuring the aws cli](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html))


## Flags and Actions

### Flag: -region

Specify a region, the default region is `us-east-1`. Note, the region flag must come before any actions

example, specifying the `us-west-1` region

```
awstool -region us-west-1
```

### Action: instances

Will list instances in a region and associated security groups

```
awstool instances
```

Output
```
2015/07/05 13:10:35 AWS Region: us-east-1
2015/07/05 13:10:35 Obtaining instances
2015/07/05 13:10:36 Obtained instances 15
Reservation r-2942d9d7, owner: 461758718275
i-bbaf326c [sg-6bf4a603]
Reservation r-8e8424a2, owner: 461758718275
i-17bfd8f9 [sg-296cc641]
Reservation r-858bd9af, owner: 461758718275
i-f6caf41b [sg-296cc641]
Reservation r-22435bd9, owner: 461758718275
i-da2a9709 [sg-9e402bfa]
Reservation r-0a459277, owner: 461758718275
i-7f46d31d [sg-6bf4a603]
Reservation r-b384139f, owner: 461758718275
i-16a88ff8 [sg-640a9301]
Reservation r-2d5af15c, owner: 461758718275
i-9bfe3fc9 [sg-820990e7, sg-640a9301]
```

### Action: security-groups

Output security groups and instances, sorted by those security groups with instances
```
awstool security-groups
```

or, output dry-run AWS CLI security group delete statements

```
awstool security-groups with-delete
```


## Examples

Using the default profile

```
awstool security-groups
```

Using a profile

```
AWS_PROFILE=bespoke awstool security-groups
```

## output

```bash
$ AWS_PROFILE=bespoke awstool security-groups with-delete
          id                 name  in out   i
 sg-5e866b36        quicklaunch-1   2   0   0
   tcp   22-  22 0.0.0.0/0
   tcp   80-  80 0.0.0.0/0
 sg-4bd5b526         cmdline-test   1   0   0
   tcp   80-  80 0.0.0.0/0
 sg-dc876ab4              default   3   0   0
   icmp   -1-  -1 all
   tcp    0-65535 all
   udp    0-65535 all
 sg-094e6562         doge-launch1   1   0   0
   tcp   22-  22 0.0.0.0/0
 sg-d7406bbc        lite-launch-1   1   0   0
   tcp   22-  22 0.0.0.0/0
 sg-96cdb6fe                 pega   4   0   0
   tcp   80-  80 0.0.0.0/0
   tcp 3389-3389 0.0.0.0/0
   tcp 9090-9090 0.0.0.0/0
   tcp 9443-9443 0.0.0.0/0
aws ec2 delete-security-group --group-id sg-5e866b36 --dry-run
aws ec2 delete-security-group --group-id sg-4bd5b526 --dry-run
aws ec2 delete-security-group --group-id sg-dc876ab4 --dry-run
aws ec2 delete-security-group --group-id sg-094e6562 --dry-run
aws ec2 delete-security-group --group-id sg-d7406bbc --dry-run
aws ec2 delete-security-group --group-id sg-96cdb6fe --dry-run
 ```


## Development Notes

 This golang project uses [gb](http://getgb.io) project structure.

 Cross-compilation with [gox](https://github.com/mitchellh/gox)
