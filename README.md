# EC2 Security Group Tool

This app uses the [AWS Go SDK](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials) to perform some basic operations.

Currently, it queries existing security groups, listing the total number of Incoming IP Permissions, Outgoing IP Permissions and the EC2 Instances using them. Additionally, it outputs an AWS CLI to delete unused security groups.

This can be run with the [shared AWS credentials file](https://github.com/aws/aws-sdk-go/wiki/Getting-Started-Credentials) (more info at [configuring the aws cli](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html))

## Examples

Using the default profile

```
awstool
```

Using a profile

```
AWS_PROFILE=bespoke awstool
```

## output

```bash
$ AWS_PROFILE=bespoke awstool
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

 Cross-compilation by gox
