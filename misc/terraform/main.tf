resource "aws_instance" "rinha" {
  ami           = var.ami
  instance_type = var.instance_type
  subnet_id     = var.subnet_id

  iam_instance_profile = aws_iam_instance_profile.ssm_profile.name

  security_groups = [
    aws_security_group.main.id
  ]

  tags = {
    Name = format("%s-host", var.setup_prefix)
  }
}