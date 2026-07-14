# Terraform Modules

## Directory Structure

```
terraform/
├── environments/
│   ├── dev/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── outputs.tf
│   │   └── terraform.tfvars
│   ├── staging/
│   │   └── ...
│   └── prod/
│       └── ...
├── modules/
│   ├── networking/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   ├── eks/
│   │   └── ...
│   ├── rds/
│   │   └── ...
│   ├── redis/
│   │   └── ...
│   ├── kafka/
│   │   └── ...
│   ├── s3/
│   │   └── ...
│   ├── monitoring/
│   │   └── ...
│   └── iam/
│       └── ...
├── backend.tf
├── providers.tf
└── versions.tf
```

## Module Breakdown

### Networking Module
```hcl
# modules/networking/main.tf
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "${var.environment}-vpc"
  }
}

resource "aws_subnet" "private" {
  count             = length(var.private_subnet_cidrs)
  vpc_id            = aws_vpc.main.id
  cidr_block        = var.private_subnet_cidrs[count.index]
  availability_zone = var.availability_zones[count.index]

  tags = {
    Name = "${var.environment}-private-${count.index}"
  }
}

resource "aws_subnet" "public" {
  count             = length(var.public_subnet_cidrs)
  vpc_id            = aws_vpc.main.id
  cidr_block        = var.public_subnet_cidrs[count.index]
  availability_zone = var.availability_zones[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.environment}-public-${count.index}"
  }
}
```

### RDS Module
```hcl
# modules/rds/main.tf
resource "aws_db_instance" "postgres" {
  identifier = "${var.environment}-payment-db"

  engine         = "postgres"
  engine_version = "16.3"
  instance_class = var.instance_class

  allocated_storage     = var.allocated_storage
  storage_encrypted     = true
  storage_type          = "gp3"

  db_name  = "paymentdb"
  username = var.db_username
  password = var.db_password

  multi_az               = var.multi_az
  backup_retention_period = var.backup_retention_days
  backup_window          = "02:00-03:00"
  maintenance_window     = "sun:04:00-sun:05:00"

  vpc_security_group_ids = var.security_group_ids
  db_subnet_group_name   = aws_db_subnet_group.main.name

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  deletion_protection = true
  skip_final_snapshot = false
  final_snapshot_identifier = "${var.environment}-final-${formatdate("YYYY-MM-DD", timestamp())}"

  tags = {
    Environment = var.environment
  }
}
```

### Redis Module
```hcl
# modules/redis/main.tf
resource "aws_elasticache_replication_group" "redis" {
  replication_group_id = "${var.environment}-redis"
  description         = "Redis cluster for ${var.environment}"

  engine         = "redis"
  engine_version = "7.1"
  node_type      = var.node_type

  num_cache_clusters = var.num_shards * 2  # primary + replica per shard

  parameter_group_name       = "default.redis7.cluster.on"
  port                       = 6379
  automatic_failover_enabled = true
  multi_az_enabled           = true

  subnet_group_name  = aws_elasticache_subnet_group.main.name
  security_group_ids = var.security_group_ids

  at_rest_encryption_enabled = true
  transit_encryption_enabled  = true

  tags = {
    Environment = var.environment
  }
}
```

### S3 Module
```hcl
# modules/s3/main.tf
resource "aws_s3_bucket" "main" {
  bucket = "${var.environment}-openpayment-${var.purpose}"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "main" {
  bucket = aws_s3_bucket.main.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "main" {
  bucket = aws_s3_bucket.main.id

  rule {
    id     = "expire-after-90-days"
    status = var.expire_days != null ? "Enabled" : "Disabled"

    expiration {
      days = var.expire_days
    }
  }
}

resource "aws_s3_bucket_public_access_block" "main" {
  bucket = aws_s3_bucket.main.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}
```

## State Management
```hcl
# terraform/environments/prod/backend.tf
terraform {
  backend "s3" {
    bucket         = "openpayment-terraform-state"
    key            = "prod/terraform.tfstate"
    region         = "ap-south-1"
    encrypt        = true
    dynamodb_table = "terraform-state-lock"
  }
}
```

## CI/CD for Terraform
```yaml
# .github/workflows/terraform.yml
name: Terraform
on:
  pull_request:
    paths: ['terraform/**']
  push:
    branches: [main]
    paths: ['terraform/**']

jobs:
  plan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: hashicorp/setup-terraform@v3

      - name: Terraform Init
        run: terraform init

      - name: Terraform Format Check
        run: terraform fmt -check

      - name: Terraform Plan
        run: terraform plan -no-color

  apply:
    if: github.ref == 'refs/heads/main'
    needs: [plan]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: hashicorp/setup-terraform@v3

      - name: Terraform Init
        run: terraform init

      - name: Terraform Apply
        run: terraform apply -auto-approve
```
